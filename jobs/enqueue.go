package jobs

import (
	"context"
	"fmt"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"

	"core/database"
)

var (
	DEFAULT_NAMESPACE = "enterprise-core"
)

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", ":6379")
	},
}

var Enqueuer = work.NewEnqueuer(DEFAULT_NAMESPACE, redisPool)

var Pool = work.NewWorkerPool(struct{}{}, 25, DEFAULT_NAMESPACE, redisPool)

func Initialize() {
	Register("onboard_tenant", func(q work.Q) error {
		name, ok := q["name"].(string)
		if database.Exist(name) {
			return fmt.Errorf("database already exists")
		}
		if !ok {
			return fmt.Errorf("name is not a string")
		}
		err := database.Create(name)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("store_name is not a string")
		}
		database.MigrateTenants(database.GetTenantConnection(name))
		return nil
	})
	start()
}

func start() error {
	var ctx = context.Background()
	go func() {
		select {
		case <-ctx.Done():
			Stop()
		}
	}()
	Pool.Start()
	return nil
}

func Stop() error {
	Pool.Stop()
	return nil
}

func Register(name string, h func(work.Q) error) error {
	Pool.Job(name, func(job *work.Job) error {
		return h(job.Args)
	})
	return nil
}

func Perform(job string, args work.Q) error {
	_, err := Enqueuer.Enqueue(job, args)
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

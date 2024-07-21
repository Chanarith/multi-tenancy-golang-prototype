package middleware

import (
	"github.com/gin-gonic/gin"

	"core/database"
)

type ForceCentralConnection struct {
}

func NewForceCentralConnection() *ForceCentralConnection {
	return &ForceCentralConnection{}
}

func (f *ForceCentralConnection) Handle(c *gin.Context) {
	Db := database.GetCentralConnection()
	c.Set("DB", Db)
	c.Next()
	sqlx, err := Db.DB()
	if err != nil {
		panic(err)
	}
	sqlx.Close()
}

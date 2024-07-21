package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"core/models"
	models_tenants "core/models/tenants"
)

func connect(dbName string) (db *gorm.DB, err error) {
	connectionString := fmt.Sprintf("host=localhost user=root password=root dbname=%s port=5432 sslmode=disable", dbName)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: connectionString,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}

func Create(dbName string) error {
	db := GetCentralConnection().Exec(fmt.Sprintf("CREATE DATABASE \"%s\" WITH TEMPLATE=template0", dbName))
	dbSQL, err := db.DB()
	if err != nil {
		log.Fatal("Failed to create database.")
	}
	defer dbSQL.Close()
	if db.Error != nil {
		log.Fatal(db.Error)
		return db.Error
	}
	return nil
}

func Exist(dbName string) bool {
	var count int64
	db := GetCentralConnection().Raw(fmt.Sprintf("SELECT COUNT(*) FROM pg_database WHERE datname = '%s'", dbName)).Scan(&count)
	dbSQL, _ := db.DB()
	defer dbSQL.Close()
	return count > 0
}

func GetCentralConnection() *gorm.DB {
	db, err := connect("tenants")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return db
}

func GetTenantConnection(tenantName string) *gorm.DB {
	db, err := connect(tenantName)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return db
}

func MigrateTenants(db *gorm.DB) {
	db.AutoMigrate(&models_tenants.Invoice{}, &models_tenants.Product{})
	dbSQL, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get the sql instance from gormDb.")
	}
	defer dbSQL.Close()
}

func MigrateCentral(db *gorm.DB) {
	db.AutoMigrate(&models.Tenant{})
	dbSQL, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get the sql instance from gormDb.")
	}
	defer dbSQL.Close()
}

package models

import (
	"fmt"
	"log"
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Connect() *gorm.DB {
	host := os.Getenv("PGHOST")
	user := os.Getenv("PGUSER")
	pass := os.Getenv("PGPASSWORD")
	dbName := os.Getenv("PGDATABASE")
	dbPort := os.Getenv("PGPORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, pass, dbName, dbPort)

	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Panic("Unable to connect to the database")
	}

	d.AutoMigrate(&Company{})

	log.Println(":INFO: Database connected")
	return d

}

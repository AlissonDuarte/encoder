package main

import (
	"encoder/application/services"
	"encoder/framework/database"
	"encoder/framework/queue"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var db database.Database

func init() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	autoMigrate, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))

	if err != nil {
		panic(err)
	}
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))

	if err != nil {
		panic(err)
	}

	db.AutoMigrateDb = autoMigrate
	db.Debug = debug

	db.DsnTest = os.Getenv("DSN_TEST")
	db.Dsn = os.Getenv("DSN")
	db.DbType = os.Getenv("DB_TYPE")
	db.DbTypeTest = os.Getenv("DB_TYPE_TEST")
	db.Env = os.Getenv("ENV")
}

func main() {
	messageChannel := make(chan amqp.Delivery)
	jobReturn := make(chan services.JobWorkerResult)

	dbConnection, err := db.Connect()

	if err != nil {
		panic(err)
	}

	defer dbConnection.Close()

	rabbit := queue.NewRabbit()

	ch := rabbit.Connect()

	defer ch.Close()

	jobManager := services.NewJobManager(dbConnection, messageChannel, jobReturn, rabbit)

	jobManager.Start(ch)
}

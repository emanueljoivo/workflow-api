package main

import (
	"github.com/joho/godotenv"
	"github.com/nuveo/log"
	"gitlab.com/emanueljoivo/workflows/api"
	"gitlab.com/emanueljoivo/workflows/messaging"
	"gitlab.com/emanueljoivo/workflows/storage"
	"os"
)

const (
	MsgAddrKey string = "MESSAGING_SERVER_ADDRESS"
	MsgPortKey string = "MESSAGING_SERVER_PORT"
	MsgUserKey string = "MESSAGING_SERVER_USER"
	MsgPasswordKey string = "MESSAGING_SERVER_PASSWORD"
	DBAddrKey     string = "DATABASE_ADDRESS"
	DBPortKey     string = "DATABASE_PORT"
	DBUserKey     string = "DATABASE_USER"
	DBNameKey     string = "DATABASE_NAME"
	DBPasswordKey string = "DATABASE_PASSWORD"
)

func init() {
	createDataDir()
}

func main() {
	const DefaultServerPort = "5000"

	err := godotenv.Load()

	if err != nil {
		log.Warningln("No .env file found. Looking for environment variables")
	}

	db := storage.NewStorage(
		os.Getenv(DBAddrKey),
		os.Getenv(DBPortKey),
		os.Getenv(DBUserKey),
		os.Getenv(DBNameKey),
		os.Getenv(DBPasswordKey))

	db.Setup()

	defer db.DB().Close()

	sender := messaging.NewSender(os.Getenv(MsgAddrKey), os.Getenv(MsgPortKey),
		os.Getenv(MsgUserKey), os.Getenv(MsgPasswordKey))

	consumer := messaging.NewConsumer(os.Getenv(MsgAddrKey), os.Getenv(MsgPortKey),
		os.Getenv(MsgUserKey), os.Getenv(MsgPasswordKey))

	a := api.NewAPI(db, sender, consumer)

	if err := a.Start(DefaultServerPort); err != nil {
		log.Warningln(err.Error())
	}
}

func createDataDir() {
	if _, err := os.Stat("data/"); os.IsNotExist(err) {
		err = os.Mkdir("data/", os.ModePerm)
		if err != nil {
			log.Errorf("%s\n", err.Error())
		}
	}
}
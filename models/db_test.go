package models

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// load env variables
	err := godotenv.Load("../env/test.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Connect()
	m.Run()
	Close()
}

package main

import (
	"fmt"
	"os"

	"github.com/duds-fw/go-badger-orm/badgerorm"
	"github.com/sirupsen/logrus"
)

func main() {
	config := badgerorm.Config{
		DBPath:     "data",
		LogLevel:   logrus.DebugLevel,
		LogOutput:  "console",
		MemoryMode: true,
		SyncWrites: true,
	}
	db, _ := badgerorm.NewBadgerORM(config)
	defer db.Close()

	if len(os.Args) < 2 {
		fmt.Println("Usage: badgercli <command> <table> <key> [value]")
		return
	}

	command := os.Args[1]
	table := os.Args[2]
	key := os.Args[3]

	switch command {
	case "get":
		var result map[string]interface{}
		err := db.Get(table, key, &result)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Result:", result)
		}
	case "delete":
		err := db.Delete(table, key)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Deleted:", key)
		}
	default:
		fmt.Println("Invalid command")
	}
}

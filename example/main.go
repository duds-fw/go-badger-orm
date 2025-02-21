package main

import (
	"fmt"
	"time"

	"github.com/duds-fw/go-badger-orm/badgerorm"
)

type User struct {
	Name string
	Age  int
}

func main() {
	db, _ := badgerorm.NewBadgerORM("data")
	defer db.Close()

	// Save a user
	user := User{Name: "Alice", Age: 30}
	db.Save("users", "1", user, time.Hour)

	// Get user
	var retrieved User
	db.Get("users", "1", &retrieved)
	fmt.Println("Retrieved User:", retrieved)

	// Delete user
	db.Delete("users", "1")
}

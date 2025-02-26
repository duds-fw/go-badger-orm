package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/duds-fw/go-badger-orm/badgerorm"
	"github.com/sirupsen/logrus"
)

// User represents a user in the system
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Example usage of Indexing
func indexUsers(orm *badgerorm.BadgerORM) {
	users := []User{
		{ID: "1", Email: "alice@example.com", Role: "admin"},
		{ID: "2", Email: "bob@example.com", Role: "user"},
		{ID: "3", Email: "alice@example.com", Role: "user"},
	}

	for _, user := range users {
		err := orm.Index("users", user.Email, user.ID)
		if err != nil {
			log.Printf("Failed to index user %s: %v", user.ID, err)
		}
		err = orm.Index("users", user.Role, user.ID)
		if err != nil {
			log.Printf("Failed to index user %s: %v", user.ID, err)
		}
	}
}

// Example usage of Querying
func queryUsers(orm *badgerorm.BadgerORM) {
	email := "alice@example.com"
	recordKeys, err := orm.QueryIndex("users", email)
	if err != nil {
		log.Printf("Failed to query index by email %s: %v", email, err)
	} else {
		fmt.Printf("Users indexed by email '%s': %v\n", email, recordKeys)
	}

	role := "user"
	recordKeys, err = orm.QueryIndex("users", role)
	if err != nil {
		log.Printf("Failed to query index by role %s: %v", role, err)
	} else {
		fmt.Printf("Users indexed by role '%s': %v\n", role, recordKeys)
	}
}

func main() {
	config := badgerorm.Config{
		DBPath:     "data",
		LogLevel:   logrus.DebugLevel,
		LogOutput:  "console",
		MemoryMode: true,
		SyncWrites: true,
	}
	db, err := badgerorm.NewBadgerORM(config)
	orm := db
	if err != nil {
		log.Fatalf("Failed to initialize BadgerORM: %v", err)
	}
	defer db.Close()

	// Save a user
	user := User{ID: "1", Email: "alice@example.com", Role: "User"}
	user2 := User{ID: "2", Email: "alice2@example.com", Role: "User"}
	db.Save("users", user.ID, user, time.Hour)
	db.Save("users", user2.ID, user2, time.Hour)

	// Index users
	indexUsers(db)

	// Query indexed users
	queryUsers(db)

	// Get user
	var retrieved User
	db.Get("users", "1", &retrieved)
	result, _ := json.Marshal(retrieved)
	fmt.Println("Retrieved User:", string(result))

	// Delete user
	db.Delete("users", "1")
	fmt.Println("Deleting User:", 1)

	// Get user
	var retrieved2 User
	db.Get("users", "1", &retrieved2)
	result, _ = json.Marshal(retrieved2)
	fmt.Println("Retrieved User:", string(result))

	// Save some users
	db.Save("users", "user1", User{Name: "Alice", Email: "alice@example.com", Role: "admin"}, 0)
	db.Save("users", "user2", User{Name: "Bob", Email: "bob@example.com", Role: "user"}, 0)
	db.Save("users", "user3", User{Name: "Charlie", Email: "charlie@example.com", Role: "user"}, 0)

	// Query with prefix
	usersWithPrefix, _ := db.QueryPrefix("users", "user")
	fmt.Println("Users with prefix 'user':", usersWithPrefix)

	// Query range
	usersInRange, _ := db.QueryRange("users", "user1", "user3")
	fmt.Println("Users in range 'user1' to 'user3':", usersInRange)

	// Query with pagination
	paginatedUsers, _ := db.QueryWithPagination("users", 0, 2) // Page 0, 2 items per page
	fmt.Println("Paginated users (page 0):", paginatedUsers)

	// Example data for batch insert
	records := map[string]interface{}{
		"user1": map[string]string{"name": "Alice", "email": "alice@example.com"},
		"user2": map[string]string{"name": "Bob", "email": "bob@example.com"},
	}

	// Batch Insert
	if err := orm.BatchInsert("users", records, 24*time.Hour); err != nil {
		log.Fatalf("Batch insert failed: %v", err)
	}
	fmt.Println("Batch insert successful")

	// Example data for batch update
	updates := map[string]interface{}{
		"user1": map[string]string{"name": "Alice Smith", "email": "alice.smith@example.com"},
		"user2": map[string]string{"name": "Bob Johnson", "email": "bob.johnson@example.com"},
	}

	// Batch Update
	if err := orm.BatchUpdate("users", updates); err != nil {
		log.Fatalf("Batch update failed: %v", err)
	}
	fmt.Println("Batch update successful")

	// Example keys for batch delete
	keysToDelete := []string{"user1", "user2"}

	// Batch Delete
	if err := orm.BatchDelete("users", keysToDelete); err != nil {
		log.Fatalf("Batch delete failed: %v", err)
	}
	fmt.Println("Batch delete successful")
}

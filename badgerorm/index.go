package badgerorm

import (
	"encoding/json"

	"fmt"

	"github.com/dgraph-io/badger/v4"
)

// Index stores index data for efficient querying, supporting composite and multi-value indexes
func (orm *BadgerORM) Index(table string, indexKey string, recordKeys ...string) error {
	return orm.db.Update(func(txn *badger.Txn) error {
		// Create a composite key for the index
		indexCompositeKey := GenerateIndexKey(table, indexKey)

		// Retrieve existing record keys for the index
		var existingKeys []string
		item, err := txn.Get([]byte(indexCompositeKey))
		if err == nil {
			// If the index already exists, unmarshal the existing keys
			var existingValues []string
			if err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &existingValues)
			}); err != nil {
				return err
			}
			existingKeys = existingValues
		}

		// Combine existing keys with new record keys
		allKeys := append(existingKeys, recordKeys...)

		// Marshal the combined keys to JSON
		data, err := json.Marshal(allKeys)
		if err != nil {
			return err
		}

		// Store the updated index
		return txn.Set([]byte(indexCompositeKey), data)
	})
}

// QueryIndex retrieves indexed record keys
func (orm *BadgerORM) QueryIndex(table, indexKey string) ([]string, error) {
	var recordKeys []string
	err := orm.db.View(func(txn *badger.Txn) error {
		indexCompositeKey := GenerateIndexKey(table, indexKey)
		item, err := txn.Get([]byte(indexCompositeKey))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &recordKeys)
		})
	})
	return recordKeys, err
}

// Helper function to generate index keys
func GenerateIndexKey(table, indexKey string) string {
	return fmt.Sprintf("index:%s:%s", table, indexKey)
}

package badgerorm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
)

// BatchInsert inserts multiple records in a single transaction
func (orm *BadgerORM) BatchInsert(table string, records map[string]interface{}, ttl time.Duration) error {
	return orm.db.Update(func(txn *badger.Txn) error {
		for key, value := range records {
			data, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("failed to marshal value for key %s: %v", key, err)
			}

			e := badger.NewEntry([]byte(fmt.Sprintf("%s:%s", table, key)), data)
			if ttl > 0 {
				e.WithTTL(ttl)
			}
			if err := txn.SetEntry(e); err != nil {
				return fmt.Errorf("failed to set entry for key %s: %v", key, err)
			}
		}
		return nil
	})
}

// BatchUpdate updates multiple records in a single transaction
func (orm *BadgerORM) BatchUpdate(table string, records map[string]interface{}) error {
	return orm.db.Update(func(txn *badger.Txn) error {
		for key, value := range records {
			data, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("failed to marshal value for key %s: %v", key, err)
			}

			if err := txn.Set([]byte(fmt.Sprintf("%s:%s", table, key)), data); err != nil {
				return fmt.Errorf("failed to update entry for key %s: %v", key, err)
			}
		}
		return nil
	})
}

// BatchDelete deletes multiple records in a single transaction
func (orm *BadgerORM) BatchDelete(table string, keys []string) error {
	return orm.db.Update(func(txn *badger.Txn) error {
		for _, key := range keys {
			if err := txn.Delete([]byte(fmt.Sprintf("%s:%s", table, key))); err != nil {
				return fmt.Errorf("failed to delete entry for key %s: %v", key, err)
			}
		}
		return nil
	})
}

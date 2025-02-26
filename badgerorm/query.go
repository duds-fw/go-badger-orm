package badgerorm

import (
	"bytes"
	"encoding/json"
	"fmt"

	"time"

	"github.com/dgraph-io/badger/v4"
)

// Save data with a key
func (orm *BadgerORM) Save(table, key string, value interface{}, ttl time.Duration) error {
	if ttl < 0 {
		orm.logger.Warnf("Negative TTL provided for key: %s:%s", table, key)
		return fmt.Errorf("invalid TTL: %v", ttl)
	}

	data, err := json.Marshal(value)
	if err != nil {
		orm.logger.Errorf("Failed to marshal value for key: %s:%s, error: %v", table, key, err)
		return err
	}

	err = orm.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(GenerateIndexKey(table, key)), data)
		if ttl > 0 {
			e.WithTTL(ttl)
		}
		return txn.SetEntry(e)
	})

	if err == nil {
		orm.logger.Infof("Saved key: %s:%s", table, key)
	} else {
		orm.logger.Errorf("Failed to save key: %s:%s, error: %v", table, key, err)
	}
	return err
}

// Get data by key
func (orm *BadgerORM) Get(table, key string, result interface{}) error {
	return orm.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(GenerateIndexKey(table, key)))
		if err != nil {
			orm.logger.Errorf("Failed to get key: %s:%s, error: %v", table, key, err)
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, result)
		})
	})
}

// Delete a key
func (orm *BadgerORM) Delete(table, key string) error {
	err := orm.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(GenerateIndexKey(table, key)))
	})

	if err == nil {
		orm.logger.Infof("Deleted key: %s:%s", table, key)
	} else {
		orm.logger.Errorf("Failed to delete key: %s:%s, error: %v", table, key, err)
	}
	return err
}

// QueryPrefix retrieves all records that start with a given prefix
func (orm *BadgerORM) QueryPrefix(table, prefix string) ([]string, error) {
	var results []string
	err := orm.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefixKey := []byte(fmt.Sprintf("%s:%s", table, prefix))
		for it.Seek(prefixKey); it.ValidForPrefix(prefixKey); it.Next() {
			item := it.Item()
			var value string
			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &value)
			})
			if err != nil {
				return err
			}
			results = append(results, value)
		}
		return nil
	})
	return results, err
}

// QueryRange retrieves records within a specified range of keys
func (orm *BadgerORM) QueryRange(table, startKey, endKey string) ([]string, error) {
	var results []string
	err := orm.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		start := []byte(GenerateIndexKey(table, startKey))
		end := []byte(GenerateIndexKey(table, endKey))
		for it.Seek(start); it.Valid() && bytes.Compare(it.Item().Key(), end) <= 0; it.Next() {
			item := it.Item()
			results = append(results, string(item.Key()))
		}
		return nil
	})
	return results, err
}

// PaginatedQuery retrieves records with pagination support
func (orm *BadgerORM) QueryWithPagination(table string, page, pageSize int) ([]string, error) {
	var results []string
	err := orm.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		start := page * pageSize
		count := 0
		for it.Seek([]byte(fmt.Sprintf("%s:", table))); it.Valid(); it.Next() {
			if count >= start && count < start+pageSize {
				item := it.Item()
				var value string
				err := item.Value(func(val []byte) error {
					return json.Unmarshal(val, &value)
				})
				if err != nil {
					return err
				}
				results = append(results, value)
			}
			count++
		}
		return nil
	})
	return results, err
}

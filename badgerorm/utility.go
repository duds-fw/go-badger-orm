package badgerorm

import (
	"github.com/dgraph-io/badger/v4"
)

// CountRecords returns the total number of records in a specified table
func (orm *BadgerORM) CountRecords(table string) (int, error) {
	count := 0
	err := orm.db.View(func(txn *badger.Txn) error {
		prefix := []byte(table + ":")
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			count++
		}
		return nil
	})
	return count, err
}

// GetAllKeys retrieves all keys for a given table
func (orm *BadgerORM) GetAllKeys(table string) ([]string, error) {
	var keys []string
	err := orm.db.View(func(txn *badger.Txn) error {
		prefix := []byte(table + ":")
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			keys = append(keys, string(item.Key()))
		}
		return nil
	})
	return keys, err
}

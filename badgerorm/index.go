package badgerorm

import (
	"github.com/dgraph-io/badger/v4"
)

// Transaction executes a set of operations atomically
func (orm *BadgerORM) Transaction(fn func(txn *badger.Txn) error) error {
	return orm.db.Update(fn)
}

// Indexes stores index data for efficient querying
func (orm *BadgerORM) Index(table, indexKey, recordKey string) error {
	return orm.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("index:"+table+":"+indexKey), []byte(recordKey))
	})
}

// QueryIndex retrieves indexed keys
func (orm *BadgerORM) QueryIndex(table, indexKey string) ([]string, error) {
	var keys []string
	return keys, orm.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("index:" + table + ":" + indexKey)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			keys = append(keys, string(item.Key()))
		}
		return nil
	})
}

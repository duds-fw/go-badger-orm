package badgerorm

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/dgraph-io/badger/v4"
)

// Backup exports the entire database to a JSON file
func (orm *BadgerORM) Backup(filePath string) error {
	data := make(map[string]interface{})

	err := orm.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())
			var value interface{}
			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &value)
			})
			if err != nil {
				return err
			}
			data[key] = value
		}
		return nil
	})

	if err != nil {
		return err
	}

	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, file, 0644)
}

// Restore imports data from a JSON file into the database
func (orm *BadgerORM) Restore(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make(map[string]interface{})
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}

	return orm.db.Update(func(txn *badger.Txn) error {
		for key, value := range data {
			val, err := json.Marshal(value)
			if err != nil {
				return err
			}
			if err := txn.Set([]byte(key), val); err != nil {
				return err
			}
		}
		return nil
	})
}

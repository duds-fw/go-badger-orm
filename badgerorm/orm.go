package badgerorm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

// BadgerORM struct
type BadgerORM struct {
	db     *badger.DB
	logger *logrus.Logger
}

// NewBadgerORM initializes the database
func NewBadgerORM(dbPath string) (*BadgerORM, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})

	opts := badger.DefaultOptions(dbPath).WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &BadgerORM{db: db, logger: logger}, nil
}

// Save data with a key
func (orm *BadgerORM) Save(table, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = orm.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(fmt.Sprintf("%s:%s", table, key)), data)
		if ttl > 0 {
			e.WithTTL(ttl)
		}
		return txn.SetEntry(e)
	})

	if err == nil {
		orm.logger.Infof("Saved key: %s:%s", table, key)
	}
	return err
}

// Get data by key
func (orm *BadgerORM) Get(table, key string, result interface{}) error {
	return orm.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(fmt.Sprintf("%s:%s", table, key)))
		if err != nil {
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
		return txn.Delete([]byte(fmt.Sprintf("%s:%s", table, key)))
	})

	if err == nil {
		orm.logger.Infof("Deleted key: %s:%s", table, key)
	}
	return err
}

// Close database
func (orm *BadgerORM) Close() {
	orm.db.Close()
}

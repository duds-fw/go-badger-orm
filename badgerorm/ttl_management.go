package badgerorm

import (
	"time"

	"github.com/dgraph-io/badger/v4"
)

// CleanupExpiredRecords removes expired records from the database
func (orm *BadgerORM) CleanupExpiredRecords() error {
	return orm.db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if item.ExpiresAt() != 0 && time.Now().UnixNano() > int64(item.ExpiresAt()) {
				if err := txn.Delete(item.Key()); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// QueryRecordsNearExpiration retrieves records that are close to expiration
func (orm *BadgerORM) QueryRecordsNearExpiration(threshold time.Duration) ([]string, error) {
	var nearExpirationKeys []string
	err := orm.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if item.ExpiresAt() != 0 {
				expirationTime := time.Unix(0, int64(item.ExpiresAt()))
				if time.Until(expirationTime) < threshold {
					nearExpirationKeys = append(nearExpirationKeys, string(item.Key()))
				}
			}
		}
		return nil
	})
	return nearExpirationKeys, err
}

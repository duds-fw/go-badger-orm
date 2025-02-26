package badgerorm

import (
	"time"
)

// RunGC triggers Badger's garbage collection in a background goroutine
func (orm *BadgerORM) StartGC(interval time.Duration) {
	go func() {
		for {
			err := orm.db.RunValueLogGC(0.5) // Adjust the threshold as needed
			if err != nil {
				orm.logger.Warnf("Garbage collection error: %v", err)
			}
			time.Sleep(interval)
		}
	}()
}

// ManualGC triggers garbage collection manually with a specified threshold
func (orm *BadgerORM) ManualGC(threshold float64) error {
	return orm.db.RunValueLogGC(threshold)
}

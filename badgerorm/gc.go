package badgerorm

import "time"

// RunGC triggers Badger's garbage collection
func (orm *BadgerORM) RunGC() {
	for {
		err := orm.db.RunValueLogGC(0.5)
		if err != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

package badgerorm

import (
	"os"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

// Config holds the configuration options for BadgerORM
type Config struct {
	DBPath     string
	LogLevel   logrus.Level
	LogOutput  string // e.g., "console" or "file"
	MemoryMode bool   // Use memory mode for Badger
	SyncWrites bool   // Enable synchronous writes
}

// BadgerORM struct
type BadgerORM struct {
	db     *badger.DB
	logger *logrus.Logger
	mu     sync.RWMutex // Read-write mutex for concurrency control
}

// NewBadgerORM initializes the database
func NewBadgerORM(config Config) (*BadgerORM, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})

	// Set log level
	logger.SetLevel(config.LogLevel)

	// Configure logging output
	if config.LogOutput == "file" {
		file, err := os.OpenFile("badgerorm.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		logger.SetOutput(file)
	} else {
		logger.SetOutput(os.Stdout)
	}

	// Set Badger options
	opts := badger.DefaultOptions(config.DBPath).WithLoggingLevel(badger.ERROR)
	if config.MemoryMode {
		opts.InMemory = true
	}
	if config.SyncWrites {
		opts.SyncWrites = true
	}

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &BadgerORM{db: db, logger: logger}, nil
}

// Close database
func (orm *BadgerORM) Close() {
	orm.db.Close()
}

// Example function to rebuild indexes in the background
func (orm *BadgerORM) RebuildIndexes() {
	go func() {
		for {
			// Logic to rebuild indexes
			time.Sleep(10 * time.Minute) // Adjust the interval as needed
		}
	}()
}

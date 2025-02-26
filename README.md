[![Release BadgerCLI](https://github.com/duds-fw/go-badger-orm/actions/workflows/release.yml/badge.svg)](https://github.com/duds-fw/go-badger-orm/actions/workflows/release.yml)

# Badger ORM

Badger ORM is a lightweight Object-Relational Mapping (ORM) library for Go, built on top of the Badger key-value database. It provides a simple and efficient way to interact with your data using Go structs.

## Features

- **CRUD Operations**: Easily save, retrieve, update, and delete records.
- **Indexing**: Support for single, composite, and multi-value indexing.
- **Querying Capabilities**: Prefix-based searches, range queries, pagination, and JSON field filtering.
- **Batch Operations**: Bulk insert, update, and delete for efficiency.
- **TTL Management**: Automatic expiration of records with Time-To-Live (TTL) support.
- **Backup & Restore**: Export and import data in JSON format.
- **Concurrency Support**: Safe concurrent access with read-write locks.

## Installation

To install Badger ORM, use the following command:

```bash
go get github.com/duds-fw/go-badger-orm/badgerorm
```

## Usage

### Initialization

```go
package main

import (
	"log"
	"github.com/duds-fw/go-badger-orm/badgerorm"
)

func main() {
	orm, err := badgerorm.NewBadgerORM("path/to/db")
	if err != nil {
		log.Fatalf("Failed to initialize BadgerORM: %v", err)
	}
	defer orm.Close()
}
```

# Documentation

For detailed documentation, please refer to the [Badger ORM WIKI](https://github.com/duds-fw/go-badger-orm/wiki)
.

# Contributing

Contributions are welcome! Please open an issue or submit a pull request.

# License

This library is under **MIT** License.

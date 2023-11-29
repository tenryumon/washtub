package inmemdb

import "github.com/hashicorp/go-memdb"

func InitDB(tables map[string]*memdb.TableSchema) (error, *memdb.MemDB) {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: tables,
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return err, db
	}

	return err, db
}

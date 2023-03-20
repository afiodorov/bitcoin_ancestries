package main

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func createSchema(db *pg.DB) error {
	err := db.CreateTable((*Ancestry)(nil), &orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS hash ON ancestries USING hash (hash);")
	if err != nil {
		return err
	}
	return nil
}

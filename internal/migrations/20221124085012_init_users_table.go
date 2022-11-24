package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upInitUsersTable, downInitUsersTable)
}

func upInitUsersTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create table if not exists users (
			id varchar(36) not null,
			email varchar(225) not null unique,
			username varchar(225),
			password varchar(225) not null,
			tokenhash varchar(15) not null,
			created_at timestamp not null,
			updated_at timestamp not null,
			primary key (id)
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downInitUsersTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop table if exists users;
	`)
	if err != nil {
		return err
	}
	return nil
}

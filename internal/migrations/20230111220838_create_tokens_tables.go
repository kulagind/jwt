package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateTokensTables, downCreateTokensTables)
}

func upCreateTokensTables(tx *sql.Tx) error {
	_, err := tx.Exec(`
		create table if not exists updated_tokens (
			id varchar(36) not null,
			old_token varchar(225) not null unique,
			new_token varchar(225) not null,
			created_at timestamp not null,
			updated_at timestamp not null,
			primary key (id)
		);
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		create table if not exists black_list (
			id varchar(36) not null,
			token varchar(225) not null unique,
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

func downCreateTokensTables(tx *sql.Tx) error {
	_, err := tx.Exec(`
		drop table if exists updated_tokens;
	`)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		drop table if exists black_list;
	`)
	if err != nil {
		return err
	}
	return nil
}

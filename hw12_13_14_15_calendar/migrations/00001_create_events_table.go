package migrations

import (
	"database/sql"
)

func Up0001(tx *sql.Tx) error {
	query := `
		CREATE TABLE events (
			id varchar(36) PRIMARY KEY,
			title varchar(255) NOT NULL,
			starts_at timestamptz NOT NULL DEFAULT current_timestamp,
			duration bigint NOT NULL DEFAULT 0,
			description text NOT NULL DEFAULT '',
			owner_id varchar(36) NOT NULL,
			notify_before bigint NOT NULL DEFAULT 0,
			notification_sent boolean NOT NULL DEFAULT false
		);
	`

	if _, err := tx.Exec(query); err != nil {
		return err
	}

	return nil
}

func Down0001(tx *sql.Tx) error {
	if _, err := tx.Exec("DROP TABLE events;"); err != nil {
		return err
	}

	return nil
}

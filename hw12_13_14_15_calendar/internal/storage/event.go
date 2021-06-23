package storage

import "time"

type Event struct {
	ID           string        `db:"id"`
	Title        string        `db:"title"`
	StartsAt     time.Time     `db:"starts_at"`
	Duration     time.Duration `db:"duration"`
	Description  string        `db:"description"`
	OwnerID      string        `db:"owner_id"`
	NotifyBefore time.Duration `db:"notify_before"`
}

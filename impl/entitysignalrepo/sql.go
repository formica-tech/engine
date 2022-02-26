package entitysignalrepo

import (
	"database/sql"
	"github.com/hamzali/formica-engine/domain"
	"github.com/lib/pq"
)

var copyColumns = []string{"entity_id", "event", "payload", "timestamp"}

const table = "formic_signal"

type EntitySignalSqlRepo struct {
	db *sql.DB
}

func NewSql(db *sql.DB) *EntitySignalSqlRepo {
	return &EntitySignalSqlRepo{db: db}
}

func (d *EntitySignalSqlRepo) BatchSave(signal []domain.EntitySignal) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	str := pq.CopyIn(table, copyColumns...)
	stmt, err := tx.Prepare(str)
	if err != nil {
		return err
	}

	for _, s := range signal {
		_, err = stmt.Exec(s.EntityID, s.Event, s.Payload, s.Timestamp)
		if err != nil {
			return err
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

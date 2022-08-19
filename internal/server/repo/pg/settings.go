package pg

import (
	"context"
	"database/sql"
	"time"
)

func (p *Pg) openDB(pgConn string) error {
	db, err := sql.Open("pgx", pgConn)
	if err != nil {
		return err
	}
	p.db = db
	return err
}

func (p *Pg) setDefaultConfig() error {
	if p.db == nil {
		return ErrHaveNotOpenedDBConnection
	}

	p.db.SetMaxOpenConns(20)
	p.db.SetMaxIdleConns(20)
	p.db.SetConnMaxIdleTime(time.Second * 30)
	p.db.SetConnMaxLifetime(time.Minute * 2)

	return nil
}

func (p *Pg) setGetMetricStmt(ctx context.Context) error {
	if p.db == nil {
		return ErrHaveNotOpenedDBConnection
	}

	stmt, err := p.db.PrepareContext(ctx,
		`SELECT
				id,
				type,
				delta,
				value
			FROM
				metrics
			WHERE id=$1`)
	if err != nil {
		return err
	}
	p.getMetricStmt = stmt
	return nil
}

func (p *Pg) setSetMetricStmt(ctx context.Context) error {
	if p.db == nil {
		return ErrHaveNotOpenedDBConnection
	}

	stmt, err := p.db.PrepareContext(ctx,
		`INSERT INTO metrics
			(id, type, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			id=$1, type=$2, delta=$3, value=$4`)
	if err != nil {
		return err
	}
	p.setMetricStmt = stmt
	return nil
}

func (p *Pg) setGetDataStmt(ctx context.Context) error {
	if p.db == nil {
		return ErrHaveNotOpenedDBConnection
	}

	stmt, err := p.db.PrepareContext(ctx,
		`SELECT
				id, 
				type, 
				delta,
				value 
			FROM 
				metrics`)
	if err != nil {
		return err
	}
	p.getDataStmt = stmt
	return nil
}

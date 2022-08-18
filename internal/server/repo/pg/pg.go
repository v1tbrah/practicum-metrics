package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

type PgStorage struct {
	dbPool *pgxpool.Pool
}

// New returns new postgres storage.
func New(pgConn string) (*PgStorage, error) {
	log.Debug().Str("pgConn", pgConn).Msg("pg.New started")
	defer log.Debug().Msg("pg.New ended")

	ctx := context.Background()
	dbPool, err := pgxpool.Connect(ctx, pgConn)
	if err != nil {
		return nil, err
	}

	_, err = dbPool.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS metrics (
				id varchar(255) PRIMARY KEY, 
				type  varchar(30) NOT NULL, 
				delta bigint, 
				value double PRECISION)`)

	if err != nil {
		return nil, err
	}

	return &PgStorage{dbPool: dbPool}, nil
}

func (p *PgStorage) GetMetric(ctx context.Context, ID string) (metric.Metrics, bool, error) {
	log.Debug().Str("MID", ID).Msg("pg.GetMetric started")
	defer log.Debug().Msg("pg.GetMetric ended")

	resultMetric := metric.Metrics{}

	row := p.dbPool.QueryRow(ctx,
		`SELECT
				id,
				type,
				delta,
				value
			FROM
				metrics
			WHERE id=$1`,
		ID)

	var delta sql.NullInt64
	var value sql.NullFloat64
	err := row.Scan(&resultMetric.ID, &resultMetric.MType, &delta, &value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return resultMetric, false, nil
		} else {
			return resultMetric, false, err
		}
	}

	if resultMetric.MType == "gauge" && value.Valid {
		resultMetric.Value = &value.Float64
	} else if resultMetric.MType == "counter" && delta.Valid {
		resultMetric.Delta = &delta.Int64
	}

	return resultMetric, true, nil
}

func (p *PgStorage) SetMetric(ctx context.Context, thisMetric metric.Metrics) error {
	log.Debug().Str("thisMetric", fmt.Sprint(thisMetric)).Msg("pg.SetMetric started")
	defer log.Debug().Msg("pg.SetMetric ended")

	tx, err := p.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx,
		`INSERT INTO metrics
				(id, type, delta, value)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE SET
				id=$1, type=$2, delta=$3, value=$4`,
		thisMetric.ID, thisMetric.MType, thisMetric.Delta, thisMetric.Value); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *PgStorage) SetListMetrics(ctx context.Context, listMetrics []metric.Metrics) error {
	log.Debug().Str("listMetrics", fmt.Sprint(listMetrics)).Msg("pg.SetListMetrics started")
	defer log.Debug().Msg("pg.SetListMetrics ended")

	tx, err := p.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, currMetric := range listMetrics {
		if _, err = tx.Exec(ctx,
			`INSERT INTO metrics
					(id, type, delta, value)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (id) DO UPDATE SET
					id=$1, type=$2, delta=$3, value=$4`,
			currMetric.ID, currMetric.MType, currMetric.Delta, currMetric.Value); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (p *PgStorage) GetData(ctx context.Context) (*model.Data, error) {
	log.Debug().Msg("pg.GetData started")
	defer log.Debug().Msg("pg.GetData ended")

	rows, err := p.dbPool.Query(ctx,
		`SELECT
				id, 
				type, 
				delta,
				value 
			FROM 
				metrics`)
	if err != nil {
		return nil, err
	}

	newData := model.NewData()
	for rows.Next() {
		currMetric := metric.Metrics{}
		var delta sql.NullInt64
		var value sql.NullFloat64
		rows.Scan(&currMetric.ID, &currMetric.MType, &delta, &value)
		if currMetric.MType == "gauge" && value.Valid {
			currMetric.Value = &value.Float64
		} else if currMetric.MType == "counter" && delta.Valid {
			currMetric.Delta = &delta.Int64
		}
		newData.Metrics[currMetric.ID] = currMetric
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return newData, nil
}

func (p *PgStorage) Ping(ctx context.Context) error {
	log.Debug().Msg("pg.Ping started")
	defer log.Debug().Msg("pg.Ping ended")

	return p.dbPool.Ping(ctx)
}

func (p *PgStorage) ClosePoolConn() {
	log.Debug().Msg("pg.CloseConnection started")
	defer log.Debug().Msg("pg.CloseConnection ended")

	p.dbPool.Close()
}

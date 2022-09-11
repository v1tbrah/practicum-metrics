package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"

	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
)

var ErrHaveNotOpenedDBConnection = errors.New("haven`t opened the database connection")

type Pg struct {
	db            *sql.DB
	stmtGetMetric *sql.Stmt
	stmtSetMetric *sql.Stmt
	stmtGetData   *sql.Stmt
}

// New returns new postgres storage.
func New(pgConn string) (*Pg, error) {
	log.Debug().Str("pgConn", pgConn).Msg("pg.New started")
	var err error
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("pg.New ended")
		} else {
			log.Debug().Msg("pg.New ended")
		}
	}()

	newPg := Pg{}

	db, err := sql.Open("pgx", pgConn)
	if err != nil {
		return nil, err
	}
	newPg.db = db

	newPg.db.SetMaxOpenConns(20)
	newPg.db.SetMaxIdleConns(20)
	newPg.db.SetConnMaxIdleTime(time.Second * 30)
	newPg.db.SetConnMaxLifetime(time.Minute * 2)

	ctx := context.Background()
	if err = newPg.initTables(ctx); err != nil {
		return nil, err
	}

	newPg.stmtGetMetric, err = newPg.db.PrepareContext(ctx, stmtGetMetric)
	if err != nil {
		return nil, err
	}

	newPg.stmtSetMetric, err = newPg.db.PrepareContext(ctx, stmtSetMetric)
	if err != nil {
		return nil, err
	}

	newPg.stmtGetData, err = newPg.db.PrepareContext(ctx, stmtGetData)
	if err != nil {
		return nil, err
	}

	return &newPg, nil
}

func (p *Pg) GetMetric(ctx context.Context, ID string) (metric.Metrics, bool, error) {
	log.Debug().Str("MID", ID).Msg("pg.GetMetric started")
	var err error
	defer func() {
		if err != nil {
			log.Debug().Msg("pg.GetMetric ended")
		} else {
			log.Error().Err(err).Msg("pg.GetMetric ended")
		}
	}()

	resultMetric := metric.Metrics{}

	row := p.stmtGetMetric.QueryRowContext(ctx, ID)

	var delta sql.NullInt64
	var value sql.NullFloat64
	err = row.Scan(&resultMetric.ID, &resultMetric.MType, &delta, &value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

func (p *Pg) SetMetric(ctx context.Context, thisMetric metric.Metrics) error {
	log.Debug().Str("thisMetric", fmt.Sprint(thisMetric)).Msg("pg.SetMetric started")
	var err error
	defer func() {
		if err != nil {
			log.Debug().Msg("pg.SetMetric ended")
		} else {
			log.Error().Err(err).Msg("pg.SetMetric ended")
		}
	}()

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txStmt := tx.StmtContext(ctx, p.stmtSetMetric)
	if _, err = txStmt.ExecContext(ctx,
		thisMetric.ID, thisMetric.MType, thisMetric.Delta, thisMetric.Value); err != nil {
		return err
	}

	return tx.Commit()
}

func (p *Pg) SetListMetrics(ctx context.Context, listMetrics []metric.Metrics) error {
	log.Debug().Str("listMetrics", fmt.Sprint(listMetrics)).Msg("pg.SetListMetrics started")
	var err error
	defer func() {
		if err != nil {
			log.Debug().Msg("pg.SetListMetrics ended")
		} else {
			log.Error().Err(err).Msg("pg.SetListMetrics ended")
		}
	}()

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txStmt := tx.StmtContext(ctx, p.stmtSetMetric)

	for _, currMetric := range listMetrics {
		if _, err = txStmt.ExecContext(ctx,
			currMetric.ID, currMetric.MType, currMetric.Delta, currMetric.Value); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (p *Pg) GetData(ctx context.Context) (model.Data, error) {
	log.Debug().Msg("pg.GetData started")
	var err error
	defer func() {
		if err != nil {
			log.Debug().Msg("pg.GetData ended")
		} else {
			log.Error().Err(err).Msg("pg.GetData ended")
		}
	}()

	rows, err := p.stmtGetData.QueryContext(ctx)
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
		newData[currMetric.ID] = currMetric
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return newData, nil
}

func (p *Pg) Ping(ctx context.Context) error {
	log.Debug().Msg("pg.Ping started")
	var err error
	defer func() {
		if err != nil {
			log.Debug().Msg("pg.Ping ended")
		} else {
			log.Error().Err(err).Msg("pg.Ping ended")
		}
	}()

	err = p.db.PingContext(ctx)
	return err
}

func (p *Pg) ClosePoolConn() {
	log.Debug().Msg("pg.CloseConnection started")
	defer log.Debug().Msg("pg.CloseConnection ended")

	p.db.Close()
}

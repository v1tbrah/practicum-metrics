package pg

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
	"log"
)

type PgStorage struct {
	Data   *model.Data
	dbPool *pgxpool.Pool
}

// New returns new postgres storage.
func New(pgConn string) *PgStorage {

	dbPool, err := pgxpool.Connect(context.Background(), pgConn)
	if err != nil {
		log.Fatalln("Unable to connect to database:", err)
	}
	//defer dbPool.Close() TODO

	_, err = dbPool.Exec(context.Background(),
		"CREATE TABLE IF NOT EXISTS metrics ( "+
			"id varchar(255) PRIMARY KEY, "+
			"type  varchar(30) NOT NULL, "+
			"delta bigint, "+
			"value double PRECISION "+
			")")

	if err != nil {
		log.Fatalln("Creating table Metrics. Error. Reason:", err)
	}

	return &PgStorage{Data: model.NewData(), dbPool: dbPool}
}

func (p *PgStorage) GetMetric(ID string) (metric.Metrics, bool, error) {
	thisMetric := metric.Metrics{}

	row := p.dbPool.QueryRow(context.Background(), "SELECT id, type, delta, value FROM metrics WHERE id=$1", ID)

	var delta sql.NullInt64
	var value sql.NullFloat64
	err := row.Scan(&thisMetric.ID, &thisMetric.MType, &delta, &value)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return thisMetric, false, nil
		} else {
			return thisMetric, false, err
		}
	}

	if thisMetric.MType == "gauge" && value.Valid {
		thisMetric.Value = &value.Float64
	} else if thisMetric.MType == "counter" && delta.Valid {
		thisMetric.Delta = &delta.Int64
	}

	return thisMetric, true, nil
}

func (p *PgStorage) SetMetric(ID string, thisMetric metric.Metrics) error {

	ctx := context.Background()
	tx, err := p.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, "DELETE FROM metrics WHERE id=$1", ID); err != nil && err != sql.ErrNoRows {
		log.Println("#1", err)
		tx.Rollback(ctx)
		return err
	}
	if _, err = tx.Exec(ctx, "INSERT INTO metrics (id, type, delta, value) values ($1, $2, $3, $4)",
		ID, thisMetric.MType, thisMetric.Delta, thisMetric.Value); err != nil && err != sql.ErrNoRows {
		log.Println("#2", err)
		tx.Rollback(ctx)
		return err
	}
	tx.Commit(ctx)
	return nil
}

func (p *PgStorage) GetData() (*model.Data, error) {
	dataFromStorage, err := p.dataFromStorage()
	if err != nil {
		return nil, err
	}
	return dataFromStorage, nil
}

func (p *PgStorage) dataFromStorage() (*model.Data, error) {
	rows, err := p.dbPool.Query(context.Background(),
		"SELECT "+
			"id, "+
			"type, "+
			"delta, "+
			"value "+
			"FROM "+
			"metrics")
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

func (p *PgStorage) CloseConnection() {
	p.dbPool.Close()
}

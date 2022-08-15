package pg

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/v1tbrah/metricsAndAlerting/internal/server/model"
	"github.com/v1tbrah/metricsAndAlerting/pkg/metric"
	"log"
)

type pgStorage struct {
	Data   *model.Data
	dbPool *pgxpool.Pool
}

// New returns new postgres storage.
func New(pgConn string) *pgStorage {

	dbPool, err := pgxpool.Connect(context.Background(), pgConn)
	if err != nil {
		log.Fatalln("Unable to connect to database:", err)
	}
	//defer dbPool.Close() TODO

	_, err = dbPool.Exec(context.Background(),
		"CREATE TABLE IF NOT EXISTS metrics ( "+
			"id varchar(255) PRIMARY KEY, "+
			"type  varchar(30) NOT NULL, "+
			"delta integer, "+
			"value double PRECISION "+
			")")

	if err != nil {
		log.Fatalln("Creating table Metrics. Error. Reason:", err)
	}

	return &pgStorage{Data: model.NewData(), dbPool: dbPool}
}

func (p *pgStorage) GetData() *model.Data {
	return p.Data
}

func (p *pgStorage) SetData(data *model.Data) {
	p.Data = data
}

func (p *pgStorage) RestoreData() error {

	dataFromStorage, err := p.dataFromStorage()
	if err != nil {
		return err
	}
	p.SetData(dataFromStorage)
	return nil
}

func (p *pgStorage) StoreData() error {

	tx, err := p.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}
	if _, err = tx.Exec(context.Background(), "DELETE FROM metrics"); err != nil {
		tx.Rollback(context.Background())
		return err
	}
	for name, currMetric := range p.Data.Metrics {
		if _, err = tx.Exec(context.Background(), "INSERT INTO metrics (id, type, delta, value) values ($1, $2, $3, $4)",
			name, currMetric.MType, currMetric.Delta, currMetric.Value); err != nil {
			tx.Rollback(context.Background())
			return err
		}
	}

	err = tx.Commit(context.Background())
	return err
}

func (p *pgStorage) dataFromStorage() (*model.Data, error) {
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
		if !currMetric.TypeIsValid() {
			log.Println("Restore metric. Error. Reason:", err)
			continue
		}
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

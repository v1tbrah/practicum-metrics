package pg

import (
	"context"

	"github.com/rs/zerolog/log"
)

const (
	queryCreateTableMetrics = `CREATE TABLE IF NOT EXISTS metrics
		(
			id varchar(255) PRIMARY KEY, 
			type  varchar(30) NOT NULL, 
			delta bigint, 
			value double PRECISION
    	);`
)

func (p *Pg) initTables(ctx context.Context) error {
	log.Debug().Msg("pg.initTables started")
	var err error
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("pg.initTables ended")
		} else {
			log.Debug().Msg("pg.initTables ended")
		}
	}()

	if p.db == nil {
		return ErrHaveNotOpenedDBConnection
	}

	_, err = p.db.ExecContext(ctx, queryCreateTableMetrics)

	return err
}

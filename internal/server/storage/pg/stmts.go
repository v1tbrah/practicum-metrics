package pg

const (
	stmtGetMetric = `SELECT
				id,
				type,
				delta,
				value
         FROM
				metrics
		 WHERE id=$1`

	stmtSetMetric = `INSERT INTO metrics
			(id, type, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			id=$1, type=$2, delta=$3, value=$4`

	stmtGetData = `SELECT
				id, 
				type, 
				delta,
				value 
			FROM 
				metrics`
)

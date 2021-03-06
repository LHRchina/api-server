package orm

import (
	"fmt"
	"sync"

	"github.com/go-pg/pg/internal"
)

// Placeholder that is replaced with count(*).
const placeholder = `'_go_pg_placeholder'`

// https://wiki.postgresql.org/wiki/Count_estimate
var pgCountEstimateFunc = fmt.Sprintf(`
CREATE OR REPLACE FUNCTION _go_pg_count_estimate_v2(query text, threshold int)
RETURNS int AS $$
DECLARE
  rec record;
  nrows int;
BEGIN
  FOR rec IN EXECUTE 'EXPLAIN ' || query LOOP
    nrows := substring(rec."QUERY PLAN" FROM ' rows=(\d+)');
    EXIT WHEN nrows IS NOT NULL;
  END LOOP;

  -- Return the estimation if there are too many rows.
  IF nrows > threshold THEN
    RETURN nrows;
  END IF;

  -- Otherwise execute real count query.
  query := replace(query, 'SELECT '%s'', 'SELECT count(*)');
  EXECUTE query INTO nrows;

  IF nrows IS NULL THEN
    nrows := 0;
  END IF;

  RETURN nrows;
END;
$$ LANGUAGE plpgsql;
`, placeholder)

// CountEstimate uses EXPLAIN to get estimated number of rows matching the query.
// If that number is bigger than the threshold it returns the estimation.
// Otherwise it executes another query using count aggregate function and
// returns the result.
//
// Based on https://wiki.postgresql.org/wiki/Count_estimate
func (q *Query) CountEstimate(threshold int) (int, error) {
	if q.stickyErr != nil {
		return 0, q.stickyErr
	}

	query, err := q.countSelectQuery(placeholder).AppendQuery(nil)
	if err != nil {
		return 0, err
	}

	for i := 0; i < 3; i++ {
		var count int
		_, err = q.db.QueryOne(
			Scan(&count),
			"SELECT _go_pg_count_estimate_v2(?, ?)",
			string(query), threshold,
		)
		if err != nil {
			if pgerr, ok := err.(internal.PGError); ok && pgerr.Field('C') == "42883" {
				// undefined_function
				err = q.createCountEstimateFunc()
				if err != nil {
					return 0, err
				}
				continue
			}
		}
		return count, err
	}

	return 0, err
}

func (q *Query) createCountEstimateFunc() error {
	_, err := q.db.Exec(pgCountEstimateFunc)
	return err
}

// SelectAndCountEstimate runs Select and CountEstimate in two goroutines,
// waits for them to finish and returns the result.
func (q *Query) SelectAndCountEstimate(threshold int, values ...interface{}) (count int, err error) {
	if q.stickyErr != nil {
		return 0, q.stickyErr
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var mu sync.Mutex

	go func() {
		defer wg.Done()
		if e := q.Select(values...); e != nil {
			mu.Lock()
			err = e
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		var e error
		count, e = q.CountEstimate(threshold)
		if e != nil {
			mu.Lock()
			err = e
			mu.Unlock()
		}
	}()

	wg.Wait()
	return count, err
}

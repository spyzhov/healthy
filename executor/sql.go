package executor

import (
	"database/sql"
	"errors"
	"fmt"
	_ "runtime/cgo" // necessary for some drivers

	"github.com/spyzhov/healthy/executor/internal"
	. "github.com/spyzhov/healthy/executor/internal/args"
	"github.com/spyzhov/healthy/step"
	"github.com/spyzhov/safe"

	//	region drivers
	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-adodb"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/nakagami/firebirdsql"
	//  endregion
)

type SqlArgs struct {
	Driver  string         `json:"driver"`
	URL     string         `json:"url"`
	SQL     string         `json:"sql"`
	Args    []interface{}  `json:"args"`
	Require SqlArgsRequire `json:"require"`
}

func (e *Executor) Sql(args *SqlArgs) (step.Function, error) {
	if err := args.Validate(); err != nil {
		return nil, safe.Wrap(err, "sql")
	}

	var db *sql.DB
	id := fmt.Sprintf("%s_%s", args.Driver, args.URL)
	reconnect := func() error {
		return e.protected(func() (err error) {
			if old, ok := e.connections[id]; ok {
				defer safe.Close(old, "Executor:db_connections:"+args.Driver+":old")
				delete(e.connections, id)
			}

			db, err = sql.Open(args.Driver, args.URL)
			if err != nil {
				return safe.Wrap(err, "connect")
			}
			e.connections[id] = db
			return nil
		})
	}
	getDB := func() error {
		return e.protected(func() (err error) {
			var ok bool
			if db, ok = e.connections[id]; ok {
				return db.Ping()
			}
			return fmt.Errorf("not found")
		})
	}
	if err := getDB(); err != nil {
		if err := reconnect(); err != nil {
			return nil, safe.Wrap(err, "sql")
		}
	}

	var exec func(force bool) (*sql.Rows, error)
	exec = func(first bool) (rows *sql.Rows, err error) {
		rows, err = db.QueryContext(e.ctx, args.SQL, args.Args...)
		if errors.Is(err, sql.ErrConnDone) || (err != nil && err.Error() == "sql: database is closed") {
			if first {
				rec := reconnect()
				if rec != nil {
					return nil, safe.Wrap(rec, "sql: reconnect")
				}
				return exec(false)
			}
			return nil, err
		}
		return rows, err
	}

	return func() (*step.Result, error) {
		var (
			err    error
			rows   *sql.Rows
			header []string
			data   [][]interface{}
		)
		if rows, err = exec(true); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, safe.Wrap(err, "sql")
			}
		}
		if header, err = rows.Columns(); err != nil {
			return nil, safe.Wrap(err, "sql: column names")
		}
		if errors.Is(err, sql.ErrNoRows) {
			data = make([][]interface{}, 0)
		} else {
			for rows.Next() {
				row := make([]interface{}, len(header))
				ref := make([]interface{}, len(header))
				for i := 0; i < len(header); i++ {
					ref[i] = &row[i]
				}
				err = rows.Scan(ref...)
				if err != nil {
					return nil, safe.Wrap(err, "sql: scan rows")
				}
				data = append(data, row)
			}
		}

		err = args.Require.Match(data)
		if err != nil {
			return step.NewResultError(err.Error()), nil
		}

		return step.NewResultSuccess("OK"), nil
	}, nil
}

func (a *SqlArgs) Validate() (err error) {
	available := sql.Drivers()
	if !internal.StringInSlice(a.Driver, available) {
		return fmt.Errorf("url: driver %s not available, you can use only: %v", a.Driver, available)
	}
	if a.SQL == "" {
		return fmt.Errorf("sql: blank request")
	}
	if err = a.Require.Validate(); err != nil {
		return safe.Wrap(err, "require")
	}
	return nil
}

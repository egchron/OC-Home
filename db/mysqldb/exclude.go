package mysqldb

import (
	"database/sql"
	"time"

	"okumoto.net/db"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func (m *MysqlDB) AddExclude(s *db.ExcludeEntry) (bool, error) {

	duration := s.Duration / time.Second
	q := "" +
		"INSERT INTO Exclude (" +
		"	`desc`,`start`,`end`,`duration`,`interval`,`units`" +
		") VALUES (?, ?, ?, ?, ?, ?)"

	result, err := m.Exec(q,
		s.Desc, s.Start, s.End, duration, s.Units, s.Interval)
	switch e := err.(type) {
	case nil:
		// handle after switch.
	case *mysql.MySQLError:
		if e.Number == mysqlerr.ER_DUP_ENTRY {
			return true, nil // duplicate entry
		}
		return false, errors.Wrap(err, "AddExclude")

	default:
		return false, errors.Wrapf(err, "AddExclude Unexpected %T", e)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return false, errors.Wrap(err, "AddExclude LastInsertId")
	}

	s.Id = id
	return false, nil
}

func (m *MysqlDB) GetExcludes() ([]db.ExcludeEntry, error) {

	q := "" +
		"SELECT " +
		"	`id`, `desc`, `start`,`end`,`duration`,`interval`,`units`" +
		"FROM Exclude"

	rows, err := m.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "GetExcludes Query")
	}
	defer rows.Close()

	excludes := make([]db.ExcludeEntry, 0)
	for rows.Next() {
		var s db.ExcludeEntry
		var end sql.NullTime
		var duration int32
		var interval sql.NullInt32
		var units sql.NullString
		err := rows.Scan(
			&s.Id,
			&s.Desc,
			&s.Start,
			&duration,
			&end,
			&interval,
			&units,
		)
		if err != nil {
			return nil, errors.Wrap(err, "GetExcludes Scan")
		}

		s.Duration = time.Duration(duration) * time.Second
		if end.Valid {
			s.End = &end.Time
		}
		if interval.Valid {
			s.Interval = &interval.Int32
		}
		if units.Valid {
			var v db.Units
			switch units.String {
			case "seconds":
				v = db.Seconds
			case "minutes":
				v = db.Minutes
			case "hours":
				v = db.Hours
			case "days":
				v = db.Days
			case "weeks":
				v = db.Weeks
			case "months":
				v = db.Months
			}
			s.Units = &v
		}

		excludes = append(excludes, s)
	}

	return excludes, nil
}

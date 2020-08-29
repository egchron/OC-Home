package mysqldb

import (
	"okumoto.net/db"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func (m *MysqlDB) AddSRLink(s *db.SRLinkEntry) (bool, error) {

	q := "" +
		"INSERT INTO SRLink (station, relay, seq)" +
		"VALUES (?, ?, 0)"

	result, err := m.Exec(q, s.Station, s.Relay)
	switch e := err.(type) {
	case nil:
		// handle after switch.
	case *mysql.MySQLError:
		if e.Number == mysqlerr.ER_DUP_ENTRY {
			return true, nil // duplicate entry
		}
		return false, errors.Wrap(err, "AddSRLink Exec")

	default:
		return false, errors.Wrapf(err, "AddSRLink Unexpected %T", e)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return false, errors.Wrap(err, "AddSRLink LastInsertId")
	}

	s.Id = id
	return false, nil
}

func (m *MysqlDB) GetSRLinks() ([]db.SRLinkEntry, error) {

	q := "" +
		"SELECT id, station, relay " +
		"FROM SRLink"

	rows, err := m.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "GetSRLinks Query")
	}
	defer rows.Close()

	links := make([]db.SRLinkEntry, 0)
	for rows.Next() {
		var s db.SRLinkEntry
		err := rows.Scan(
			&s.Id,
			&s.Station,
			&s.Relay,
		)
		if err != nil {
			return nil, errors.Wrap(err, "GetSRLinks Scan")
		}

		links = append(links, s)
	}

	return links, nil
}

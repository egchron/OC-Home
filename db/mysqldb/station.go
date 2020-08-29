package mysqldb

import (
	"time"

	"okumoto.net/db"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func (m *MysqlDB) AddStation(s *db.StationEntry) (bool, error) {

	q := "" +
		"INSERT INTO Station (name, budget)" +
		"VALUES (?, ?)"

	result, err := m.Exec(q, s.Name, s.Budget)
	switch e := err.(type) {
	case nil:
		// handle after switch.
	case *mysql.MySQLError:
		if e.Number == mysqlerr.ER_DUP_ENTRY {
			return true, nil // duplicate entry
		}
		return false, errors.Wrap(err, "AddStation")

	default:
		return false, errors.Wrapf(err, "AddStation Unexpected %T", e)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return false, errors.Wrap(err, "AddStation LastInsertId")
	}

	s.Id = id
	return false, nil
}

func (m *MysqlDB) GetStations() ([]db.StationEntry, error) {

	q := "" +
		"SELECT id, name, budget " +
		"FROM Station"

	rows, err := m.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "GetStations Query")
	}
	defer rows.Close()

	stations := make([]db.StationEntry, 0)
	for rows.Next() {
		var s db.StationEntry
		var budget int
		err := rows.Scan(
			&s.Id,
			&s.Name,
			&budget,
		)
		if err != nil {
			return nil, errors.Wrap(err, "GetStations Scan")
		}
		s.Budget = time.Duration(budget) * time.Second

		stations = append(stations, s)
	}

	return stations, nil
}

func (m *MysqlDB) GetStationByName(name string) (*db.StationEntry, error) {

	q := "" +
		"SELECT id, name " +
		"FROM Station " +
		"WHERE name = ?"

	row := m.QueryRow(q, name)

	var s db.StationEntry
	err := row.Scan(
		&s.Id,
		&s.Name,
	)
	if err != nil {
		return nil, errors.Wrap(err, "GetStations Scan")
	}

	return &s, nil
}

func (m *MysqlDB) GetStationById(id int64) (*db.StationEntry, error) {

	q := "" +
		"SELECT id, name " +
		"FROM Station " +
		"WHERE id = ?"

	row := m.QueryRow(q, id)

	var s db.StationEntry
	err := row.Scan(
		&s.Id,
		&s.Name,
	)
	if err != nil {
		return nil, errors.Wrap(err, "GetStations Scan")
	}

	return &s, nil
}

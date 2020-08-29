package mysqldb

import (
	"okumoto.net/db"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func (m *MysqlDB) AddRelay(s *db.RelayEntry) (bool, error) {

	q := "" +
		"INSERT INTO Relay (name)" +
		"VALUES (?)"

	result, err := m.Exec(q, s.Name)
	switch e := err.(type) {
	case nil:
		// handle after switch.
	case *mysql.MySQLError:
		if e.Number == mysqlerr.ER_DUP_ENTRY {
			return true, nil // duplicate entry
		}
		return false, errors.Wrap(err, "AddRelay Exec")

	default:
		return false, errors.Wrapf(err, "AddRelay Unexpected %T", e)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return false, errors.Wrap(err, "AddRelay LastInsertId")
	}

	s.Id = id
	return false, nil
}

func (m *MysqlDB) GetRelays() ([]db.RelayEntry, error) {

	q := "" +
		"SELECT id, name " +
		"FROM Relay"

	rows, err := m.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "GetRelays Query")
	}
	defer rows.Close()

	relays := make([]db.RelayEntry, 0)
	for rows.Next() {
		var s db.RelayEntry
		err := rows.Scan(
			&s.Id,
			&s.Name,
		)
		if err != nil {
			return nil, errors.Wrap(err, "GetRelays Scan")
		}

		relays = append(relays, s)
	}

	return relays, nil
}

func (m *MysqlDB) GetRelayByName(name string) (*db.RelayEntry, error) {

	q := "" +
		"SELECT id, name " +
		"FROM Relay " +
		"WHERE name = ?"

	row := m.QueryRow(q, name)

	var r db.RelayEntry
	err := row.Scan(
		&r.Id,
		&r.Name,
	)
	if err != nil {
		return nil, errors.Wrap(err, "GetRelays Scan")
	}

	return &r, nil
}

func (m *MysqlDB) GetRelayById(id int64) (*db.RelayEntry, error) {

	q := "" +
		"SELECT id, name " +
		"FROM Relay " +
		"WHERE id = ?"

	row := m.QueryRow(q, id)

	var r db.RelayEntry
	err := row.Scan(
		&r.Id,
		&r.Name,
	)
	if err != nil {
		return nil, errors.Wrap(err, "GetRelays Scan")
	}

	return &r, nil
}

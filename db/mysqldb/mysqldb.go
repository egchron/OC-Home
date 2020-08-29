package mysqldb

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"okumoto.net/db"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type MysqlDB struct {
	*sql.DB
	dbname string
}

// New establishes a connection to the MySQL server.
func New(c db.DBConfig) (*MysqlDB, error) {

	config := mysql.NewConfig()
	config.User = c.User
	config.Passwd = c.Pass
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%s", c.Server, c.Port)
	config.DBName = c.Name

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, errors.Wrap(err, "New Open")
	}

	if err := db.Ping(); err == nil {
		// Great, we can access the DB.
		return &MysqlDB{DB: db, dbname: c.Name}, nil
	}

	//--------------------------------------------------
	// Database needs to be created.  Open a new SQL
	// connection which doesn't have DBName set.
	//--------------------------------------------------
	config.DBName = ""

	db2, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, errors.Wrap(err, "New Open2")
	}
	_, err = db2.Exec("CREATE DATABASE `Water`")
	if err != nil {
		return nil, errors.Wrap(err, "New Exec")
	}
	db2.Close()

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "New Ping2")
	}
	return &MysqlDB{DB: db, dbname: c.Name}, nil
}

func (db *MysqlDB) MaybeSetupDB(schemaFile string) error {

	schema, err := os.Open(schemaFile)
	if err != nil {
		return errors.Wrap(err, "MaybeSetupDB Open")
	}
	defer schema.Close()

	var (
		tableTag       = "-- Table structure for table "
		createTableTag = "CREATE TABLE "
	)

	state := 0
	var tableDef []string
	scanner := bufio.NewScanner(schema)
	for scanner.Scan() {
		line := scanner.Text()
		switch state {
		case 0: // Want "structure for table"
			if strings.HasPrefix(line, tableTag) {
				state = 1 // Collect table def.
			}
		case 1: // Want table def
			if strings.HasPrefix(line, createTableTag) {
				tableDef = append(tableDef, line)
				state = 2
			}
		case 2: // Want end table def
			if strings.HasPrefix(line, "--") {
				if err := db.SetupTable(tableDef); err != nil {
					return errors.Wrap(err, "MaybeSetupDB SetupTable")
				}
				tableDef = []string{}
				state = 0
			} else if strings.HasPrefix(line, "/*!40101 SET ") {
				// Ignore
			} else {
				tableDef = append(tableDef, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "MaybeSetupDB Err")
	}

	if err = db.SetupTable(tableDef); err != nil {
		return errors.Wrap(err, "MaybeSetupDB last SetupTable")
	}

	return nil
}

func (db *MysqlDB) SetupTable(tableDef []string) error {

	_, err := db.Exec(strings.Join(tableDef, "\n"))
	if err != nil {
		return errors.Wrap(err, "SetupTable Exec")
	}

	return nil
}

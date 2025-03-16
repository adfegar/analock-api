package storage

import (
	"database/sql"

	"github.com/adfer-dev/analock-api/database"
	"github.com/adfer-dev/analock-api/utils"
)

type Storage interface {
	Get(uint) (interface{}, error)
	Create(interface{}) error
	Update(interface{}) error
	Delete(uint) error
	Scan(*sql.Rows) (interface{}, error)
}

var databaseConnection *sql.DB = database.GetDatabaseInstance().GetConnection()
var storageLogger *utils.CustomLogger = utils.GetCustomLogger()

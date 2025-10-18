package repository

import (
	"database/sql"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/models"
)

type IRepositoryLog interface {
	GetLogs() ([]models.Log, error)
	GetLogById(id int) (models.Log, error)
	InsertLog(user models.Log) (sql.Result, error)
}

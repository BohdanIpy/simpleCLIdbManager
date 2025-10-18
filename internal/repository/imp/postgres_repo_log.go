package imp

import (
	"database/sql"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/models"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/repository"
)

type PostgresRepoLog struct {
	db *sql.DB
}

func (p PostgresRepoLog) GetLogs() ([]models.Log, error) {
	rows, err := p.db.Query("SELECT * FROM logs;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	logs := make([]models.Log, 0)
	for rows.Next() {
		var log models.Log
		err = rows.Scan(&log.Id, &log.LogTime, &log.LogMessage)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (p PostgresRepoLog) GetLogById(id int) (models.Log, error) {
	var log models.Log
	err := p.db.QueryRow("SELECT * FROM logs WHERE id = $1;", id).Scan(&log.Id, &log.LogTime, &log.LogMessage)
	if err != nil {
		return models.Log{}, err
	}
	return log, nil
}

func (p PostgresRepoLog) InsertLog(user models.Log) (sql.Result, error) {
	res, err := p.db.Exec("INSERT INTO logs (log_time, log_message) VALUES ($1, $2)", user.LogTime, user.LogMessage)
	return res, err
}

func NewPostgresRepoLog(db *sql.DB) repository.IRepositoryLog {
	return &PostgresRepoLog{db: db}
}

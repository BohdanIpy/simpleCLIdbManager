package imp

import (
	"github.com/BohdanIpy/simpleCLIdbManager/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"regexp"
	"testing"
	"time"
)

func logEquals(this models.Log, other models.Log) bool {
	return this.Id == other.Id &&
		this.LogMessage == other.LogMessage
}

func TestPostgresRepoLog_GetLogs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepoLog(db)
	logs := []models.Log{
		models.Log{
			Id:         1,
			LogTime:    time.Now(),
			LogMessage: "msg1",
		},
		models.Log{
			Id:         2,
			LogTime:    time.Now(),
			LogMessage: "msg2",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "log_time", "log_message"}).
		AddRow(logs[0].Id, logs[0].LogTime, logs[0].LogMessage).
		AddRow(logs[1].Id, logs[1].LogTime, logs[1].LogMessage)

	mock.ExpectQuery(`SELECT \* FROM logs;`).
		WillReturnRows(rows)

	data, err := repo.GetLogs()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when getting logs", err)
	}
	if len(data) != len(logs) {
		t.Fatalf("logs length should be %d, but %d", len(data), len(logs))
	}
	for i, log := range data {
		if !logEquals(log, logs[i]) {
			t.Fatalf("logs should be %v, but %v", logs[i], logs[i])
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresRepoLog_GetLogById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepoLog(db)

	log := models.Log{
		Id:         1,
		LogTime:    time.Now(),
		LogMessage: "msg1",
	}

	rows := sqlmock.NewRows([]string{"id", "log_time", "log_message"}).AddRow(log.Id, log.LogTime, log.LogMessage)
	mock.ExpectQuery(`SELECT \* FROM logs WHERE id = \$1;`).
		WithArgs(1).
		WillReturnRows(rows)

	data, err := repo.GetLogById(1)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when getting log", err)
	}
	if !logEquals(data, log) {
		t.Fatalf("logs should be %v, but %v", log, data)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresRepoLog_CreateLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepoLog(db)

	log := models.Log{
		Id:         1,
		LogTime:    time.Now(),
		LogMessage: "msg1",
	}

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO logs (log_time, log_message) VALUES ($1, $2)")).
		WithArgs(sqlmock.AnyArg(), log.LogMessage).
		WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := repo.InsertLog(log)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when inserting log", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when getting affected rows", err)
	}
	if affected != 1 {
		t.Fatalf("affected should be 1, but %d", affected)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

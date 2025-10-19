package middleware

import (
	"github.com/BohdanIpy/simpleCLIdbManager/internal/models"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/repository/imp"
	"github.com/DATA-DOG/go-sqlmock"
	"regexp"
	"testing"
	"time"
)

func userEquals(this models.User, other models.User) bool {
	return this.Email == other.Email &&
		this.ID == other.ID &&
		this.Name == other.Name &&
		this.Password == other.Password &&
		this.RegisteredAt == other.RegisteredAt
}

func TestLoggerMiddleware_GetUsers(t *testing.T) {
	dbLog, mockLog, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbLog.Close()

	dbUser, mockUser, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbUser.Close()

	rUser := imp.NewPostgresRepoUser(dbUser)
	repo := LoggerMiddleware{
		logDb: imp.NewPostgresRepoLog(dbLog),
		next:  rUser,
	}

	userRows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "registered_at",
	}).AddRow(1, "John", "john@example.com", "secret", time.Now())

	mockUser.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users;")).
		WillReturnRows(userRows)

	mockLog.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO logs (log_time, log_message) VALUES ($1, $2);",
	)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockLog.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO logs (log_time, log_message) VALUES ($1, $2);",
	)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(2, 1))

	_, err = repo.GetUsers()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when getting users", err)
	}

	if err := mockLog.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
	if err := mockUser.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestLoggerMiddleware_GetUserById(t *testing.T) {
	dbLog, mockLog, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbLog.Close()

	dbUser, mockUser, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbUser.Close()

	rUser := imp.NewPostgresRepoUser(dbUser)
	repo := NewLoggerMiddleware(imp.NewPostgresRepoLog(dbLog), rUser)

	mUser := models.User{
		ID:           1,
		Name:         "John",
		Email:        "john@example.com",
		Password:     "secret",
		RegisteredAt: time.Now(),
	}

	userRows := sqlmock.NewRows([]string{"id", "name", "email", "password", "registered_at"}).
		AddRow(mUser.ID, mUser.Name, mUser.Email, mUser.Password, mUser.RegisteredAt)

	mockUser.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id = $1;")).
		WithArgs(1).
		WillReturnRows(userRows)

	mockLog.ExpectExec(regexp.QuoteMeta("INSERT INTO logs (log_time, log_message) VALUES ($1, $2);")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockLog.ExpectExec(regexp.QuoteMeta("INSERT INTO logs (log_time, log_message) VALUES ($1, $2);")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	us, err := repo.GetUserById(1)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when getting user by id", err)
	}
	if err := mockLog.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
	if err := mockUser.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
	if !userEquals(us, mUser) {
		t.Fatalf("users are not equal")
	}
}

func TestLoggerMiddleware_InsertUser(t *testing.T) {
	dbLog, mockLog, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbLog.Close()

	dbUser, mockUser, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbUser.Close()

	rUser := imp.NewPostgresRepoUser(dbUser)
	repo := LoggerMiddleware{imp.NewPostgresRepoLog(dbLog), rUser}

	mUser := models.User{
		ID:           1,
		Name:         "John",
		Email:        "john@example.com",
		Password:     "secret",
		RegisteredAt: time.Now(),
	}

	mockLog.ExpectExec(regexp.QuoteMeta("INSERT INTO logs (log_time, log_message) VALUES ($1, $2);")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockLog.ExpectExec(regexp.QuoteMeta("INSERT INTO logs (log_time, log_message) VALUES ($1, $2);")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockUser.ExpectExec(regexp.QuoteMeta("INSERT INTO users (name, email, password, registered_at) VALUES ($1, $2, $3, $4);")).
		WithArgs(mUser.Name, mUser.Email, mUser.Password, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = repo.InsertUser(mUser)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when inserting user", err)
	}
	if err := mockLog.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
	if err := mockUser.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestLoggerMiddleware_UpdateUser(t *testing.T) {
	dbLog, mockLog, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbLog.Close()

	dbUser, mockUser, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbUser.Close()

	rUser := imp.NewPostgresRepoUser(dbUser)
	repo := LoggerMiddleware{imp.NewPostgresRepoLog(dbLog), rUser}

	mUser := models.User{
		ID:           1,
		Name:         "John",
		Email:        "john@example.com",
		Password:     "secret",
		RegisteredAt: time.Now(),
	}

	mockLog.ExpectExec(regexp.QuoteMeta("INSERT INTO logs (log_time, log_message) VALUES ($1, $2);")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockLog.ExpectExec(regexp.QuoteMeta("INSERT INTO logs (log_time, log_message) VALUES ($1, $2);")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockUser.ExpectExec(regexp.QuoteMeta("UPDATE users SET name = $1, email = $2, password = $3, registered_at = $4 WHERE id = $5;")).
		WithArgs(mUser.Name, mUser.Email, mUser.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = repo.UpdateUserById(mUser.ID, mUser)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when inserting user", err)
	}
	if err := mockLog.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
	if err := mockUser.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestLoggerMiddleware_DeleteUser(t *testing.T) {
	dbLog, mockLog, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbLog.Close()

	dbUser, mockUser, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbUser.Close()

	rUser := imp.NewPostgresRepoUser(dbUser)
	repo := LoggerMiddleware{imp.NewPostgresRepoLog(dbLog), rUser}

	mUser := models.User{
		ID:           1,
		Name:         "John",
		Email:        "john@example.com",
		Password:     "secret",
		RegisteredAt: time.Now(),
	}

	mockLog.ExpectExec(regexp.QuoteMeta("INSERT INTO logs (log_time, log_message) VALUES ($1, $2);")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockLog.ExpectExec(regexp.QuoteMeta("INSERT INTO logs (log_time, log_message) VALUES ($1, $2);")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockUser.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = $1;")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id := mUser.ID
	_, err = repo.DeleteUserById(id)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when inserting user", err)
	}
	if err := mockLog.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
	if err := mockUser.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

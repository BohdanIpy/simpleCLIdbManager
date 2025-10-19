package middleware

import (
	"github.com/BohdanIpy/simpleCLIdbManager/internal/models"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/repository/imp"
	"github.com/DATA-DOG/go-sqlmock"
	"regexp"
	"testing"
	"time"
)

func TestTransactionalMiddleware_GetUsers(t *testing.T) {
	dbUser, mockUser, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbUser.Close()

	rUser := imp.NewPostgresRepoUser(dbUser)
	repo := TransactionalMiddleware{
		db:   dbUser,
		next: rUser,
	}

	userRows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "registered_at",
	}).AddRow(1, "John", "john@example.com", "secret", time.Now())

	mockUser.ExpectBegin()

	mockUser.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users;")).
		WillReturnRows(userRows)

	mockUser.ExpectCommit()

	_, err = repo.GetUsers()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when getting users", err)
	}

	if err := mockUser.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestTransactionalMiddleware_GetUserById(t *testing.T) {
	dbUser, mockUser, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbUser.Close()

	rUser := imp.NewPostgresRepoUser(dbUser)
	repo := TransactionalMiddleware{db: dbUser, next: rUser}

	mUser := models.User{
		ID:           1,
		Name:         "John",
		Email:        "john@example.com",
		Password:     "secret",
		RegisteredAt: time.Now(),
	}

	userRows := sqlmock.NewRows([]string{"id", "name", "email", "password", "registered_at"}).
		AddRow(mUser.ID, mUser.Name, mUser.Email, mUser.Password, mUser.RegisteredAt)

	mockUser.ExpectBegin()

	mockUser.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id = $1;")).
		WithArgs(1).
		WillReturnRows(userRows)

	mockUser.ExpectCommit()

	us, err := repo.GetUserById(1)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when getting user by id", err)
	}
	if err := mockUser.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
	if !userEquals(us, mUser) {
		t.Fatalf("users are not equal")
	}
}

func TestTransactionalMiddleware_InsertUser(t *testing.T) {
	dbUser, mockUser, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbUser.Close()

	rUser := imp.NewPostgresRepoUser(dbUser)
	repo := TransactionalMiddleware{db: dbUser, next: rUser}

	mUser := models.User{
		ID:           1,
		Name:         "John",
		Email:        "john@example.com",
		Password:     "secret",
		RegisteredAt: time.Now(),
	}

	mockUser.ExpectBegin()

	mockUser.ExpectExec(regexp.QuoteMeta("INSERT INTO users (name, email, password, registered_at) VALUES ($1, $2, $3, $4);")).
		WithArgs(mUser.Name, mUser.Email, mUser.Password, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockUser.ExpectCommit()

	_, err = repo.InsertUser(mUser)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when inserting user", err)
	}
	if err := mockUser.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestTransactionalMiddleware_UpdateUser(t *testing.T) {
	dbUser, mockUser, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbUser.Close()

	rUser := imp.NewPostgresRepoUser(dbUser)
	repo := TransactionalMiddleware{db: dbUser, next: rUser}

	mUser := models.User{
		ID:           1,
		Name:         "John",
		Email:        "john@example.com",
		Password:     "secret",
		RegisteredAt: time.Now(),
	}

	mockUser.ExpectBegin()

	mockUser.ExpectExec(regexp.QuoteMeta("UPDATE users SET name = $1, email = $2, password = $3, registered_at = $4 WHERE id = $5;")).
		WithArgs(mUser.Name, mUser.Email, mUser.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockUser.ExpectCommit()

	_, err = repo.UpdateUserById(mUser.ID, mUser)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when inserting user", err)
	}
	if err := mockUser.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestTransactionalMiddleware_DeleteUser(t *testing.T) {
	dbUser, mockUser, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbUser.Close()

	rUser := imp.NewPostgresRepoUser(dbUser)
	repo := TransactionalMiddleware{db: dbUser, next: rUser}

	mUser := models.User{
		ID:           1,
		Name:         "John",
		Email:        "john@example.com",
		Password:     "secret",
		RegisteredAt: time.Now(),
	}

	mockUser.ExpectBegin()

	mockUser.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = $1;")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockUser.ExpectCommit()

	id := mUser.ID
	_, err = repo.DeleteUserById(id)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when inserting user", err)
	}
	if err := mockUser.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

package imp

import (
	"github.com/BohdanIpy/simpleCLIdbManager/internal/models"
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

// GetById
func TestPostgresRepoUser_GetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepoUser(db)

	expectedUser := &models.User{
		ID:           1,
		Name:         "John Doe",
		Email:        "john@example.com",
		Password:     "dasfsa",
		RegisteredAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "registered_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.RegisteredAt)

	mock.ExpectQuery(`SELECT \* FROM users WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	user, err := repo.GetUserById(1)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when getting user by id", err)
	}
	if !userEquals(user, *expectedUser) {
		t.Fatalf("got %+v, expected %+v", user, expectedUser)
	}

	if mock.ExpectationsWereMet() != nil {
		t.Fatalf("there were unfulfilled expectations: %s", mock.ExpectationsWereMet())
	}
}

// GetAll
func TestPostgresRepoUser_GetUsers2(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	expectedUsers := []models.User{
		models.User{
			ID:           1,
			Name:         "John Doe",
			Email:        "john@example.com",
			Password:     "dasfsa",
			RegisteredAt: time.Now(),
		},
		models.User{
			ID:           2,
			Name:         "Name",
			Email:        "mail@example.com",
			Password:     "dasfsa",
			RegisteredAt: time.Now(),
		},
	}

	repo := NewPostgresRepoUser(db)

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "registered_at"}).
		AddRow(expectedUsers[0].ID, expectedUsers[0].Name, expectedUsers[0].Email, expectedUsers[0].Password, expectedUsers[0].RegisteredAt).
		AddRow(expectedUsers[1].ID, expectedUsers[1].Name, expectedUsers[1].Email, expectedUsers[1].Password, expectedUsers[1].RegisteredAt)

	mock.ExpectQuery(`SELECT \* FROM users;`).
		WillReturnRows(rows)

	users, err := repo.GetUsers()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when getting users", err)
	}
	if len(users) != len(expectedUsers) {
		t.Fatalf("got %d users, expected %d", len(users), len(expectedUsers))
	}
	for i, user := range users {
		if !userEquals(user, expectedUsers[i]) {
			t.Fatalf("got %+v, expected %+v", user, expectedUsers[i])
		}
	}
	if mock.ExpectationsWereMet() != nil {
		t.Fatalf("there were unfulfilled expectations: %s", mock.ExpectationsWereMet())
	}
}

// InsertUser
func TestPostgresRepoUser_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepoUser(db)

	expectedUser := &models.User{
		ID:           1,
		Name:         "John Doe",
		Email:        "john@example.com",
		Password:     "dasfsa",
		RegisteredAt: time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO users (name, email, password, registered_at) VALUES ($1, $2, $3, $4);",
	)).
		WithArgs(expectedUser.Name, expectedUser.Email, expectedUser.Password, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := repo.InsertUser(*expectedUser)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when inserting user", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when getting affected rows", err)
	}
	if affected != 1 {
		t.Fatalf("got %d rows affected, expected %d", affected, 1)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Fatalf("there were unfulfilled expectations: %s", mock.ExpectationsWereMet())
	}
}

// update user
func TestPostgresRepoUser_UpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepoUser(db)

	user := models.User{
		ID:           1,
		Name:         "Alice",
		Email:        "alice@example.com",
		Password:     "securepassword",
		RegisteredAt: time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE users 
		SET name = $1, email = $2, password = $3, registered_at = $4
		WHERE id = $5
	`)).WithArgs(user.Name, user.Email, user.Password, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.UpdateUserById(1, user)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when updating user", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when getting affected rows", err)
	}
	if affected != 1 {
		t.Fatalf("got %d rows affected, expected %d", affected, 1)
	}
	if mock.ExpectationsWereMet() != nil {
		t.Fatalf("there were unfulfilled expectations: %s", mock.ExpectationsWereMet())
	}
}

// delete by id
func TestPostgresRepoUser_DeleteUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := 1
	repo := NewPostgresRepoUser(db)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = $1;")).
		WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.DeleteUserById(id)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when deleting user", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when getting affected rows", err)
	}
	if affected != 1 {
		t.Fatalf("got %d rows affected, expected %d", affected, 1)
	}
}

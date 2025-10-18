package imp

import (
	"database/sql"
	"errors"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/models"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/repository"
	"time"
)

type PostgresRepoUser struct {
	db *sql.DB
}

func (p PostgresRepoUser) GetUsers() ([]models.User, error) {
	rows, err := p.db.Query("SELECT * FROM users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RegisteredAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (p PostgresRepoUser) GetUserById(id int) (models.User, error) {
	var user models.User
	err := p.db.QueryRow("SELECT * FROM users WHERE id = $1;", id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RegisteredAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, sql.ErrNoRows
		} else {
			return models.User{}, err
		}
	}
	return user, nil
}

func (p PostgresRepoUser) InsertUser(user models.User) (sql.Result, error) {
	res, err := p.db.Exec("INSERT INTO users (name, email, password, registered_at) VALUES ($1, $2, $3, $4);", user.Name, user.Email, user.Password, time.Now())
	return res, err
}

func (p PostgresRepoUser) DeleteUserById(id int) (sql.Result, error) {
	res, err := p.db.Exec("DELETE FROM users WHERE id = $1;", id)
	return res, err
}

func (p PostgresRepoUser) UpdateUserById(id int, user models.User) (sql.Result, error) {
	res, err := p.db.Exec("UPDATE users SET name = $1, email = $2, password = $3, registered_at = $4 WHERE id = $5;", user.Name, user.Email, user.Password, time.Now(), user.ID)
	return res, err
}

func NewPostgresRepoUser(db *sql.DB) repository.IRepositoryUser {
	return &PostgresRepoUser{db: db}
}

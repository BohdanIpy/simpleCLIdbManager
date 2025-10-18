package middleware

import (
	"database/sql"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/models"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/repository"
)

type TransactionalMiddleware struct {
	db   *sql.DB
	next repository.IRepositoryUser
}

func (t *TransactionalMiddleware) GetUsers() ([]models.User, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := t.next.GetUsers()
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (t *TransactionalMiddleware) GetUserById(id int) (models.User, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return models.User{}, err
	}
	defer tx.Rollback()

	result, err := t.next.GetUserById(id)
	if err != nil {
		return models.User{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.User{}, err
	}

	return result, nil
}

func (t *TransactionalMiddleware) InsertUser(user models.User) (sql.Result, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := t.next.InsertUser(user)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (t *TransactionalMiddleware) DeleteUserById(id int) (sql.Result, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := t.next.DeleteUserById(id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (t *TransactionalMiddleware) UpdateUserById(id int, user models.User) (sql.Result, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := t.next.UpdateUserById(id, user)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func NewTransactionalMiddleware(db *sql.DB, next repository.IRepositoryUser) repository.IRepositoryUser {
	return &TransactionalMiddleware{db: db, next: next}
}

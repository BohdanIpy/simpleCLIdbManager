package repository

import (
	"database/sql"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/models"
)

type IRepositoryUser interface {
	GetUsers() ([]models.User, error)
	GetUserById(id int) (models.User, error)
	InsertUser(user models.User) (sql.Result, error)
	UpdateUserById(id int, user models.User) (sql.Result, error)
	DeleteUserById(id int) (sql.Result, error)
}

package middleware

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/BohdanIpy/simpleCLIdbManager/internal/models"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/repository"
)

type LoggerMiddleware struct {
	logDb repository.IRepositoryLog
	next  repository.IRepositoryUser
}

func (l *LoggerMiddleware) safeLog(msg string) {
	if _, err := l.logDb.InsertLog(models.Log{
		LogMessage: msg,
		LogTime:    time.Now(),
	}); err != nil {
		log.Printf("[WARN] failed to insert log: %v", err)
	}
}

func (l *LoggerMiddleware) GetUsers() ([]models.User, error) {
	l.safeLog("Getting all users")

	data, err := l.next.GetUsers()

	status := "succeeded"
	if err != nil {
		status = "failed"
	}
	l.safeLog(fmt.Sprintf("Getting all users %s", status))

	return data, err
}

func (l *LoggerMiddleware) GetUserById(id int) (models.User, error) {
	l.safeLog(fmt.Sprintf("Getting user by id: %d", id))

	data, err := l.next.GetUserById(id)

	status := "succeeded"
	if err != nil {
		status = "failed"
	}
	l.safeLog(fmt.Sprintf("Getting user by id: %d -- %s", id, status))

	return data, err
}

func (l *LoggerMiddleware) InsertUser(user models.User) (sql.Result, error) {
	l.safeLog("Trying to insert user")

	data, err := l.next.InsertUser(user)

	status := "succeeded"
	if err != nil {
		status = "failed"
	}
	l.safeLog(fmt.Sprintf("Insert user %s", status))

	return data, err
}

func (l *LoggerMiddleware) DeleteUserById(id int) (sql.Result, error) {
	l.safeLog(fmt.Sprintf("Started deleting user with id %d", id))

	res, err := l.next.DeleteUserById(id)

	status := "succeeded"
	if err != nil {
		status = fmt.Sprintf("failed: %v", err)
	}
	l.safeLog(fmt.Sprintf("Deleting user with id %d %s", id, status))

	return res, err
}

func (l *LoggerMiddleware) UpdateUserById(id int, user models.User) (sql.Result, error) {
	l.safeLog(fmt.Sprintf("Trying to update user with id %d", id))

	res, err := l.next.UpdateUserById(id, user)

	status := "succeeded"
	if err != nil {
		status = fmt.Sprintf("failed: %v", err)
	}
	l.safeLog(fmt.Sprintf("Updating user with id %d %s", id, status))

	return res, err
}

func NewLoggerMiddleware(repo repository.IRepositoryLog, next repository.IRepositoryUser) repository.IRepositoryUser {
	return &LoggerMiddleware{logDb: repo, next: next}
}

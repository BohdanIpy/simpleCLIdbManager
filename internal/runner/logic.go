package runner

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/models"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/repository"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func Logic(ctx context.Context, repo repository.IRepositoryUser) {
	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down...")
			return
		default:
			printMenu()
			fmt.Print("Enter command: ")

			line, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("error reading input: %v\n", err)
				continue
			}

			cmd := strings.Fields(strings.TrimSpace(line))
			if len(cmd) == 0 {
				continue
			}
			select {
			case <-ctx.Done():
				fmt.Println("Shutting down...")
				return
			default:
				if err := handleCommand(repo, cmd); err != nil {
					fmt.Printf("Error: %v\n\n", err)
				}
			}
		}
	}
}

func printMenu() {
	fmt.Println("\nAvailable operations:")
	fmt.Println("1                            - Get all users")
	fmt.Println("2 <id>                       - Get user by ID")
	fmt.Println("3 <name> <email> <password>  - Insert user")
	fmt.Println("4 <id>                       - Delete user by ID")
	fmt.Println("5 <id> <name> <email> <pwd>  - Update user by ID")
	fmt.Println("q                            - Quit")
}

func handleCommand(repo repository.IRepositoryUser, cmd []string) error {
	switch cmd[0] {
	case "q", "quit", "exit":
		os.Exit(0)

	case "1":
		return handleGetAll(repo)

	case "2":
		if len(cmd) < 2 {
			return errors.New("usage: 2 <id>")
		}
		return handleGetByID(repo, cmd[1])

	case "3":
		if len(cmd) < 4 {
			return errors.New("usage: 3 <name> <email> <password>")
		}
		return handleInsert(repo, cmd[1], cmd[2], cmd[3])

	case "4":
		if len(cmd) < 2 {
			return errors.New("usage: 4 <id>")
		}
		return handleDelete(repo, cmd[1])

	case "5":
		if len(cmd) < 5 {
			return errors.New("usage: 5 <id> <name> <email> <password>")
		}
		return handleUpdate(repo, cmd[1], cmd[2], cmd[3], cmd[4])

	default:
		return fmt.Errorf("unknown command: %s", cmd[0])
	}
	return nil
}

func handleGetAll(repo repository.IRepositoryUser) error {
	users, err := repo.GetUsers()
	if err != nil {
		return err
	}
	for _, u := range users {
		fmt.Println(u)
	}
	return nil
}

func handleGetByID(repo repository.IRepositoryUser, idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	user, err := repo.GetUserById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("User not found")
			return nil
		}
		return err
	}
	fmt.Println(user)
	return nil
}

func handleInsert(repo repository.IRepositoryUser, name, email, password string) error {
	user := models.User{
		Name:         name,
		Email:        email,
		Password:     password,
		RegisteredAt: time.Now(),
	}
	res, err := repo.InsertUser(user)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	fmt.Printf("Inserted %d row(s)\n", rows)
	return nil
}

func handleDelete(repo repository.IRepositoryUser, idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	res, err := repo.DeleteUserById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("User not found")
			return nil
		}
		return err
	}
	rows, _ := res.RowsAffected()
	fmt.Printf("Deleted %d row(s)\n", rows)
	return nil
}

func handleUpdate(repo repository.IRepositoryUser, idStr, name, email, password string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	user, err := repo.GetUserById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("User not found")
			return nil
		}
		return err
	}

	user.Name = name
	user.Email = email
	user.Password = password

	res, err := repo.UpdateUserById(id, user)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	fmt.Printf("Updated %d row(s)\n", rows)
	return nil
}

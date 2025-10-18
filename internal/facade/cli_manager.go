package facade

import (
	"context"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/db"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/middleware"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/repository/imp"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/runner"
)

func RunCliManager(ctx context.Context, host string, port string, user string, dbname string, password string) error {
	con, err := db.Connect(host, port, user, dbname, password)
	if err != nil {
		return err
	}
	err = db.Migrate(con)
	if err != nil {
		return err
	}

	repoUser := imp.NewPostgresRepoUser(con)
	repoTx := middleware.NewTransactionalMiddleware(con, repoUser)
	repoLogging := middleware.NewLoggerMiddleware(imp.NewPostgresRepoLog(con), repoTx)

	go func() {
		defer con.Close()
		runner.Logic(ctx, repoLogging)
	}()
	return nil
}

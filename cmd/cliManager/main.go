package main

import (
	"context"
	"fmt"
	"github.com/BohdanIpy/simpleCLIdbManager/internal/facade"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 6 {
		fmt.Println("Usage: ./a.out <host> <port> <user> <dbname> <password>")
		return
	}
	host := args[1]
	port := args[2]
	user := args[3]
	dbname := args[4]
	password := args[5]

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cancelingSignals := make(chan os.Signal, 1)
	err := facade.RunCliManager(ctx, host, port, user, dbname, password)
	if err != nil {
		fmt.Println(err)
		cancel()
		return
	}
	<-cancelingSignals

}

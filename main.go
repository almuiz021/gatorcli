package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/almuiz021/gatorcli/internal/config"
	"github.com/almuiz021/gatorcli/internal/database"
)

type state struct {
	db         *database.Queries
	cfg        *config.Config
	httpClient http.Client
	baseUrl    string
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{
		cfg: &cfg,
		db:  dbQueries,
		httpClient: http.Client{
			Timeout: 10 * time.Second,
		},
		baseUrl: "https://www.wagslane.dev/index.xml",
	}

	cmds := &commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerAllUsers)
	cmds.register("agg", handlerGetRequests)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)

	cmdLineArgs := os.Args
	if len(cmdLineArgs) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := strings.ToLower(cmdLineArgs[1])
	cmdArgs := cmdLineArgs[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}

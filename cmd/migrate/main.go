package main

import (
	"flag"
	"log"
	"os"
	"reminder-server/internal/initializers"
	"strconv"

	"github.com/pressly/goose/v3"
)

var (
	GOOSE_DRIVER   string
	GOOSE_DBSTRING string
)

func main() {
	initializers.LoadEnv()
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		log.Fatal("Please provide a goose command (up, down, status, create)")
	}

	command := args[0]

	initializers.ConnectDB()
	GOOSE_DRIVER = os.Getenv("GOOSE_DRIVER")
	GOOSE_DBSTRING = initializers.GetDBString()

	db, err := initializers.GetDB().DB()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Failed to close DB: %v", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	// Set the migration directory
	if err := goose.SetDialect(GOOSE_DRIVER); err != nil {
		log.Fatal(err)
	}

	migrationsDir := "migrations"

	switch command {
	case "create":
		if len(args) < 2 {
			log.Fatal("Please provide migration name")
		}
		name := args[1]
		if err := goose.Create(db, migrationsDir, name, "sql"); err != nil {
			log.Fatalf("goose create: %v", err)
		}
	case "up":
		if err := goose.Up(db, migrationsDir); err != nil {
			log.Fatalf("goose up: %v", err)
		}
	case "down":
		if err := goose.Down(db, migrationsDir); err != nil {
			log.Fatalf("goose down: %v", err)
		}
	case "status":
		if err := goose.Status(db, migrationsDir); err != nil {
			log.Fatalf("goose status: %v", err)
		}
	case "reset":
		if err := goose.Reset(db, migrationsDir); err != nil {
			log.Fatalf("goose reset: %v", err)
		}
	case "down-to":
		if len(args) < 2 {
			log.Fatal("Please provide migration version")
		}

		version, err := strconv.Atoi(args[1])

		if err != nil {
			log.Fatalf("Failed to convert version to int: %v", err)
		}

		if err := goose.DownTo(db, migrationsDir, int64(version)); err != nil {
			log.Fatalf("goose down-to: %v", err)
		}
	default:
		log.Fatalf("Unknown command: %q", command)
	}
}

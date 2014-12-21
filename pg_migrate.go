package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gopkg.in/alecthomas/kingpin.v1"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	wd, _    = os.Getwd()
	dir      = kingpin.Flag("dir", "Directory where migration files are located. Current working directory will be used by default").Default(wd).String()
	host     = kingpin.Flag("host", "Server address and port").Default("localhost").String()
	port     = kingpin.Flag("port", "Server port").Default("5432").Int()
	dbname   = kingpin.Flag("db", "Database name").Default("postgres").String()
	user     = kingpin.Flag("user", "User").Default("postgres").String()
	password = kingpin.Flag("password", "Password").String()
	sslmode  = kingpin.Flag("sslmode", "").Default("disable").String()
	verbose  = kingpin.Flag("verbose", "Verbose output").Default("false").Bool()
)

func createMigrationTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name TEXT,
			time TIMESTAMP
		)`)
	if err != nil {
		return err
	}

	return nil
}

func findLatestMigration(db *sql.DB) (string, time.Time, error) {
	row := db.QueryRow("SELECT name, time FROM migrations ORDER BY time DESC LIMIT 1")

	var name string
	var migrationTime time.Time

	err := row.Scan(&name, &migrationTime)
	if err != nil {
		return "", time.Time{}, err
	}

	return name, migrationTime, nil
}

func migrate(db *sql.DB, name string) error {
	data, err := ioutil.ReadFile(*dir + "/" + name)
	if err != nil {
		panic(err)
	}

	sql := fmt.Sprintf("BEGIN;\n%v\nCOMMIT;\n", string(data))
	if *verbose {
		fmt.Println(sql)
	}

	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	db.Exec("INSERT INTO migrations (name, time) VALUES ($1, $2)", name, time.Now())

	return nil
}

func main() {
	kingpin.Parse()
	fmt.Println(*dir, *host, *user, *password)

	// docker run --rm=true -p 5432:5432 postgres:9.4
	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%v password=%v dbname=%v host=%v port=%v sslmode=%v",
		*user, *password, *dbname, *host, *port, *sslmode,
	))
	if err != nil {
		log.Fatal(err)
	}

	createMigrationTable(db)

	name, migrationTime, err := findLatestMigration(db)
	if err != nil && err.Error() != "sql: no rows in result set" {
		panic(err)
	}
	log.Printf("Latest migration: %v (migrated %v)\n", name, migrationTime)

	files, err := ioutil.ReadDir(*dir)
	if err != nil {
		panic(err)
	}

	existing := strings.SplitN(name, "-", 2)[0]
	existingNum, _ := strconv.ParseInt(existing, 10, 64)

	migrations := []string{}
	for _, file := range files {
		name := file.Name()
		if name[len(name)-4:] != ".sql" {
			continue
		}

		migration := strings.SplitN(name, "-", 2)[0]
		migrationNum, err := strconv.ParseInt(migration, 10, 64)
		if err != nil {
			panic(err)
		}

		if migrationNum <= existingNum {
			continue
		}

		migrations = append(migrations, file.Name())
		fmt.Println(file.Name())
	}

	if len(migrations) == 0 {
		panic("No migrations found")
	}

	sort.Strings(migrations)

	for _, migration := range migrations {
		err := migrate(db, migration)
		if err != nil {
			panic(err)
		}
		if *verbose {
			log.Printf("Migration %v successfully applied", migration)
		}
	}

	fmt.Println(migrations)
}

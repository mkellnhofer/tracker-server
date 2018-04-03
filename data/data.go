package data

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"kellnhofer.com/tracker/constant"
)

const curDbVers = 1

// --- Public methods ---

func GetDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./data/data.db")
	if err != nil {
		log.Fatalf("Could not open database connection! (Error: %s)", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not open database connection! (Error: %s)", err)
	}

	enableDbForeignKeys(db)

	return db
}

func UpdateDb(db *sql.DB) {
	dbVers := getDbVersion(db)
	if dbVers == 0 {
		log.Println("Creating database ...")
		createDb(db)
		updateDb(db, 1)
		log.Println("Successfully created database.")
	} else if dbVers < curDbVers {
		log.Println("Updating database ...")
		updateDb(db, dbVers)
		log.Println("Successfully updated database.")
	}
}

func ParseTime(timeIn string) time.Time {
	timeOut, err := time.ParseInLocation(constant.DbDateFormat, timeIn, time.UTC)
	if err != nil {
		return time.Time{}
	}
	return timeOut
}

func FormatTime(timeIn time.Time) string {
	return timeIn.UTC().Format(constant.DbDateFormat)
}

// --- Private methods ---

func enableDbForeignKeys(db *sql.DB) {
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatalf("Could not enable database foreign keys! (Error: %s)", err)
	}
}

func getDbVersion(db *sql.DB) int {
	var name string
	row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='setting'")
	err := row.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		log.Fatalf("Could not query database version! (Error: %s)", err)
	}

	var version int
	row = db.QueryRow("SELECT value FROM setting WHERE key LIKE 'db_version'")
	err = row.Scan(&version)
	if err != nil {
		log.Fatalf("Could not query database version! (Error: %s)", err)
	}

	return version
}

func createDb(db *sql.DB) {
	dbStmts := readDbFile("db_create.sql")
	executeDbStmts(db, dbStmts)
}

func updateDb(db *sql.DB, dbVers int) {
	for i := dbVers + 1; i <= curDbVers; i++ {
		fileName := fmt.Sprintf("db_update_v%d.sql", i)
		dbStmts := readDbFile(fileName)
		log.Printf("Executing database update v%d ...", i)
		executeDbStmts(db, dbStmts)
		updateDbVersion(db, i)
	}
}

func readDbFile(name string) []string {
	// Open file
	file, err := os.Open("scripts/db/" + name)
	if err != nil {
		log.Fatalf("Could not open database update script %s! (Error: %s)", name, err)
	}
	defer file.Close()

	var stmts []string

	// Read statements from file
	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()

		// If line is empty: Skip
		if line == "" {
			continue
		}

		// Add line to current statement
		buffer.WriteString(line)

		// If line is end of statement: Add current statement to result
		if strings.HasSuffix(line, ";") {
			stmts = append(stmts, buffer.String())
			buffer.Reset()
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Could not read database update script %s! (Error: %s)", name, err)
	}

	return stmts
}

func executeDbStmts(db *sql.DB, stmts []string) {
	for _, stmt := range stmts {
		executeDbStmt(db, stmt)
	}
}

func executeDbStmt(db *sql.DB, stmt string) {
	_, err := db.Exec(stmt)
	if err != nil {
		log.Fatalf("Could not execute database statement! (Error: %s)", err)
	}
}

func updateDbVersion(db *sql.DB, dbVers int) {
	_, err := db.Exec("UPDATE setting SET value = ? WHERE key = 'db_version'", dbVers)
	if err != nil {
		log.Fatalf("Could not update database version! (Error: %s)", err)
	}
}

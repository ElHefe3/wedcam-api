package db

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
    var err error
    DB, err = sql.Open("sqlite3", "./wedcam.db")
    if err != nil {
        return err
    }
	
    schema := `
    CREATE TABLE IF NOT EXISTS accounts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        token TEXT UNIQUE,
        active BOOLEAN DEFAULT true,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS qr_codes (
        token TEXT PRIMARY KEY,
        account_id INTEGER,
        uploads_allowed INTEGER DEFAULT 1,
        uploads_used INTEGER DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(account_id) REFERENCES accounts(id)
    );`

    _, err = DB.Exec(schema)
    return err
}
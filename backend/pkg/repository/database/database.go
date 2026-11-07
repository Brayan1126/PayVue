package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("error creating tables: %w", err)
	}

	// Ejecutar migración para añadir user_id si no existe
	if err := migrateUserID(db); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS debts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL DEFAULT 0,
		name TEXT NOT NULL,
		total_amount REAL NOT NULL,
		remaining_amount REAL NOT NULL,
		due_date DATETIME NOT NULL,
		interest_rate REAL NOT NULL,
		num_installments INTEGER NOT NULL,
		installment_amount REAL NOT NULL,
		payment_day INTEGER NOT NULL,
		paid BOOLEAN DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS incomes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL DEFAULT 0,
		amount REAL NOT NULL,
		source TEXT NOT NULL,
		date DATETIME NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS payments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL DEFAULT 0,
		amount REAL NOT NULL,
		debt_id INTEGER NOT NULL,
		receipt_filename TEXT,
		date DATETIME NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (debt_id) REFERENCES debts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_debts_user_id ON debts(user_id);
	CREATE INDEX IF NOT EXISTS idx_incomes_user_id ON incomes(user_id);
	CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);
	`

	_, err := db.Exec(schema)
	return err
}

// Migración para añadir user_id a tablas existentes
func migrateUserID(db *sql.DB) error {
	// Verificar si la columna user_id existe en debts
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('debts') WHERE name='user_id'").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		// Añadir columna user_id a tablas existentes
		migrations := []string{
			"ALTER TABLE debts ADD COLUMN user_id INTEGER NOT NULL DEFAULT 0",
			"ALTER TABLE incomes ADD COLUMN user_id INTEGER NOT NULL DEFAULT 0",
			"ALTER TABLE payments ADD COLUMN user_id INTEGER NOT NULL DEFAULT 0",
		}
		for _, m := range migrations {
			_, err := db.Exec(m)
			if err != nil {
				log.Printf("Migration step skipped (may already exist): %v", err)
			}
		}
	}

	return nil
}

package main

import (
	"database/sql"

	// Driver para SQLite3
	_ "github.com/mattn/go-sqlite3"
)

// Puntero a la estructura DB, nos permite manejar la
// base de datos
var db *sql.DB

func GetConnection() *sql.DB {
	// Para evitar realizar una nueva conexión en cada llamada a
	// la función GetConnection.
	if db != nil {
		return db
	}
	// Declaramos la variable err para poder usar el operador
	// de asignación “=” en lugar que el de asignación corta,
	// para evitar que cree una nueva variable db en este scope y
	// en su lugar que inicialice la variable db que declaramos a
	// nivel de paquete.
	var err error
	// Conexión a la base de datos
	db, err = sql.Open("sqlite3", "data.sqlite")
	if err != nil {
		panic(err)
	}
	return db
}

func MakeMigrations() error {
	db := GetConnection()
	q := `CREATE TABLE IF NOT EXISTS blocks (
				hash VARCHAR(64),
				size VARCHAR(64),
				stripped_size VARCHAR(64),
				weight VARCHAR(64),
				number VARCHAR(64),
				version VARCHAR(64),
				merkle_root VARCHAR(64),
				timestamp DATETIME,
				nonce VARCHAR(64) PRIMARY KEY,
				bits VARCHAR(64),
				coinbase_param VARCHAR(64),
				transaction_count VARCHAR(64)
			);`
	_, err := db.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

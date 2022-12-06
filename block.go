package main

import (
	"errors"
	"time"
)

type Block struct {
	Hash         string    `json:"hash" bigquery:"hash"`
	Size         int       `json:"size" bigquery:"size"`
	StrippedSize int       `json:"stripped_size" bigquery:"stripped_size"`
	Weight       int       `json:"weight" bigquery:"weight"`
	Number       int       `json:"number" bigquery:"number"`
	Version      int       `json:"version" bigquery:"version"`
	MerkleRoot   string    `json:"merkle_root" bigquery:"merkle_root"`
	Timestamp    time.Time `json:"timestamp" bigquery:"timestamp"`
	// TimestampMonth   Date type       `json:"timestamp_month" bigquery:"timestamp_month"`
	Nonce            string `json:"nonce" bigquery:"nonce"`
	Bits             string `json:"bits" bigquery:"bits"`
	CoinbaseParam    string `json:"coinbase_param" bigquery:"coinbase_param"`
	TransactionCount int    `json:"transaction_count" bigquery:"transaction_count"`
}

func (b Block) Create() error {
	db := GetConnection()
	q := `SELECT
            *
            FROM blocks WHERE nonce=?`
	// Ejecutamos la query
	row, err := db.Query(q, b.Nonce)
	if err != nil {
		return err
	}
	defer row.Close()
	if row.Next() {
		return nil
	}
	// Query para insertar los datos en la tabla notes
	exec := `INSERT INTO blocks (
		hash,
		size,
		stripped_size,
		weight,
		number,
		version,
		merkle_root,
		timestamp,
		nonce,
		bits,
		coinbase_param,
		transaction_count
		)VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	// Preparamos la petición para insertar los datos de manera
	// segura
	stmt, err := db.Prepare(exec)
	if err != nil {
		return err
	}
	// Nos aseguramos de cerrar el recurso antes de finalizar la
	// función gracias a defer.
	defer stmt.Close()
	// Ejecutamos la petición pasando los datos correspondientes.
	// El orden es importante, corresponde con los “?” del
	// string q.
	r, err := stmt.Exec(b.Hash, b.Size, b.StrippedSize, b.Weight, b.Number, b.Version, b.MerkleRoot,
		b.Timestamp, b.Nonce, b.Bits, b.CoinbaseParam, b.TransactionCount)
	if err != nil {
		return err
	}
	// Confirmamos que una fila fuera afectada, debido a que
	// insertamos un registro en la tabla. En caso contrario
	// devolvemos un error.
	if i, err := r.RowsAffected(); err != nil || i != 1 {
		return errors.New("ERROR: Se esperaba una fila afectada")
	}
	// Si llegamos a este punto consideramos que todo el proceso
	// fue exitoso y retornamos un nil para confirmar que no
	// existe un error.
	return nil
}

func (b *Block) GetAll() ([]Block, error) {
	db := GetConnection()
	q := `SELECT
            *
            FROM blocks`
	// Ejecutamos la query
	rows, err := db.Query(q)
	if err != nil {
		return []Block{}, err
	}
	// Cerramos el recurso
	defer rows.Close()
	// Declaramos un slice de bloques para que almacene las
	// bloques que retorna la petición.
	blocks := []Block{}
	// El método Next retorna un bool, mientras sea true indicará
	// que existe un valor siguiente para leer.
	for rows.Next() {
		// Escaneamos el valor actual de la fila e insertamos el
		// retorno en los correspondientes campos de la nota.
		if err := rows.Scan(
			&b.Hash, &b.Size, &b.StrippedSize, &b.Weight, &b.Number,
			&b.Version, &b.MerkleRoot, &b.Timestamp,
			&b.Nonce, &b.Bits, &b.CoinbaseParam, &b.TransactionCount,
		); err != nil {
			return nil, err
		}
		// Añadimos cada nueva nota al slice de bloques que
		// declaramos antes.
		blocks = append(blocks, *b)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

func (b *Block) GetOne(nonce string) (Block, error) {
	db := GetConnection()
	q := `SELECT * FROM blocks WHERE nonce=?`
	row, err := db.Query(q, nonce)
	if err != nil {
		return Block{}, err
	}
	defer row.Close()
	for row.Next() {
		if err := row.Scan(
			&b.Hash, &b.Size, &b.StrippedSize, &b.Weight, &b.Number,
			&b.Version, &b.MerkleRoot, &b.Timestamp,
			&b.Nonce, &b.Bits, &b.CoinbaseParam, &b.TransactionCount,
		); err != nil {
			return Block{}, err
		}
	}
	err = row.Err()
	if err != nil {
		return Block{}, err
	}

	return *b, nil
}

func (b Block) Delete(nonce string) error {
	db := GetConnection()
	q := `DELETE FROM blocks
            WHERE nonce=?`
	stmt, err := db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	r, err := stmt.Exec(nonce)
	if err != nil {
		return err
	}
	if i, err := r.RowsAffected(); err != nil || i != 1 {
		return errors.New("ERROR: Se esperaba una fila afectada")
	}
	return nil
}

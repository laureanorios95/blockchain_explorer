package main

import (
	"errors"
	"time"
)

type Block struct {
	Hash      string    `json:"hash"`
	PrevBlock string    `json:"prev_block"`
	Time      int       `json:"time"`
	Bits      int       `json:"bits"`
	Fee       int       `json:"fee"`
	Nonce     int       `json:"nonce"`
	NTx       int       `json:"n_tx"`
	CreatedAt time.Time `json:"created_at"`
}

func (b Block) Create() error {
	// Obtenemos la conexión a la base de datos.
	db := GetConnection()
	q := `SELECT
            *
            FROM blocks WHERE nonce=?`
	// Ejecutamos la query
	_, err := db.Query(q, b.Nonce)
	if err != nil {
		return err
	}

	// Query para insertar los datos en la tabla notes
	exec := `INSERT INTO blocks (
		hash,
		prev_block,
		time,
		bits,
		fee,
		nonce,
		n_tx,
		created_at
		)VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
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
	r, err := stmt.Exec(b.Hash, b.PrevBlock, b.Time, b.Bits, b.Fee, b.Nonce, b.NTx, time.Now())
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
		rows.Scan(
			&b.Hash,
			&b.PrevBlock,
			&b.Time,
			&b.Bits,
			&b.Fee,
			&b.Nonce,
			&b.NTx,
			&b.CreatedAt,
		)
		// Añadimos cada nueva nota al slice de bloques que
		// declaramos antes.
		blocks = append(blocks, *b)
	}
	return blocks, nil
}

func (b *Block) GetOne(nonce int) (Block, error) {
	db := GetConnection()
	q := `SELECT * FROM blocks WHERE nonce=?`
	row, err := db.Query(q, nonce)
	if err != nil {
		return Block{}, err
	}
	defer row.Close()
	for row.Next() {
		row.Scan(
			&b.Hash,
			&b.PrevBlock,
			&b.Time,
			&b.Bits,
			&b.Fee,
			&b.Nonce,
			&b.NTx,
			&b.CreatedAt,
		)
	}
	return *b, nil
}

func (b Block) Delete(nonce int) error {
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

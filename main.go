// make a functional CRUD with SQLite
// all requests will be get type
// get the latest block from blockchain explorer API https://www.blockchain.com/explorer/api and save it
// get SQLite database with all the blocks saved until the moment refresh it weekly or some period
// get specific block saved in database
// continue with the next steps
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
)

func main() {

	// flag para realizar la creación de las tablas en la base
	// de datos.
	migrate := flag.Bool(
		"migrate", false, "Crea las tablas en la base de datos",
	)
	// Parseando todas las flags
	flag.Parse()
	if *migrate {
		if err := MakeMigrations(); err != nil {
			log.Fatal(err)
		}
	}

	// Instancia de http.DefaultServerMux
	mux := http.NewServeMux()
	// Ruta a manejar
	mux.HandleFunc("/blocks", BlocksHandler)
	// Servidor escuchando en el puerto 8080
	http.ListenAndServe(":8080", mux)
}

// BlocksHandler nos permite manejar la petición a la ruta ‘/blocks // y pasa el control a la función correspondiente según el método
// de la petición.
func BlocksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Get("nonce") != "" {
			GetBlock(w, r)
		} else {
			GetBlocks(w, r)
		}
	case http.MethodPost:
		AddBlock(w, r)
	case http.MethodDelete:
		DeleteBlock(w, r)
	default:
		// Caso por defecto en caso de que se realice una
		// petición con un método diferente a los esperados.
		http.Error(w, "Metodo no permitido",
			http.StatusBadRequest)
		return
	}
}

func GetBlock(w http.ResponseWriter, r *http.Request) {
	// obtenemos el valor pasado en la url como query
	// correspondiente a id, del tipo /blocks?nonce=0123456789.
	nonceStr := r.URL.Query().Get("nonce")
	// Convertimos el valor obtenido del query a un int, de ser
	// posible.
	nonce, err := strconv.Atoi(nonceStr)
	if err != nil {
		http.Error(w, "Query nonce debe ser un número",
			http.StatusBadRequest)
		return
	}
	var blocks Block
	// Solicitando el bloque descrito en el parametro URL
	block, err := blocks.GetOne(nonce)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Convirtiendo el slice de bloques a formato JSON,
	// retorna un []byte y un error.
	j, err := json.Marshal(block)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Escribiendo el código de respuesta.
	w.WriteHeader(http.StatusOK)
	// Estableciendo el tipo de contenido del cuerpo de la
	// respuesta.
	w.Header().Set("Content-Type", "application/json")
	// Escribiendo la respuesta, es decir nuestro slice de bloques
	// en formato JSON.
	w.Write(j)
}

func GetBlocks(w http.ResponseWriter, r *http.Request) {
	// Puntero a una estructura de tipo Block.
	b := new(Block)
	// Solicitando todos los bloques en la base de datos.
	blocks, err := b.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Convirtiendo el slice de bloques a formato JSON,
	// retorna un []byte y un error.
	j, err := json.Marshal(blocks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Escribiendo el código de respuesta.
	w.WriteHeader(http.StatusOK)
	// Estableciendo el tipo de contenido del cuerpo de la
	// respuesta.
	w.Header().Set("Content-Type", "application/json")
	// Escribiendo la respuesta, es decir nuestro slice de bloques
	// en formato JSON.
	w.Write(j)
}

func AddBlock(w http.ResponseWriter, r *http.Request) {
	var (
		block, hash Block
	)
	// Tomando el cuerpo de la petición, en formato JSON, y
	// decodificándola e la variable block que acabamos de
	// declarar.
	resp, err := http.Get("https://blockchain.info/latestblock")
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp.Body.Close()

	res, err := http.Get("https://blockchain.info/rawblock/" + hash.Hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	err = json.NewDecoder(res.Body).Decode(&block)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Creamos la nueva nota gracias al método Create.
	err = block.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteBlock(w http.ResponseWriter, r *http.Request) {
	// obtenemos el valor pasado en la url como query
	// correspondiente a id, del tipo /notes?id=3.
	nonceStr := r.URL.Query().Get("nonce")
	// Verificamos que no esté vacío.
	if nonceStr == "" {
		http.Error(w, "Query nonce es requerido",
			http.StatusBadRequest)
		return
	}
	// Convertimos el valor obtenido del query a un int, de ser
	// posible.
	nonce, err := strconv.Atoi(nonceStr)
	if err != nil {
		http.Error(w, "Query nonce debe ser un número",
			http.StatusBadRequest)
		return
	}
	var block Block
	// Borramos la nota con el id correspondiente.
	err = block.Delete(nonce)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

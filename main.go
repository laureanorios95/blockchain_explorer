package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"
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
		VerifyBlocks(w, r)
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
	nonce := r.URL.Query().Get("nonce")

	var blocks Block
	// Solicitando el bloque descrito en el parametro URL
	block, err := blocks.GetOne(nonce)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	j, err := json.Marshal(block)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Escribiendo el código de respuesta.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Escribiendo la respuesta, es decir nuestro slice de bloques
	// en formato JSON.
	w.Write(j)
}

func GetBlocks(w http.ResponseWriter, r *http.Request) {
	b := new(Block)
	// Solicitando todos los bloques en la base de datos.
	blocks, err := b.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	j, err := json.Marshal(blocks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Escribiendo la respuesta, es decir nuestro slice de bloques
	// en formato JSON.
	w.Write(j)
}

func VerifyBlocks(w http.ResponseWriter, r *http.Request) {
	b := new(Block)
	type Period struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}
	var period Period

	if err := json.NewDecoder(r.Body).Decode(&period); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if period.End.Sub(period.Start) > time.Duration(2*time.Hour) {
		http.Error(w, "Time period is longer than 2 hours", http.StatusBadRequest)
		return
	}

	mainCount, err := mainCount([]time.Time{period.Start, period.End})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	selfCount, err := b.selfCount([]time.Time{period.Start, period.End})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if selfCount == mainCount {
		blocks, err := b.GetMany([]time.Time{period.Start, period.End})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(blocks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
		return
	}
	blocks, err := ExecPipeline([]time.Time{period.Start, period.End})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, block := range blocks {
		if err := block.Create(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	j, err := json.Marshal(blocks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Escribiendo la respuesta, es decir nuestro slice de bloques
	// en formato JSON.
	w.Write(j)
}

func DeleteBlock(w http.ResponseWriter, r *http.Request) {
	// obtenemos el valor pasado en la url como query
	// correspondiente a id, del tipo /blocks?nonce=123456789.
	nonce := r.URL.Query().Get("nonce")
	if nonce == "" {
		http.Error(w, "Query nonce es requerido",
			http.StatusBadRequest)
		return
	}

	var block Block
	// Borramos la nota con el id correspondiente.
	err := block.Delete(nonce)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

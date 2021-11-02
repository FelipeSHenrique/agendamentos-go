package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"module/config"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", config.GetDSN())
	if err != nil {
		log.Fatalf("erro ao abrir banco de dados: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("erro ao conectar: %v", err)
	}

	http.Handle("/especialidade", http.HandlerFunc(index))
	http.Handle("/especialidade-create", http.HandlerFunc(create))

	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("não foi possivel iniciar o servidor: %v", err)
	}
}

type specialties struct {
	ID        int
	Name      string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func index(rw http.ResponseWriter, r *http.Request) {
	result := []specialties{} // Retorna array vazia
	rows, err := db.Query("SELECT * FROM specialties")

	if err != nil {
		log.Printf("não foi encontrado nenhum dado: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var r specialties
		err := rows.Scan(&r.ID, &r.Name, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			log.Printf("não foi encontrado nenhum dado: %v", err)
			return
		}
		result = append(result, r)
	}

	err = json.NewEncoder(rw).Encode(result)
	if err != nil {
		log.Printf("não foi possivel fazer o encode: %v", err)
	}
}

type resultOK struct {
	OK bool
}

type specialtiesCreateRequest struct {
	Name string
}

func create(rw http.ResponseWriter, r *http.Request) {
	req := specialtiesCreateRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("não foi possivel recuperar os dados: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("bad request"))
		return
	}
	if req.Name == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("bad request"))
		return
	}
	_, err = db.Exec("INSERT INTO specialties (name, created_at, updated_at) VALUES($1, NOW(), NOW())", req.Name)
	if err != nil {
		log.Printf("não foi possivel executar a query: %v", err)
		return
	}

	err = json.NewEncoder(rw).Encode(resultOK{OK: true})
	if err != nil {
		log.Printf("não foi possivel fazer o encode: %v", err)
	}
}

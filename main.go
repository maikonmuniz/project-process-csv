package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/upload-csv", uploadCSVHandler)
	fmt.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadCSVHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(32 << 20)

	file, _, err := r.FormFile("arquivo")
	if err != nil {
		http.Error(w, "Erro ao obter o arquivo", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	reader.ReuseRecord = true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Erro ao ler o arquivo CSV", http.StatusInternalServerError)
			return
		}

		fmt.Println(record)
	}

	fmt.Fprintln(w, "Arquivo processado com sucesso")
}

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
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

	lines := make(chan []string)
	errChan := make(chan error)
	var wg sync.WaitGroup

	go func() {
		defer close(lines)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errChan <- err
				return
			}

			lines <- record
		}
	}()

	const numWorkers = 10000
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lines {
				fmt.Println(line)
			}
		}()
	}

	go func() {
		for e := range errChan {
			fmt.Println("Erro ao processar arquivo CSV:", e)
			http.Error(w, "Erro ao processar arquivo CSV", http.StatusInternalServerError)
		}
	}()

	wg.Wait()

	fmt.Fprintln(w, "Arquivo processado com sucesso")
}

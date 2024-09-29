package main

import (
	"fmt"
	"log"
	"net/http"
	"wpp-integration/controllers"

	"github.com/go-chi/chi"
)

func main() {

	r := chi.NewRouter()

	r.Post("/api/consultas", controllers.HandleConsulta)

	fmt.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8000", r))
}

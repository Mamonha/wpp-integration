package main

import (
	"fmt"
	"log"
	"net/http"
	"wpp-integration/controllers"
	"wpp-integration/services"

	"github.com/go-chi/chi"
)

func main() {

	dsn := "root:root@tcp(127.0.0.1:3306)/clinifono?charset=utf8mb4&parseTime=True&loc=Local"
	if err := services.InitDB(dsn); err != nil {
		log.Fatalf("Erro ao inicializar o banco de dados: %v", err)
	}

	r := chi.NewRouter()

	r.Post("/api/consultas", controllers.HandleConsulta)

	r.Post("/api/webhook", controllers.Webhook)
	fmt.Println("Servidor rodando na porta 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

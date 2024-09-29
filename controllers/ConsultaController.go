package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wpp-integration/models"
)

func HandleConsulta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var consultaRequest models.ConsultaRequest

	if err := json.NewDecoder(r.Body).Decode(&consultaRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("ConsultaID: %d, PacienteNome: %s, DataAgendamento: %s\n",

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(consultaID)
}

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

	var consultaID models.ConsultaID

	if err := json.NewDecoder(r.Body).Decode(&consultaID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("ConsultaID: %d, PacienteID: %d\n", consultaID.ConsultaID, consultaID.PacienteID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(consultaID)
}

package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"wpp-integration/models"
	"wpp-integration/services"
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
		consultaRequest.ConsultaID, consultaRequest.HoraDeInicio, consultaRequest.PacienteNome, consultaRequest.HoraDoFim, consultaRequest.Telefone)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(consultaRequest)

	services.SendMessage(consultaRequest)
}

func Webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	fmt.Println("webhook" + string(body))

	values, err := url.ParseQuery(string(body))
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusInternalServerError)
		return
	}

	event := values.Get("event")
	messageBody := values.Get("message_body")
	phoneNumber := values.Get("contact_phone_number")

	if event == "message" {
		if err := HandleMessageRecord(phoneNumber, messageBody); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Message status updated successfully")
	} else {
		http.Error(w, "Invalid event", http.StatusBadRequest)
	}
}

func HandleMessageRecord(phoneNumber string, messageBody string) error {
	db := services.GetDB()

	var record models.MessageRecord
	if err := db.Where("contact_phone_number = ? AND status = ?", phoneNumber, models.Pending).First(&record).Error; err != nil {
		return fmt.Errorf("Pending message not found")
	}

	switch messageBody {
	case "1":
		fmt.Println("Confirmada")
		record.Status = string(models.Confirmed)

		services.ConfirmedMessage(record)
		services.HandleConsultaStatus(record, "CONFIRMED")

	case "2":
		fmt.Println("Cancelada")
		record.Status = string(models.Canceled)
		services.CanceledMessage(record)
		services.HandleConsultaStatus(record, "CANCELLED")
	default:
		services.RetryMessage(record)
		return fmt.Errorf("Invalid response")
	}

	if err := db.Save(&record).Error; err != nil {
		return fmt.Errorf("Failed to update message status")
	}
	return nil
}

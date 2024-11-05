package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"wpp-integration/models"
	"wpp-integration/services"

	"github.com/joho/godotenv"
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
	fmt.Println("Chegou aqui")
	fmt.Print(consultaRequest)

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

func Saldo(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	endpointUrl := os.Getenv("WHATSGW_BALANCE")
	if endpointUrl == "" {
		http.Error(w, "URL não configurada", http.StatusInternalServerError)
		return
	}

	apiKey := os.Getenv("WHATSGW_API_KEY")
	if apiKey == "" {
		http.Error(w, "Chave de API não configurada", http.StatusInternalServerError)
		return
	}
	data := url.Values{}
	data.Set("apikey", apiKey)

	req, err := http.NewRequest("GET", endpointUrl, nil)
	if err != nil {
		http.Error(w, "Erro ao criar a requisição", http.StatusInternalServerError)
		fmt.Println("Erro ao criar a requisição:", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erro ao fazer a requisição", http.StatusInternalServerError)
		fmt.Println("Erro ao fazer a requisição:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler a resposta", http.StatusInternalServerError)
		fmt.Println("Erro ao ler a resposta:", err)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusInternalServerError)
		fmt.Println("Erro ao decodificar o JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Erro ao enviar resposta JSON", http.StatusInternalServerError)
		fmt.Println("Erro ao enviar resposta JSON:", err)
		return
	}
}

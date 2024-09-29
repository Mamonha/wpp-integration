package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"wpp-integration/models"

	"github.com/joho/godotenv"
)

var ApiKey string
var sendMessageUrl string
var apiURL string
var OriginPhoneNumber string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}
	ApiKey = os.Getenv("WHATSGW_API_KEY")
	sendMessageUrl = os.Getenv("WHATSGW_SEND_MESSAGE")
	apiURL = os.Getenv("CONFIRM_CONSULTA")
	OriginPhoneNumber = os.Getenv("WHATSGW_PHONE_NUMBER")
}

func SendMessage(consulta models.ConsultaRequest) error {
	dataAgendamentoTime, err := time.Parse("2006-01-02", consulta.DataAgendamento)
	if err != nil {
		return fmt.Errorf("erro ao converter DataAgendamento: %v", err)
	}

	dataFormatada := dataAgendamentoTime.Format("02/01/2006")

	message := models.Message{
		APIKey:             ApiKey,
		PhoneNumber:        OriginPhoneNumber,
		ContactPhoneNumber: consulta.Telefone,
		MessageCustomID:    strconv.Itoa(int(consulta.ConsultaID)),
		MessageType:        "text",
		MessageBody: "Olá, " + consulta.PacienteNome + "!\n\n" +
			"Sua consulta está agendada para o dia " + dataFormatada +
			" às " + consulta.HoraDeInicio + ".\n\n" +
			"Por favor, responda com:\n1️⃣ para confirmar a consulta\n2️⃣ para cancelar a consulta\n\nObrigado!",
	}

	if err := sendMessage(sendMessageUrl, message); err != nil {
		return err
	}

	messageRecord := models.MessageRecord{
		ConsultaID:         consulta.ConsultaID,
		ContactPhoneNumber: consulta.Telefone,
		MessageBody:        message.MessageBody,
		MessageCustomID:    message.MessageCustomID,
		Status:             string(models.Pending),
		CreatedAt:          time.Now(),
		PacienteNome:       consulta.PacienteNome,
		DataAgendamento:    consulta.DataAgendamento,
		HoraDeInicio:       consulta.HoraDeInicio,
	}

	if err := SaveMessageToDB(messageRecord); err != nil {
		return fmt.Errorf("erro ao salvar mensagem no banco de dados: %v", err)
	}

	fmt.Println("Mensagem enviada e registrada com sucesso!")
	return nil
}

func ConfirmedMessage(consulta models.MessageRecord) error {
	dataAgendamentoTime, err := time.Parse("2006-01-02", consulta.DataAgendamento)
	if err != nil {
		return fmt.Errorf("erro ao converter DataAgendamento: %v", err)
	}

	dataFormatada := dataAgendamentoTime.Format("02/01/2006")

	message := models.Message{
		APIKey:             ApiKey,
		PhoneNumber:        OriginPhoneNumber,
		ContactPhoneNumber: consulta.ContactPhoneNumber,
		MessageCustomID:    strconv.Itoa(int(consulta.ConsultaID)),
		MessageType:        "text",
		MessageBody: "Olá, " + consulta.PacienteNome + "!\n\n" +
			"Sua consulta está confirmada para o dia " + dataFormatada +
			" às " + consulta.HoraDeInicio + ".\n\n" +
			"Obrigado por confirmar!",
	}

	return sendMessage(sendMessageUrl, message)
}

func CanceledMessage(consulta models.MessageRecord) error {
	dataAgendamentoTime, err := time.Parse("2006-01-02", consulta.DataAgendamento)
	if err != nil {
		return fmt.Errorf("erro ao converter DataAgendamento: %v", err)
	}

	dataFormatada := dataAgendamentoTime.Format("02/01/2006")

	message := models.Message{
		APIKey:             ApiKey,
		PhoneNumber:        OriginPhoneNumber,
		ContactPhoneNumber: consulta.ContactPhoneNumber,
		MessageCustomID:    strconv.Itoa(int(consulta.ConsultaID)),
		MessageType:        "text",
		MessageBody: "Olá, " + consulta.PacienteNome + "!\n\n" +
			"Sua consulta marcada para o dia " + dataFormatada +
			" às " + consulta.HoraDeInicio + " foi cancelada.\n\n" +
			"Se desejar reagendar, entre em contato diretamente conosco.",
	}
	return sendMessage(sendMessageUrl, message)
}

func RetryMessage(consulta models.MessageRecord) error {
	message := models.Message{
		APIKey:             ApiKey,
		PhoneNumber:        OriginPhoneNumber,
		ContactPhoneNumber: consulta.ContactPhoneNumber,
		MessageCustomID:    strconv.Itoa(int(consulta.ConsultaID)),
		MessageType:        "text",
		MessageBody: "Desculpe, a opção escolhida não é válida.\n\n" +
			"Por favor, responda com:\n1️⃣ para confirmar a consulta\n" +
			"2️⃣ para cancelar a consulta\n\nTente novamente!",
	}
	return sendMessage(sendMessageUrl, message)
}

func sendMessage(url string, message models.Message) error {
	fmt.Println("URL:", url)
	fmt.Println("Mensagem:", message)

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("erro ao converter mensagem para JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao fazer a requisição POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("recebido status não-200: %s", resp.Status)
	}

	return nil
}

func GetPendingMessagesByPhoneNumber(phoneNumber string) ([]models.MessageRecord, error) {
	var messages []models.MessageRecord
	if err := db.Where("phone_number = ? AND status = ?", phoneNumber, models.Pending).Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar mensagens pendentes: %v", err)
	}
	return messages, nil
}

func UpdateMessageStatus(message models.MessageRecord) error {
	if err := db.Save(&message).Error; err != nil {
		return fmt.Errorf("erro ao atualizar status da mensagem: %v", err)
	}
	return nil
}

func HandleConsultaStatus(record models.MessageRecord, status string) {
	consultaID := record.ConsultaID
	url := apiURL + fmt.Sprintf("%d", consultaID)

	payload := map[string]string{
		"status": status,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Erro ao converter payload para JSON: %v\n", err)
		return
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Erro ao criar requisição: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Erro ao fazer a requisição PUT: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erro ao atualizar status da consulta, status: %s\n", resp.Status)
		return
	}

	fmt.Printf("Consulta %s com sucesso!\n", status)
}
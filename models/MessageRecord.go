package models

import (
	"time"
)

type MessageRecord struct {
	ID                 uint64    `gorm:"primaryKey"`       
	ConsultaID         uint64    `json:"consulta_id"`         
	PhoneNumber        string    `json:"phone_number"`      
	ContactPhoneNumber string    `json:"contact_phone_number"` 
	MessageBody        string    `json:"message_body"`     
	MessageCustomID    string    `json:"message_custom_id"`    
	Status             string    `json:"status"`              
	CreatedAt          time.Time `json:"created_at"`         
	PacienteNome       string    `json:"paciente_nome"`      
	DataAgendamento    string    `json:"data_agendamento"`    
	HoraDeInicio       string    `json:hora_inicio"`
}

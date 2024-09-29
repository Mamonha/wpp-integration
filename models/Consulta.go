package models

import "time"

type Consulta struct {
	ID              uint64         `json:"id"`
	DataAgendamento time.Time      `json:"dataAgendamento"`
	HoraDeInicio    string         `json:"horaDeInicio"`
	HoraDoFim       string         `json:"horaDoFim"`
	Descricao       string         `json:"descricao"`
	Status          ConsultaStatus `json:"status"`
	Usuario         Usuario        `json:"usuario"`
	Paciente        Paciente       `json:"paciente"`
}

type Usuario struct {
	ID uint64 `json:"id"`
}

type Paciente struct {
	ID uint64 `json:"id"`
}

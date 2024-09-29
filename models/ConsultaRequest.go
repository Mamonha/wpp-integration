package models

type ConsultaRequest struct {
	ConsultaID      uint64 `json:"consultaId"`
	PacienteNome    string `json:"pacienteNome"`
	DataAgendamento string `json:"dataAgendamento"`
	HoraDeInicio    string `json:"horaDeInicio"`
	HoraDoFim       string `json:"horaDoFim"`
	Descricao       string `json:"descricao"`
	Telefone        string `json:"telefone"`
}
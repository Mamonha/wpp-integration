package models

type ConsultaStatus string

const (
	Pending   ConsultaStatus = "PENDING"
	Confirmed ConsultaStatus = "CONFIRMED"
	Canceled  ConsultaStatus = "CANCELED"
)

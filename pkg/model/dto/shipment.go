package dto

type Shipment struct {
	Id                string    `json:"id"`
	VaccineName       string    `json:"VaccineName"`
	Quantity          uint64    `json:"quantity"`
	ManufacturingDate int64     `json:"ManufacturingDate"`
	ManufacturerId    string    `json:"ManufacturerId"`
	ExpirationDate    int64     `json:"ExpirationDate"`
	AuthorityId       string    `json:"AuthorityId"`
	CustomerId        string    `json:"CustomerId"`
	Events            []Event   `json:"events"`
	MessagesSent      []Message `json:"MessagesSent"`
	MessagesReceived  []Message `json:"MessagesReceived"`
}

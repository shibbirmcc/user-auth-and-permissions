package services

type PasswordDeliveryType string

const (
	POSTGRESQL  PasswordDeliveryType = "POSTGRESQL"
	REDIS       PasswordDeliveryType = "REDIS"
	KAFKA_TOPIC PasswordDeliveryType = "KAFKA_TOPIC"
)

func (pst PasswordDeliveryType) String() string {
	switch pst {
	case POSTGRESQL:
		return "POSTGRESQL"
	case REDIS:
		return "REDIS"
	case KAFKA_TOPIC:
		return "KAFKA_TOPIC"
	default:
		return ""
	}
}

type PasswordDeliveryService interface {
	SendPassword(email, firstName, middleName, lastName, password string) error
}

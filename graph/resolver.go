package graph

import (
	"fio-service/iface"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Ps     iface.PersonService
	logger iface.Ilogger
}

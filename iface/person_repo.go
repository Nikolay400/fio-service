package iface

import "fio-service/model"

type PersonRepo interface {
	GetPeopleByFilters(*model.Filters) ([]*model.Person, error)
	GetPersonById(id int) (*model.Person, error)
	DeletePersonById(id int) (*model.Person, error)
	UpdatePerson(*model.UpdatePerson) (*model.Person, error)
	AddPerson(*model.Person) (int, error)
}

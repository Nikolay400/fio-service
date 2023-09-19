package convert

import (
	m "fio-service/graph/model"
	"fio-service/model"
)

func ConvertGqlToModel(gqlNewPerson *m.GqlNewPerson) *model.Person {
	var person model.Person
	person.Name = gqlNewPerson.Name
	person.Surname = gqlNewPerson.Surname

	if gqlNewPerson.Patronymic != nil {
		person.Patronymic = *gqlNewPerson.Patronymic
	}
	if gqlNewPerson.Age != nil {
		person.Age = *gqlNewPerson.Age
	}
	if gqlNewPerson.Country != nil {
		person.Country = *gqlNewPerson.Country
	}
	if gqlNewPerson.Gender != nil {
		person.Gender = *gqlNewPerson.Gender
	}
	return &person
}

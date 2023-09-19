package repo

import (
	"database/sql"
	"fio-service/env"
	"fio-service/iface"
	"fio-service/model"
)

type PersonRepo struct {
	db     *sql.DB
	logger iface.Ilogger
}

func NewPersonRepo(logger iface.Ilogger) (*PersonRepo, error) {
	dsn := env.GetDbDsn()
	conn, err := sql.Open("postgres", dsn)
	return &PersonRepo{conn, logger}, err
}

func (pr *PersonRepo) GetPeopleByFilters(filters *model.Filters) ([]*model.Person, error) {

	if filters.OnPage == 0 {
		filters.OnPage = 5
	}
	if filters.PageNum == 0 {
		filters.PageNum = 1
	}

	querySrting := "" +
		"SELECT id, name, surname, patronymic, age, gender, country " +
		"FROM people " +
		"WHERE ($1=0 OR age>=$1) AND ($2=0 OR age<=$2) " +
		"AND ($3='' OR gender=$3) " +
		"AND ($4='' OR country=$4) " +
		"AND ($5='' OR name ILIKE CONCAT('%',$5,'%') OR surname ILIKE CONCAT('%',$5,'%') OR patronymic ILIKE CONCAT('%',$5,'%')) " +
		"LIMIT $6 OFFSET $7"
	queryRes, err := pr.db.Query(querySrting, filters.AgeFrom, filters.AgeTo, filters.Gender, filters.Country, filters.Search, filters.OnPage, (filters.PageNum-1)*filters.OnPage)
	if err != nil {
		return nil, err
	}
	defer queryRes.Close()

	people := make([]*model.Person, 0)
	for queryRes.Next() {
		var person model.Person
		queryRes.Scan(&person.Id, &person.Name, &person.Surname, &person.Patronymic, &person.Age, &person.Gender, &person.Country)
		people = append(people, &person)
	}

	return people, nil
}

func (pr *PersonRepo) GetPersonById(id int) (*model.Person, error) {
	querySrting := "" +
		"SELECT id, name, surname, patronymic, age, gender, country " +
		"FROM people " +
		"WHERE id = $1"

	queryRes, err := pr.db.Query(querySrting, id)
	if err != nil {
		return nil, err
	}
	defer queryRes.Close()

	var person model.Person
	if queryRes.Next() {
		queryRes.Scan(&person.Id, &person.Name, &person.Surname, &person.Patronymic, &person.Age, &person.Gender, &person.Country)
	}
	return &person, nil
}

func (pr *PersonRepo) DeletePersonById(id int) (*model.Person, error) {
	querySrting := "DELETE FROM people WHERE id = $1 RETURNING id, name, surname, patronymic, age, gender, country"

	queryRes, err := pr.db.Query(querySrting, id)
	if err != nil {
		return nil, err
	}
	defer queryRes.Close()

	var person model.Person
	if queryRes.Next() {
		queryRes.Scan(&person.Id, &person.Name, &person.Surname, &person.Patronymic, &person.Age, &person.Gender, &person.Country)
	}
	return &person, nil
}

func (pr *PersonRepo) UpdatePerson(person *model.UpdatePerson) (*model.Person, error) {
	querySrting := "UPDATE people SET name=COALESCE($1,name), surname=COALESCE($2,surname), patronymic=COALESCE($3, patronymic), age=COALESCE($4, age), gender=COALESCE($5, gender), country=COALESCE($6, country) WHERE id = $7 RETURNING id, name, surname, patronymic, age, gender, country"

	queryRes, err := pr.db.Query(querySrting, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Country, person.Id)
	if err != nil {
		return nil, err
	}
	defer queryRes.Close()

	var returnPerson model.Person
	if queryRes.Next() {
		queryRes.Scan(&returnPerson.Id, &returnPerson.Name, &returnPerson.Surname, &returnPerson.Patronymic, &returnPerson.Age, &returnPerson.Gender, &returnPerson.Country)
	}
	return &returnPerson, nil
}

func (pr *PersonRepo) AddPerson(person *model.Person) (int, error) {
	querySrting := "INSERT INTO people(name, surname, patronymic, age, gender, country) VALUES($1, $2, $3, $4, $5, $6) RETURNING id"

	personId, err := pr.db.Query(querySrting, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Country)
	if err != nil {
		return 0, err
	}
	defer personId.Close()

	var id int
	if personId.Next() {
		personId.Scan(&id)
	}
	return id, nil
}

func (pr *PersonRepo) Close() {
	pr.db.Close()
}

package service

import (
	"encoding/json"
	"fio-service/iface"
	"fio-service/model"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type PersonService struct {
	repo        iface.PersonRepo
	redisClient iface.Icacher
	logger      iface.Ilogger
}

func NewPersonService(repo iface.PersonRepo, redisClient iface.Icacher, logger iface.Ilogger) *PersonService {
	return &PersonService{repo, redisClient, logger}
}

func (ps *PersonService) GetPeopleByFilters(filters *model.Filters) ([]*model.Person, error) {
	var people []*model.Person
	redisKey := filters.GetKeyForRedis()
	peopleStr, redisErr := ps.redisClient.Get(redisKey)
	if redisErr == redis.Nil {
		people, err := ps.repo.GetPeopleByFilters(filters)
		if err != nil {
			return nil, err
		}
		jsonStr, jsonErr := json.Marshal(people)
		if jsonErr != nil {
			return nil, jsonErr
		}
		err = ps.redisClient.Set(redisKey, jsonStr, time.Hour*24)
		if err != nil {
			return nil, err
		}
		ps.logger.Info("People were taken from DB and saved in cache.")
		return people, nil
	} else if redisErr != nil {
		return nil, redisErr
	} else {
		err := json.Unmarshal([]byte(peopleStr), &people)
		if err != nil {
			return nil, err
		}
	}
	return people, nil
}

func (ps *PersonService) GetPersonById(id int) (*model.Person, error) {
	var person *model.Person
	idStr := strconv.Itoa(id)
	personJson, err := ps.redisClient.Get("person" + idStr)
	if err == redis.Nil {
		person, err = ps.repo.GetPersonById(id)
		if err != nil {
			return nil, err
		}
		err = ps.redisClient.Set("person"+idStr, person, time.Hour*24)
		if err != nil {
			return nil, err
		}
		ps.logger.Info("Person was taken from DB and saved in cache. Person: " + ps.logger.Json(person))
	} else if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal([]byte(personJson), &person)
		if err != nil {
			return nil, err
		}
	}
	return person, nil
}

func (ps *PersonService) DeletePersonById(id int) (*model.Person, error) {
	person, err := ps.repo.DeletePersonById(id)
	if err != nil {
		return nil, err
	}
	ps.logger.Info("Person was deleted from DB. Deleted person: " + ps.logger.Json(person))
	err = ps.redisClient.Del("person" + strconv.Itoa(id))
	return person, err
}

func (ps *PersonService) UpdatePerson(updatePerson *model.UpdatePerson) (*model.Person, error) {
	person, err := ps.repo.UpdatePerson(updatePerson)
	if err != nil {
		return nil, err
	}
	ps.logger.Info("Person was updated. Updated person: " + ps.logger.Json(person))
	err = ps.redisClient.Del("person" + strconv.Itoa(updatePerson.Id))
	return person, err
}

func (ps *PersonService) AddPerson(person *model.Person) (int, error) {
	id, err := ps.repo.AddPerson(person)
	if err != nil {
		return 0, err
	}
	person.Id = id
	ps.logger.Info("New person was added. New person: " + ps.logger.Json(person))
	return id, nil
}

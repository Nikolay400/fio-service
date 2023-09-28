package test

import (
	"encoding/json"
	"errors"
	"fio-service/logger"
	mock_iface "fio-service/mock"
	"fio-service/model"
	"fio-service/service"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
)

type mocks struct {
	repo   *mock_iface.MockPersonRepo
	cacher *mock_iface.MockIcacher
}

func TestGetPeopleByFilters(t *testing.T) {
	type Expected struct {
		people []*model.Person
		isErr  bool
	}

	people := []*model.Person{
		{
			Id:         1,
			Name:       "Ivan",
			Surname:    "Ivanov",
			Patronymic: "Ivanovich",
			Age:        40,
			Gender:     "male",
			Country:    "RU",
		},
		{
			Id:         2,
			Name:       "Pavel",
			Surname:    "Pavlov",
			Patronymic: "Pavlovich",
			Age:        20,
			Gender:     "male",
			Country:    "RU",
		},
	}

	tests := []struct {
		name        string
		prepareMock func(m *mocks, filters *model.Filters)
		arg         *model.Filters
		exp         *Expected
	}{
		{
			name: "test with successful getting people from db",
			prepareMock: func(m *mocks, filters *model.Filters) {
				cacheKey := filters.GetKeyForRedis()
				jsonPeople, _ := json.Marshal(people)
				gomock.InOrder(
					m.cacher.EXPECT().Get(cacheKey).Return("", redis.Nil),
					m.repo.EXPECT().GetPeopleByFilters(filters).Return(people, nil),
					m.cacher.EXPECT().Set(cacheKey, jsonPeople, time.Hour*24).Return(nil),
				)
			},
			arg: &model.Filters{},
			exp: &Expected{
				people: people,
				isErr:  false,
			},
		},
		{
			name: "test with successful getting people from cache",
			prepareMock: func(m *mocks, filters *model.Filters) {
				cacheKey := filters.GetKeyForRedis()
				jsonPeople, _ := json.Marshal(people)
				gomock.InOrder(
					m.cacher.EXPECT().Get(cacheKey).Return(string(jsonPeople), nil),
					m.repo.EXPECT().GetPeopleByFilters(gomock.Any()).Times(0),
				)
			},
			arg: &model.Filters{},
			exp: &Expected{
				people: people,
				isErr:  false,
			},
		},
		{
			name: "test with internal error from cache",
			prepareMock: func(m *mocks, filters *model.Filters) {
				var err error = errors.New("INTERNAL ERROR")
				cacheKey := filters.GetKeyForRedis()
				gomock.InOrder(
					m.cacher.EXPECT().Get(cacheKey).Return("", err),
					m.repo.EXPECT().GetPeopleByFilters(gomock.Any()).Times(0),
				)
			},
			arg: &model.Filters{},
			exp: &Expected{
				people: nil,
				isErr:  true,
			},
		},
		{
			name: "test with internal error from db",
			prepareMock: func(m *mocks, filters *model.Filters) {
				var err error = errors.New("INTERNAL ERROR")
				cacheKey := filters.GetKeyForRedis()
				gomock.InOrder(
					m.cacher.EXPECT().Get(cacheKey).Return("", redis.Nil),
					m.repo.EXPECT().GetPeopleByFilters(filters).Return(nil, err),
					m.cacher.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
				)
			},
			arg: &model.Filters{},
			exp: &Expected{
				people: nil,
				isErr:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockPersonRepo := mock_iface.NewMockPersonRepo(ctrl)
			mockCacher := mock_iface.NewMockIcacher(ctrl)

			m := &mocks{mockPersonRepo, mockCacher}
			tt.prepareMock(m, tt.arg)
			logger, _ := logger.NewZapLogger()
			srvc := service.NewPersonService(mockPersonRepo, mockCacher, logger)
			people, err := srvc.GetPeopleByFilters(tt.arg)
			isError := (err != nil)
			if isError != tt.exp.isErr || !reflect.DeepEqual(people, tt.exp.people) {
				t.Errorf("TEST ERROR: \n"+
					"Got: people = %v isErr = %t \n"+
					"Expected: people = %v, isErr = %t \n",
					people, isError, tt.exp.people, tt.exp.isErr,
				)
			}
		})
	}
}

func TestGetPersonById(t *testing.T) {
	StandardPerson := &model.Person{
		Id:         10,
		Name:       "Ivan",
		Surname:    "Ivanov",
		Patronymic: "Ivanovich",
		Age:        40,
		Gender:     "male",
		Country:    "RU",
	}

	type Expected struct {
		person *model.Person
		isErr  bool
	}

	tests := []struct {
		name        string
		prepareMock func(m *mocks, id int)
		arg         int
		exp         *Expected
	}{
		{
			name: "test with successful getting new person from db",
			prepareMock: func(m *mocks, id int) {
				gomock.InOrder(
					m.cacher.EXPECT().Get("person"+strconv.Itoa(id)).Return("", redis.Nil),
					m.repo.EXPECT().GetPersonById(id).Return(StandardPerson, nil),
					m.cacher.EXPECT().Set("person"+strconv.Itoa(id), StandardPerson, time.Hour*24).Return(nil),
				)
			},
			arg: 10,
			exp: &Expected{
				person: StandardPerson,
				isErr:  false,
			},
		},
		{
			name: "test with successful getting new person from cache",
			prepareMock: func(m *mocks, id int) {
				jsonPerson, _ := json.Marshal(StandardPerson)
				gomock.InOrder(
					m.cacher.EXPECT().Get("person"+strconv.Itoa(id)).Return(string(jsonPerson), nil),
					m.repo.EXPECT().GetPersonById(id).Times(0),
				)
			},
			arg: 10,
			exp: &Expected{
				person: StandardPerson,
				isErr:  false,
			},
		},
		{
			name: "test with internal error from cache",
			prepareMock: func(m *mocks, id int) {
				var err error = errors.New("INTERNAL ERROR")
				gomock.InOrder(
					m.cacher.EXPECT().Get("person"+strconv.Itoa(id)).Return("", err),
					m.repo.EXPECT().GetPersonById(id).Times(0),
				)
			},
			arg: 10,
			exp: &Expected{
				person: nil,
				isErr:  true,
			},
		},
		{
			name: "test with internal error from db",
			prepareMock: func(m *mocks, id int) {
				var err error = errors.New("INTERNAL ERROR")
				gomock.InOrder(
					m.cacher.EXPECT().Get("person"+strconv.Itoa(id)).Return("", redis.Nil),
					m.repo.EXPECT().GetPersonById(id).Return(nil, err),
					m.cacher.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(0),
				)
			},
			arg: 10,
			exp: &Expected{
				person: nil,
				isErr:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockPersonRepo := mock_iface.NewMockPersonRepo(ctrl)
			mockCacher := mock_iface.NewMockIcacher(ctrl)

			m := &mocks{mockPersonRepo, mockCacher}
			tt.prepareMock(m, tt.arg)
			logger, _ := logger.NewZapLogger()
			srvc := service.NewPersonService(mockPersonRepo, mockCacher, logger)
			person, err := srvc.GetPersonById(10)
			isError := (err != nil)
			if isError != tt.exp.isErr || !reflect.DeepEqual(person, tt.exp.person) {
				t.Errorf("TEST ERROR: \n"+
					"Got: person = %v, isErr = %t \n"+
					"Expected: person = %v, isErr = %t \n",
					person, isError, tt.exp.person, tt.exp.isErr,
				)
			}
		})
	}
}

func TestDeletePersonById(t *testing.T) {
	StandardPerson := &model.Person{
		Id:         10,
		Name:       "Ivan",
		Surname:    "Ivanov",
		Patronymic: "Ivanovich",
		Age:        40,
		Gender:     "male",
		Country:    "RU",
	}

	type Expected struct {
		person *model.Person
		isErr  bool
	}

	tests := []struct {
		name        string
		prepareMock func(m *mocks, id int)
		arg         int
		exp         *Expected
	}{
		{
			name: "test with successful deleting",
			prepareMock: func(m *mocks, id int) {
				gomock.InOrder(
					m.repo.EXPECT().DeletePersonById(id).Return(StandardPerson, nil),
					m.cacher.EXPECT().Del("person"+strconv.Itoa(id)).Return(nil),
				)
			},
			arg: 10,
			exp: &Expected{
				person: StandardPerson,
				isErr:  false,
			},
		},
		{
			name: "test with internal error from db",
			prepareMock: func(m *mocks, id int) {
				gomock.InOrder(
					m.repo.EXPECT().DeletePersonById(id).Return(nil, errors.New("INTERNAL ERROR")),
					m.cacher.EXPECT().Del(gomock.Any()).Times(0),
				)
			},
			arg: 10,
			exp: &Expected{
				person: nil,
				isErr:  true,
			},
		},
		{
			name: "test with internal error from cache",
			prepareMock: func(m *mocks, id int) {
				gomock.InOrder(
					m.repo.EXPECT().DeletePersonById(id).Return(StandardPerson, nil),
					m.cacher.EXPECT().Del("person"+strconv.Itoa(id)).Return(errors.New("INTERNAL ERROR")),
				)
			},
			arg: 10,
			exp: &Expected{
				person: StandardPerson,
				isErr:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockPersonRepo := mock_iface.NewMockPersonRepo(ctrl)
			mockCacher := mock_iface.NewMockIcacher(ctrl)

			m := &mocks{mockPersonRepo, mockCacher}
			tt.prepareMock(m, tt.arg)
			logger, _ := logger.NewZapLogger()
			srvc := service.NewPersonService(mockPersonRepo, mockCacher, logger)
			person, err := srvc.DeletePersonById(10)
			isError := (err != nil)
			if isError != tt.exp.isErr || !reflect.DeepEqual(person, tt.exp.person) {
				t.Errorf("TEST ERROR: \n"+
					"Got: person = %v, isErr = %t \n"+
					"Expected: person = %v, isErr = %t \n",
					person, isError, tt.exp.person, tt.exp.isErr,
				)
			}
		})
	}
}

func TestUpdatePerson(t *testing.T) {
	ReturnPerson := &model.Person{
		Id:         10,
		Name:       "Nikolay",
		Surname:    "Ivanov",
		Patronymic: "",
		Age:        45,
		Gender:     "male",
		Country:    "RU",
	}

	name := "Nikolay"
	patronymic := ""
	age := 45

	UpdatePerson := &model.UpdatePerson{
		Id:         10,
		Name:       &name,
		Patronymic: &patronymic,
		Age:        &age,
	}

	type Expected struct {
		person *model.Person
		isErr  bool
	}

	tests := []struct {
		name        string
		prepareMock func(m *mocks, person *model.UpdatePerson)
		arg         *model.UpdatePerson
		exp         *Expected
	}{
		{
			name: "test with successful updating",
			prepareMock: func(m *mocks, person *model.UpdatePerson) {
				gomock.InOrder(
					m.repo.EXPECT().UpdatePerson(person).Return(ReturnPerson, nil),
					m.cacher.EXPECT().Del("person"+strconv.Itoa(person.Id)).Return(nil),
				)
			},
			arg: UpdatePerson,
			exp: &Expected{
				person: ReturnPerson,
				isErr:  false,
			},
		},
		{
			name: "test with internal error from db",
			prepareMock: func(m *mocks, person *model.UpdatePerson) {
				gomock.InOrder(
					m.repo.EXPECT().UpdatePerson(UpdatePerson).Return(nil, errors.New("INTERNAL ERROR")),
					m.cacher.EXPECT().Del(gomock.Any()).Times(0),
				)
			},
			arg: UpdatePerson,
			exp: &Expected{
				person: nil,
				isErr:  true,
			},
		},
		{
			name: "test with internal error from cache",
			prepareMock: func(m *mocks, person *model.UpdatePerson) {
				gomock.InOrder(
					m.repo.EXPECT().UpdatePerson(person).Return(ReturnPerson, nil),
					m.cacher.EXPECT().Del("person"+strconv.Itoa(person.Id)).Return(errors.New("INTERNAL ERROR")),
				)
			},
			arg: UpdatePerson,
			exp: &Expected{
				person: ReturnPerson,
				isErr:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockPersonRepo := mock_iface.NewMockPersonRepo(ctrl)
			mockCacher := mock_iface.NewMockIcacher(ctrl)

			m := &mocks{mockPersonRepo, mockCacher}
			tt.prepareMock(m, tt.arg)
			logger, _ := logger.NewZapLogger()
			srvc := service.NewPersonService(mockPersonRepo, mockCacher, logger)
			person, err := srvc.UpdatePerson(tt.arg)
			isError := (err != nil)
			if isError != tt.exp.isErr || !reflect.DeepEqual(person, tt.exp.person) {
				t.Errorf("TEST ERROR: \n"+
					"Got: person = %v, isErr = %t \n"+
					"Expected: person = %v, isErr = %t \n",
					person, isError, tt.exp.person, tt.exp.isErr,
				)
			}
		})
	}
}

func TestAddPerson(t *testing.T) {
	type Expected struct {
		id    int
		isErr bool
	}

	tests := []struct {
		name        string
		prepareMock func(m *mocks, p *model.Person)
		arg         *model.Person
		exp         *Expected
	}{
		{
			name: "test with successful adding new person with all fields",
			prepareMock: func(m *mocks, p *model.Person) {
				var err error = nil
				gomock.InOrder(
					m.repo.EXPECT().AddPerson(p).Return(1, err).Times(1),
				)
			},
			arg: &model.Person{0, "Ivan", "Ivanov", "Ivanovich", 40, "male", "RU"},
			exp: &Expected{
				id:    1,
				isErr: false,
			},
		},
		{
			name: "test with successful adding new person with only 'name', 'surname' fields",
			prepareMock: func(m *mocks, p *model.Person) {
				var err error = nil
				gomock.InOrder(
					m.repo.EXPECT().AddPerson(p).Return(1, err).Times(1),
				)
			},
			arg: &model.Person{Name: "Ivan", Surname: "Ivanov"},
			exp: &Expected{
				id:    1,
				isErr: false,
			},
		},
		{
			name: "test with internal error",
			prepareMock: func(m *mocks, p *model.Person) {
				var err error = errors.New("INTERNAL ERROR")
				gomock.InOrder(
					m.repo.EXPECT().AddPerson(p).Return(0, err).Times(1),
				)
			},
			arg: &model.Person{Name: "Ivan", Surname: "Ivanov"},
			exp: &Expected{
				id:    0,
				isErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockPersonRepo := mock_iface.NewMockPersonRepo(ctrl)
			mockCacher := mock_iface.NewMockIcacher(ctrl)

			m := &mocks{mockPersonRepo, mockCacher}
			tt.prepareMock(m, tt.arg)
			logger, _ := logger.NewZapLogger()
			srvc := service.NewPersonService(mockPersonRepo, mockCacher, logger)
			id, err := srvc.AddPerson(tt.arg)
			isError := (err != nil)
			if isError != tt.exp.isErr || id != tt.exp.id {
				t.Errorf("TEST ERROR: \n"+
					"Got: id = %d isErr = %t \n"+
					"Expected: id = %d, isErr = %t \n",
					id, isError, tt.exp.id, tt.exp.isErr,
				)
			}
		})
	}
}

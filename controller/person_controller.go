package controller

import (
	"fio-service/iface"
	"fio-service/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PersonController struct {
	service iface.PersonService
	logger  iface.Ilogger
}

func NewPersonController(service iface.PersonService, logger iface.Ilogger) *PersonController {
	return &PersonController{service, logger}
}

func (pc *PersonController) GetPeopleByFilters(ctx *gin.Context) {
	pc.logger.Info(ctx.Request.Method + ": " + ctx.Request.URL.Path)
	filters, err := pc.getFilters(ctx)
	if err != nil {
		pc.logger.Info("StatusBadRequest: " + err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	people, err := pc.service.GetPeopleByFilters(filters)
	if err != nil {
		pc.logger.Error("StatusInternalServerError: " + err.Error())
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	pc.logger.Info("StatusOK: " + pc.logger.Json(people))
	ctx.JSON(http.StatusOK, people)
}

func (pc *PersonController) GetPersonById(ctx *gin.Context) {
	pc.logger.Info(ctx.Request.Method + ": " + ctx.Request.URL.Path)
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pc.logger.Info("StatusBadRequest: Wrong id")
		ctx.String(http.StatusBadRequest, "Wrong id")
		return
	}
	if id <= 0 {
		pc.logger.Info("StatusBadRequest: Id must be positive")
		ctx.String(http.StatusBadRequest, "Id must be positive")
		return
	}

	person, err := pc.service.GetPersonById(id)
	if err != nil {
		pc.logger.Error("StatusInternalServerError: " + err.Error())
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	pc.logger.Info("StatusOK: " + pc.logger.Json(person))
	ctx.JSON(http.StatusOK, person)
}

func (pc *PersonController) AddPerson(ctx *gin.Context) {
	pc.logger.Info(ctx.Request.Method + ": " + ctx.Request.URL.Path)
	var person model.Person
	err := ctx.BindJSON(&person)
	if err != nil {
		pc.logger.Info("StatusBadRequest: Wrong person's parameters. " + err.Error())
		ctx.String(http.StatusBadRequest, "Wrong person's parameters")
		return
	}
	err = person.Validate()
	if err != nil {
		pc.logger.Info("StatusBadRequest: " + err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	personId, err := pc.service.AddPerson(&person)
	if err != nil {
		pc.logger.Error("StatusInternalServerError: " + err.Error())
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	pc.logger.Info("StatusOK: " + strconv.Itoa(personId))
	ctx.String(http.StatusOK, strconv.Itoa(personId))
}

func (pc *PersonController) DeletePerson(ctx *gin.Context) {
	pc.logger.Info(ctx.Request.Method + ": " + ctx.Request.URL.Path)
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pc.logger.Info("StatusBadRequest: Wrong id")
		ctx.String(http.StatusBadRequest, "Wrong id")
		return
	}
	if id <= 0 {
		pc.logger.Info("StatusBadRequest: Id must be positive")
		ctx.String(http.StatusBadRequest, "Id must be positive")
		return
	}

	person, err := pc.service.DeletePersonById(id)
	if err != nil {
		pc.logger.Error("StatusInternalServerError: " + err.Error())
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	pc.logger.Info("StatusOK: " + pc.logger.Json(person))
	ctx.JSON(http.StatusOK, person)
}

func (pc *PersonController) UpdatePerson(ctx *gin.Context) {
	pc.logger.Info(ctx.Request.Method + ": " + ctx.Request.URL.Path)

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pc.logger.Info("StatusBadRequest: Wrong id")
		ctx.String(http.StatusBadRequest, "Wrong id")
		return
	}

	if id <= 0 {
		pc.logger.Info("StatusBadRequest: Id must be positive")
		ctx.String(http.StatusBadRequest, "Id must be positive")
		return
	}

	var updatePerson model.UpdatePerson
	err = ctx.BindJSON(&updatePerson)
	if err != nil {
		pc.logger.Info("StatusBadRequest: Wrong person's parameters. " + err.Error())
		ctx.String(http.StatusBadRequest, "Wrong person's parameters")
		return
	}
	updatePerson.Id = id
	err = updatePerson.Validate()
	if err != nil {
		pc.logger.Info("StatusBadRequest: " + err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	person, err := pc.service.UpdatePerson(&updatePerson)
	if err != nil {
		pc.logger.Error("StatusInternalServerError: " + err.Error())
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	pc.logger.Info("StatusOK: " + pc.logger.Json(person))
	ctx.JSON(http.StatusOK, person)
}

func (pc *PersonController) getFilters(ctx *gin.Context) (*model.Filters, error) {
	var filters model.Filters

	ageFromStr := ctx.Query("age_from")
	if ageFromStr != "" {
		ageFrom, err := strconv.Atoi(ageFromStr)
		if err != nil {
			return nil, err
		}
		filters.AgeFrom = ageFrom
	}

	ageToStr := ctx.Query("age_to")
	if ageToStr != "" {
		ageTo, err := strconv.Atoi(ageFromStr)
		if err != nil {
			return nil, err
		}
		filters.AgeTo = ageTo
	}

	onPageStr := ctx.Query("on_page")
	if onPageStr != "" {
		onPage, err := strconv.Atoi(onPageStr)
		if err != nil {
			return nil, err
		}
		filters.OnPage = onPage
	}

	pageNumStr := ctx.Query("page_num")
	if pageNumStr != "" {
		pageNum, err := strconv.Atoi(pageNumStr)
		if err != nil {
			return nil, err
		}
		filters.PageNum = pageNum
	}

	filters.Search = ctx.Query("search")
	filters.Gender = ctx.Query("gender")
	filters.Country = ctx.Query("country")

	err := filters.Validate()
	if err != nil {
		return nil, err
	}
	return &filters, nil
}

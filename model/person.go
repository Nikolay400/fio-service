package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Person struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	Country    string `json:"country"`
}

type UpdatePerson struct {
	Id         int     `json:"id"`
	Name       *string `json:"name"`
	Surname    *string `json:"surname"`
	Age        *int    `json:"age"`
	Patronymic *string `json:"patronymic"`
	Gender     *string `json:"gender"`
	Country    *string `json:"country"`
}

type FailedPerson struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Error      string `json:"error"`
}

type Filters struct {
	AgeFrom int    `json:"age_from"`
	AgeTo   int    `json:"age_to"`
	Gender  string `json:"gender"`
	Country string `json:"country"`
	Search  string `json:"search"`
	PageNum int    `json:"page_num"`
	OnPage  int    `json:"on_page"`
}

type nationalizeResponce struct {
	Country []countryId `json:"country"`
}

type countryId struct {
	CountryId   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

func (person *Person) GetAgeGenderCountry() error {
	resp, err := http.Get("https://api.agify.io/?name=" + person.Name)
	if err != nil {
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(&person)
	if err != nil {
		return err
	}

	resp, err = http.Get("https://api.genderize.io/?name=" + person.Name)
	if err != nil {
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(&person)
	if err != nil {
		return err
	}

	resp, err = http.Get("https://api.nationalize.io/?name=" + person.Name)
	if err != nil {
		return err
	}
	var nResponce nationalizeResponce
	err = json.NewDecoder(resp.Body).Decode(&nResponce)
	if err != nil {
		return err
	}
	person.Country = nResponce.GetCountry()
	return nil
}

func (person *Person) Validate() error {
	if person.Name == "" {
		return errors.New("Error: Empty name")
	}
	if person.Surname == "" {
		return errors.New("Error: Empty surname")
	}
	return person.validateFormat()
}

func (person *Person) validateFormat() error {
	if person.Age < 0 {
		return errors.New("Error: Negative Age")
	}
	re := regexp.MustCompile("[0-9]")
	if person.Name != "" && re.FindString(person.Name) != "" {
		return errors.New("Error: Wrong Name format")
	}
	if person.Name != "" && re.FindString(person.Surname) != "" {
		return errors.New("Error: Wrong Surname format")
	}
	if person.Patronymic != "" && re.FindString(person.Patronymic) != "" {
		return errors.New("Error: Wrong Patronymic format")
	}
	if person.Patronymic != "" && re.FindString(person.Gender) != "" {
		return errors.New("Error: Wrong Gender format")
	}
	if person.Country != "" && (len([]byte(person.Country)) != 2 || re.FindString(person.Country) != "") {
		return errors.New("Error: Wrong Country format")
	}
	return nil
}

func (updatePerson *UpdatePerson) Validate() error {
	person := updatePerson.ConvertToPerson()
	return person.validateFormat()
}

func (updatePerson *UpdatePerson) ConvertToPerson() *Person {
	person := &Person{}
	if updatePerson.Name != nil {
		person.Name = *updatePerson.Name
	}
	if updatePerson.Surname != nil {
		person.Surname = *updatePerson.Surname
	}
	if updatePerson.Patronymic != nil {
		person.Patronymic = *updatePerson.Patronymic
	}
	if updatePerson.Age != nil {
		person.Age = *updatePerson.Age
	}
	if updatePerson.Gender != nil {
		person.Gender = *updatePerson.Gender
	}
	if updatePerson.Country != nil {
		person.Country = *updatePerson.Country
	}
	return person
}

func (filters *Filters) Validate() error {
	if filters.AgeFrom < 0 {
		return errors.New("Error: Negative AgeFrom")
	}
	if filters.AgeTo < 0 {
		return errors.New("Error: Negative AgeTo")
	}
	if filters.PageNum < 0 {
		return errors.New("Error: Negative PageNum")
	}
	if filters.OnPage < 0 {
		return errors.New("Error: Negative OnPage")
	}
	return nil
}

func (nr *nationalizeResponce) GetCountry() string {
	max := float32(0.0)
	var result string
	for i := range nr.Country {
		if nr.Country[i].Probability > max {
			max = nr.Country[i].Probability
			result = nr.Country[i].CountryId
		}
	}
	return result
}

func (person *Person) MarshalBinary() ([]byte, error) {
	return json.Marshal(person)
}

func (f *Filters) GetKeyForRedis() string {
	str := fmt.Sprintf("AgeFrom=%d,AgeTo=%d,Gender=%s,Country=%s,Search=%s,PageNum=%d,OnPage=%d", f.AgeFrom, f.AgeTo, f.Gender, f.Country, f.Search, f.PageNum, f.OnPage)
	return createMd5Hash(str)
}

func createMd5Hash(text string) string {
	hasher := md5.New()
	_, err := io.WriteString(hasher, text)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

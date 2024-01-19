package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)


type OpeningYearRequest struct {
	AccessKey int64  `json:"access_key"`
	OpeningYear int `json:"year"`
}

type Request struct {
	StationId int64 `json:"station_id"`
}


func (h *Handler) issueOpeningYear(c *gin.Context) {
	var input Request
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("Год ввода в строй АЭС с ID №", input.StationId)

	c.Status(http.StatusOK)

	go func() {
		time.Sleep(4 * time.Second)
		sendOpeningYearRequest(input)
	}()
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	if min > max {
		return min
	} else {
		return rand.Intn(max-min) + min
	}
}


func sendOpeningYearRequest(request Request) {

	var year = random(1950, 2023)

	answer := OpeningYearRequest{
		AccessKey: 123,
		OpeningYear: year,
	}

	client := &http.Client{}

	jsonAnswer, _ := json.Marshal(answer)
	bodyReader := bytes.NewReader(jsonAnswer)

	requestURL := fmt.Sprintf("http://127.0.0.1:8000/api/stations/%d/update_year/", request.StationId)

	req, _ := http.NewRequest(http.MethodPut, requestURL, bodyReader)

	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending PUT request:", err)
		return
	}

	defer response.Body.Close()

	fmt.Println("Год ввода в строй: ", year)
	fmt.Println("PUT Request Status:", response.Status)
}

package handlers

import (
	"Users/vaibhav.sabharwal/verve/model"
	"Users/vaibhav.sabharwal/verve/services"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
)

func AcceptHandler(c echo.Context) error {
	idParam := c.QueryParam("id")
	endpointParam := c.QueryParam("endpoint")
	res := model.GetAcceptResponse{}

	id, err := strconv.Atoi(idParam) //validating id
	if err != nil {
		res.ResponseMessage = "failed"
		return c.JSON(http.StatusBadRequest, res)
	}

	var endpoint *url.URL
	if len(endpointParam) != 0 { //validating endpoint
		endpoint, err = url.Parse(endpointParam)
		if err != nil {
			log.Printf("Invalid endpoint: %v", err)
		}
	}

	err = services.HandleRequest(id, endpoint)
	if err != nil {
		res.ResponseMessage = "failed"
		return c.JSON(http.StatusInternalServerError, res)
	}

	res.ResponseMessage = "ok"
	return c.JSON(http.StatusOK, res)
}

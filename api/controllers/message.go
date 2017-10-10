package controllers

import "github.com/labstack/echo"
import (
	"api/models"
	"net/http"
)

// HandleMessage controller
func HandleMessage(c echo.Context) error {
	m := models.InitMessage()
	err := c.Bind(m)

	if err != nil {
		return err
	}

	err = c.Validate(m)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m)
}

package controllers

import (
	"api/models"
	"github.com/labstack/echo"
	"net/http"
	"utils"
)

// HandleMessage controller
func HandleMessage(c echo.Context) error {
	var err error

	m := models.InitMessage()

	if err = c.Bind(m); err != nil {
		return err
	}

	if err = c.Validate(m); err != nil {
		t := utils.HumaniseValidationErrors(err)
		return c.JSON(http.StatusUnprocessableEntity, t)
	}

	return c.JSON(http.StatusOK, m)
}

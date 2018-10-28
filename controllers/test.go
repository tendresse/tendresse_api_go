package controllers

import (
	"github.com/labstack/echo"
	"net/http"
)

func IsUp(c echo.Context) error {
	return c.String(http.StatusOK, "")
}

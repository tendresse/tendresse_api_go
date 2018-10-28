package middlewares

import (
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
)

func LinkDB(db *pg.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("DB", db)
			return next(c)
		}
	}
}

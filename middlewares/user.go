package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/tendresse/tendresse_api_go/models"
	"net/http"
)

func VerifyUser() echo.MiddlewareFunc {
	// TODO optimiser requetes, .Column semble faire N+1
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			db := c.Get("DB").(*pg.DB)
			token := c.Get("token").(*jwt.Token)
			claims := token.Claims.(*MyCustomClaims)
			user := new(models.User)
			user.ID = claims.UserID
			err := db.Model(user).
				Where("id = ?", user.ID).
				Column("id").
				First()
			if err != nil {
				err = errors.Wrap(err, "check if user exists by ID")
				c.Logger().Error(err)
				return c.NoContent(http.StatusUnauthorized)
			}
			tokens := new([]*models.Token)
			err = db.Model(tokens).
				Where("user_id = ?", user.ID).
				Select()
			if err != nil {
				err = errors.Wrap(err, "getting user tokens")
				c.Logger().Error(err)
				return c.NoContent(http.StatusUnauthorized)
			}
			for _, user_token := range *tokens {
				if user_token.Hash == token.Raw {
					c.Set("current_user_id", claims.UserID)
					c.Set("current_user", user)
					return next(c)
				}
			}
			c.Logger().Error("this token is not valid")
			return c.NoContent(http.StatusUnauthorized)
		}
	}
}

func LoadUser() echo.MiddlewareFunc {
	// TODO optimiser requetes, .Column semble faire N+1
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			db := c.Get("DB").(*pg.DB)
			token := c.Get("token").(*jwt.Token)
			claims := token.Claims.(*MyCustomClaims)
			user := new(models.User)
			err := db.Model(user).
				Where("id = ?", claims.UserID).
				Column("user.*", "Tokens").
				Column("user.*", "Roles").
				First()
			if err != nil {
				c.Logger().Error(err)
				return c.NoContent(http.StatusUnauthorized)
			}
			for _, user_token := range user.Tokens {
				if user_token.Hash == token.Raw {
					c.Set("current_user", user)
					return next(c)
				}
			}
			c.Logger().Error(err)
			return c.NoContent(http.StatusUnauthorized)
		}
	}
}

func RolesRequired(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			current_user := c.Get("current_user_id").(*models.User)
			db := c.Get("DB").(*pg.DB)
			user_roles := new([]*models.Role)
			err := db.Model(roles).
				Where("user_id = ?", current_user.ID).
				Select()
			if err != nil {
				err = errors.Wrap(err, "get user roles")
				c.Logger().Error(err)
				return c.NoContent(http.StatusUnauthorized)
			}
			for _, required_role := range roles {
				for _, user_role := range *user_roles {
					if required_role == user_role.Name {
						return next(c)
					}
				}
			}
			c.Logger().Error("a user without the required role tried to access this route")
			return c.NoContent(http.StatusUnauthorized)
		}
	}
}

package dao

import (
	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/tendresse/tendresse_api_go/database"
	"github.com/tendresse/tendresse_api_go/middlewares"
	"github.com/tendresse/tendresse_api_go/models"
	"os"
	"time"
)

func GenerateToken(user *models.User) (string, error) {
	db := database.GetDB()
	jwt, err := generateToken(user.ID)
	if err != nil {
		return "", err
	}
	token := &models.Token{
		Hash:   jwt,
		UserID: user.ID,
	}
	err = db.Insert(token)
	if err != nil {
		err = errors.Wrap(err, "inserting token to DB")
		return "", err
	}
	return jwt, nil
}

func GetUserTokens(user *models.User, tokens *[]*models.Token) error {
	db := database.GetDB()
	err := db.Model(tokens).
		Where("user_id = ?", user.ID).
		Select()
	return errors.Wrap(err, "getting user tokens")
}

func generateToken(user_id int) (string, error) {
	claims := &middlewares.MyCustomClaims{
		user_id,
		jwtGo.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 7*2).Unix(),
		},
	}
	token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		err = errors.Wrap(err, "creating JWT")
	}
	return tokenString, err
}

package dao

import (
	"github.com/pkg/errors"
	"github.com/tendresse/tendresse_api_go/database"
	"github.com/tendresse/tendresse_api_go/models"
)

func GetOrCreateRole(role *models.Role) error {
	db := database.GetDB()
	_, err := db.Model(role).
		Column("id").
		Where("name = ?", role.Name).
		OnConflict("DO NOTHING"). // OnConflict is optional
		Returning("id").
		SelectOrInsert()
	return err
}

func DeleteRole(role *models.Role) error {
	db := database.GetDB()
	err := db.Delete(&role)
	return errors.Wrap(err, "deleting role")
}

func GetRoles(roles *[]*models.Role) error {
	db := database.GetDB()
	err := db.Model(roles).Select()
	return errors.Wrap(err, "get all roles")
}

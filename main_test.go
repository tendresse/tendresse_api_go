package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/tendresse/tendresse_api_go/controllers"
	"github.com/tendresse/tendresse_api_go/database"
	"github.com/tendresse/tendresse_api_go/middlewares"
	"github.com/tendresse/tendresse_api_go/models"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func initEnv() {
	// docker run --name test-postgres-tendresse -e POSTGRES_PASSWORD=tendresse_password -e POSTGRES_USER=tendresse_user -e POSTGRES_DB=tendresse_test -p 127.0.0.1:5433:5432 -d postgres:10.5
	// loading configuration
	godotenv.Load("test.env")
	required_vars := []string{"PORT", "GO_ENV", "TUMBLR_API_KEY", "DATABASE_URL", "SECRET_KEY"}
	for _, required_var := range required_vars {
		if os.Getenv(required_var) == "" {
			middlewares.Logger().Fatalf("%q env var is not set.", required_var)
		}
	}
}

func TestMain(t *testing.T) {
	initEnv()
	// loading database
	database.Init()
	db := database.GetDB()
	defer database.CloseDB()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Use(middlewares.LinkDB(db))

	// PERSONNES
	t.Run("ListPersonne", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/personnes")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.ListPersonne(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("GetPersonne", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/personnes/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.GetPersonne(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
	p := &models.Personne{
		Nom:    "Dupont",
		Prenom: "Jean",
	}
	t.Run("CreatePersonne", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(p.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/personnes")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.CreatePersonne(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("GetPersonne", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/personnes/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.GetPersonne(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	p.ID = 1
	p.Prenom = "Jeanne"
	t.Run("UpdatePersonne", func(t *testing.T) {
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(p.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/personnes/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.UpdatePersonne(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	// ETABLISSEMENTS
	t.Run("ListEtablissement", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/etablissements")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.ListEtablissement(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("GetEtablissement", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/etablissements/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.GetEtablissement(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
	eta := &models.Etablissement{
		Nom:         "Thatha",
		Departement: 44,
		Contrat:     "AVEC DETENTION",
		Poste:       "NTS",
		CodeSe:      "Z4444",
		Ville:       "NANTES",
	}
	t.Run("CreateEtablissement", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(eta.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/etablissements")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.CreateEtablissement(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	eta.ID = 1
	eta.DirecteurId = 1
	eta.OseId = 1
	t.Run("UpdateEtablissement", func(t *testing.T) {
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(eta.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/etablissements/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.UpdateEtablissement(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("GetEtablissement", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/etablissements/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.GetEtablissement(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	// SIS
	t.Run("ListSi", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/sis")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.ListSi(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("GetSi", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/sis/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.GetSi(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
	si := &models.Si{
		Nom:             "SiDeTest",
		Classification:  "CD",
		EtablissementId: 1,
	}
	t.Run("CreateSi", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(si.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/sis")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.CreateSi(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("GetSi", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/sis/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.GetSi(c)) {
			t.Log(rec.Body)
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	si.ID = 1
	si.DateFinDeVie = "2222-22-22"
	t.Run("UpdateSi", func(t *testing.T) {
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(p.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/sis/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.UpdateSi(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	// ATAPS
	t.Run("ListAtap", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/ataps")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.ListAtap(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("GetAtap", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/ataps/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.GetAtap(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
	atap := &models.Atap{
		DateSignature:   "1111-11-11",
		Reference:       "REF/ATAP/TEST",
		EtablissementId: 1,
	}
	t.Run("CreateAtap", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(atap.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/ataps")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.CreateAtap(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("GetAtap", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/ataps/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.GetAtap(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	atap.ID = 1
	atap.Commentaire = "test commentaire"
	t.Run("UpdateAtap", func(t *testing.T) {
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(atap.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/ataps/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.UpdateAtap(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("LinkSiWithAtap", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/ataps/:id_atap/sis/:id_si")
		c.SetParamNames("id_atap", "id_si")
		c.SetParamValues("1", "1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.LinkSiWithAtap(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	// DSSIS
	t.Run("ListDssi", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/dssis")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.ListDssi(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("GetDssi", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/dssis/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.GetDssi(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
	dssi := &models.Dssi{
		DateCreation:    "2222-22-22",
		Reference:       "DSSI/REF/23/TEST",
		EtablissementId: 1,
	}
	t.Run("CreateDssi", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(dssi.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/dssis")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.CreateDssi(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("GetDssi", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/dssis/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.GetDssi(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	dssi.ID = 1
	dssi.DateReception = "2222-22-23"
	t.Run("UpdateDssi", func(t *testing.T) {
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(dssi.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/dssis/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.UpdateDssi(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("LinkSiWithDssi", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/dssis/:id_dssi/sis/:id_si")
		c.SetParamNames("id_dssi", "id_si")
		c.SetParamValues("1", "1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.LinkSiWithDssi(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	// ATAIS
	t.Run("ListAtai", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/atais")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.ListAtai(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("GetAtai", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/atais/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.GetAtai(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
	atai := &models.Atai{
		DateCreation: "2222-22-11",
		Reference:    "REF?ATAI/TEST",
		DssiId:       1,
	}
	t.Run("CreateAtai", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(atai.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/atais")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.CreateAtai(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("GetAtai", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/atais/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.GetAtai(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	atai.ID = 1
	atai.Document = "/doc/atai.pdf"
	t.Run("UpdateAtai", func(t *testing.T) {
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(atai.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/atais/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.UpdateAtai(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	homologation := &models.Homologation{
		SiId:   1,
		AtaiId: 1,
	}
	t.Run("AddHomologation", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(homologation.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/atais/:id_atai/sis/:id_si")
		c.SetParamNames("id_atai", "id_si")
		c.SetParamValues("1", "1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.AddHomologation(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	homologation.Reference = "REF?HOMOLOG/4/TEST"
	homologation.DateExpiration = "1111-11-11"
	t.Run("UpdateHomologation", func(t *testing.T) {
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(homologation.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/atais/:id_atai/sis/:id_si")
		c.SetParamNames("id_atai", "id_si")
		c.SetParamValues("1", "1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.UpdateHomologation(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	// VISITES
	t.Run("ListVisite", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/visites")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.ListVisite(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("GetVisite", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/visites/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.GetVisite(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
	visite := &models.Visite{
		DateDebut:       "1111-11-11",
		Heures:          2.0,
		Type:            "CONTROLE",
		EtablissementId: 1,
	}
	t.Run("CreateVisite", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(visite.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/visites")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.CreateVisite(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("GetVisite", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/visites/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.GetVisite(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	visite.ID = 1
	visite.Document = "/doc/visite_test.doc"
	t.Run("UpdateVisite", func(t *testing.T) {
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(visite.ToString()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/visites/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.UpdateVisite(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
	t.Run("LinkSiWithVisite", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/visites/:id_visite/sis/:id_si")
		c.SetParamNames("id_visite", "id_si")
		c.SetParamValues("1", "1")
		c.Set("DB", db)
		// assertions
		if assert.NoError(t, controllers.LinkSiWithVisite(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	// CLEANING
	// replace t.Run by t.Skip instead of commenting the all thing
	t.Run("Cleaning", func(t *testing.T) {
		t.Run("UnlinkSiWithAtap", func(t *testing.T) {
			req := httptest.NewRequest(echo.DELETE, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/ataps/:id_atap/sis/:id_si")
			c.SetParamNames("id_atap", "id_si")
			c.SetParamValues("1", "1")
			c.Set("DB", db)
			// assertions
			if assert.NoError(t, controllers.UnlinkSiWithAtap(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
		t.Run("UnlinkSiWithDssi", func(t *testing.T) {
			req := httptest.NewRequest(echo.DELETE, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/dssis/:id_dssi/sis/:id_si")
			c.SetParamNames("id_dssi", "id_si")
			c.SetParamValues("1", "1")
			c.Set("DB", db)
			// assertions
			if assert.NoError(t, controllers.UnlinkSiWithDssi(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
		t.Run("DeleteHomologation", func(t *testing.T) {
			req := httptest.NewRequest(echo.DELETE, "/", strings.NewReader(homologation.ToString()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/atais/:id_atai/sis/:id_si")
			c.SetParamNames("id_atai", "id_si")
			c.SetParamValues("1", "1")
			c.Set("DB", db)
			// assertions
			if assert.NoError(t, controllers.DeleteHomologation(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
		t.Run("UnlinkSiWithVisite", func(t *testing.T) {
			req := httptest.NewRequest(echo.DELETE, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/visites/:id_visite/sis/:id_si")
			c.SetParamNames("id_visite", "id_si")
			c.SetParamValues("1", "1")
			c.Set("DB", db)
			// assertions
			if assert.NoError(t, controllers.UnlinkSiWithVisite(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
		t.Run("DeleteVisite", func(t *testing.T) {
			req := httptest.NewRequest(echo.DELETE, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/visites/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			c.Set("DB", db)
			// assertions
			if assert.NoError(t, controllers.DeleteVisite(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
		t.Run("DeleteAtai", func(t *testing.T) {
			req := httptest.NewRequest(echo.DELETE, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/atais/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			c.Set("DB", db)
			// assertions
			if assert.NoError(t, controllers.DeleteAtai(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
		t.Run("DeleteDssi", func(t *testing.T) {
			req := httptest.NewRequest(echo.DELETE, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/dssis/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			c.Set("DB", db)
			// assertions
			if assert.NoError(t, controllers.DeleteDssi(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
		t.Run("DeleteAtap", func(t *testing.T) {
			req := httptest.NewRequest(echo.DELETE, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/ataps/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			c.Set("DB", db)
			// assertions
			if assert.NoError(t, controllers.DeleteAtap(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
		t.Run("DeleteSi", func(t *testing.T) {
			req := httptest.NewRequest(echo.DELETE, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/sis/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			c.Set("DB", db)
			// assertions
			if assert.NoError(t, controllers.DeleteSi(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
		t.Run("DeleteEtablissement", func(t *testing.T) {
			req := httptest.NewRequest(echo.DELETE, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/etablissements/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			c.Set("DB", db)
			// assertions
			if assert.NoError(t, controllers.DeleteEtablissement(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
		t.Run("DeletePersonne", func(t *testing.T) {
			req := httptest.NewRequest(echo.DELETE, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/personnes/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			c.Set("DB", db)
			// assertions
			if assert.NoError(t, controllers.DeletePersonne(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
	})
}

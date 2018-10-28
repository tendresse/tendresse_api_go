package app

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tendresse/tendresse_api_go/controllers"
	"github.com/tendresse/tendresse_api_go/database"
	_ "github.com/tendresse/tendresse_api_go/docs"
	"github.com/tendresse/tendresse_api_go/middlewares"
	"os"
)

func InitEnv() {
	// loading configuration
	godotenv.Load("vars.env")
	required_vars := []string{"PORT", "GO_ENV", "TUMBLR_API_KEY", "DATABASE_URL", "SECRET_KEY"}
	for _, required_var := range required_vars {
		if os.Getenv(required_var) == "" {
			middlewares.Logger().Fatalf("%q env var is not set.", required_var)
		}
	}
}

func Launch() {
	InitEnv()

	database.Init()
	db := database.GetDB()
	defer database.CloseDB()

	e := echo.New()
	e.Logger = middlewares.Logger()
	e.Use(middlewares.LoggerHook())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Use(middlewares.LinkDB(db))
	config := middleware.JWTConfig{
		Claims:     &middlewares.MyCustomClaims{},
		SigningKey: []byte(os.Getenv("SECRET_KEY")),
		ContextKey: "token",
	}

	e.POST("/signup", controllers.UserSignup)
	e.POST("/login", controllers.UserLogin)

	v1 := e.Group("/v1")
	v1.Use(middleware.JWTWithConfig(config))
	v1.Use(middlewares.VerifyUser())
	{
		v1.GET("/achievements", controllers.GetAchievementsLimited)
		v1.GET("/random", controllers.GetRandomGif)
		v1.GET("/me/profile", controllers.GetCurrentUserProfile)
		v1.GET("/me/tendresses", controllers.GetTendresses)
		v1.GET("/me/friends", controllers.GetCurrentUserFriends)
		v1.GET("/users/:username/profile", controllers.GetUserProfile)
		v1.POST("/me/friends/:username", controllers.AddFriend)
		v1.DELETE("/me/friends/:username", controllers.RemoveFriend)
		v1.POST("/me/friends/:username/tendresses", controllers.SendTendresse)
		v1.PUT("/me/tendresses/:tendresse_id", controllers.ChangeTendresseAsViewed)
	}
	restricted := v1.Group("/restricted")
	restricted.Use(middlewares.RolesRequired([]string{"admin", "contributor"}))
	{
		restricted.GET("/blogs", controllers.GetBlogs)
		restricted.POST("/blogs", controllers.AddBlog)
		restricted.DELETE("/blogs/:blog_id", controllers.DeleteBlog)
		restricted.DELETE("/blogs/:blog_id/full", controllers.DeleteBlogAndGifs)

		restricted.GET("/tags", controllers.GetTags)
		restricted.GET("/tags/:tag_id", controllers.GetTag)
		restricted.DELETE("/tags/:tag_id", controllers.DeleteTag)

		restricted.GET("/gifs", controllers.GetGifs)
		restricted.POST("/gifs", controllers.AddGif)
		restricted.PUT("/gifs/:gif_id", controllers.UpdateGifTags)
		restricted.DELETE("/gifs/:gif_id", controllers.DeleteGif)
	}

	admin := v1.Group("/admin")
	admin.Use(middlewares.RolesRequired([]string{"admin"}))
	{
		admin.GET("/check/gifs", controllers.CheckAllGifs)
		admin.GET("/roles", controllers.GetRoles)
		admin.POST("/users/:username/roles/:role_name", controllers.AddRoleToUser)
		admin.DELETE("/users/:username/roles/:role_name", controllers.RemoveRoleFromUser)
		//admin.GET("/users/:username", controllers.GetUser)

		admin.GET("/achievements", controllers.GetAchievements)
		admin.POST("/achievements", controllers.AddAchievement)
		// admin.PUT("/achievements/:achievement_id", controllers.UpdateAchievement)
		admin.DELETE("/achievements/:achievement_id", controllers.DeleteAchievement)
	}

	//e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(os.Getenv("PORT")))
}

package main

import (
	"fmt"
	"ldap/db"
	"ldap/internal"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type application struct {
	db *db.UserDB
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.Default().SetTrustedProxies(nil)

	e := gin.New()
	e.Use(sessions.Sessions(
		"session",
		cookie.NewStore([]byte(internal.SECRET)),
	))

	db, err := db.NewDatabase()
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	app := application{db: db}

	e.GET("/healthz", app.health)
	e.GET("/users", app.users)
	e.GET("/protected", internal.Auth(db), app.protected)

	e.Run(fmt.Sprintf("0.0.0.0:%s", internal.PORT))
}

func (a application) health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func (a application) users(ctx *gin.Context) {
	data, err := db.GetAll(a.db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to get all users",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    data,
	})
}

func (a application) protected(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

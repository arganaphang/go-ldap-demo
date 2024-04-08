package internal

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"ldap/db"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Auth(database *db.UserDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		var data *UserLDAPData
		var err error
		var ok bool
		if user == nil {
			if c.GetHeader("Authorization") == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized [auth]",
				})
				return
			}
			// Get Basic Auth
			auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
			if len(auth) != 2 || auth[0] != "Basic" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized [basic]",
				})
				return
			}
			payload, _ := base64.StdEncoding.DecodeString(auth[1])
			pair := strings.SplitN(string(payload), ":", 2)
			if len(pair) != 2 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized [basic]",
				})
				return
			}
			ok, data, err = AuthUsingLDAP(pair[0], pair[1])
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized [ldap incorrect]",
				})
				return
			}
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized [ldap error]",
				})
				return
			}

		} else {
			var userSession db.User
			err := json.Unmarshal([]byte(user.(string)), &userSession)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized [unmarshal session error]",
				})
				return
			}

			ok, data, err = AuthUsingLDAP(userSession.Username, "password") // IF DN with no password is allow use this function `ok, data, err = AuthUsingLDAPPasswordless(userSession.Username)`
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized [ldap incorrect]",
				})
				return
			}
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized [ldap error]",
				})
				return
			}
		}

		savedUser, err := db.Get(database, data.Name)
		if err == sql.ErrNoRows {
			if err := db.Add(database, data.Name, data.FullName, data.Email); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized [add user error]",
				})
				return
			}
			savedUser, err = db.Get(database, data.Name)
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized [get user error]",
			})
			return
		}
		// Set Session
		userData, _ := json.Marshal(savedUser)
		session.Set("user", string(userData))
		session.Save()
		// Logic Get LDAP
		c.Next()
	}
}

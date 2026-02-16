package handlers

import (
	"net/http"

	"inventory-system/database"
	"inventory-system/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{"title": "Login"})
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user models.User
	err := database.DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.PasswordHash)

	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
		return
	}

	session := sessions.Default(c)
	session.Set("user", user.ID)
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}

func Profile(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")

	var user models.User
	err := database.DB.QueryRow("SELECT id, username, email FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		c.String(http.StatusInternalServerError, "User not found")
		return
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"title": "Profile Settings",
		"user":  user,
	})
}

func UpdateProfile(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")

	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	if username == "" || email == "" {
		c.HTML(http.StatusBadRequest, "profile.html", gin.H{"error": "Username and Email are required"})
		return
	}

	if password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		_, err := database.DB.Exec("UPDATE users SET username=?, email=?, password_hash=? WHERE id=?",
			username, email, string(hash), userID)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "profile.html", gin.H{"error": "Failed to update profile"})
			return
		}
	} else {
		_, err := database.DB.Exec("UPDATE users SET username=?, email=? WHERE id=?",
			username, email, userID)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "profile.html", gin.H{"error": "Failed to update profile"})
			return
		}
	}

	c.Redirect(http.StatusFound, "/profile?success=1")
}

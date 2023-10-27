package controllers

import (
	"fmt"
	"time"

	"github.com/AYGA2K/photo-editor-web-app/webapp/database"
	"github.com/AYGA2K/photo-editor-web-app/webapp/models"
	"github.com/badoux/checkmail"
	guuid "github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	User    models.User
	Session models.Session
)

func GetUser(sessionid guuid.UUID) (User, error) {
	db := database.Database.Db
	query := Session{Sessionid: sessionid}
	found := Session{}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, err
	}
	user := User{}
	usrQuery := User{ID: found.UserRefer}
	err = db.First(&user, &usrQuery).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, err
	}
	return user, nil
}

func Login(ctx iris.Context) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	db := database.Database.Db
	json := new(LoginRequest)
	if err := ctx.ReadForm(json); err != nil {
		ctx.JSON(iris.Map{"code": 400, "message": "Invalid JSON"})

		return
	}

	found := User{}
	query := User{Username: json.Username}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(iris.Map{
			"code":    404,
			"message": "User not found",
		})
		return
	}
	if !comparePasswords(found.Password, []byte(json.Password)) {
		ctx.JSON(iris.Map{
			"code":    401,
			"message": "Invalid Password",
		})
		return
	}
	session := Session{UserRefer: found.ID, Expires: SessionExpires(), Sessionid: guuid.New()}
	db.Create(&session)
	ctx.SetCookie(&iris.Cookie{
		Name:     "sessionid",
		Path:     "/",
		Expires:  SessionExpires(),
		Value:    session.Sessionid.String(),
		HttpOnly: true,
	})
	ctx.Redirect("/")
}

func Logout(ctx iris.Context) {
	db := database.Database.Db
	json := new(Session)
	if err := ctx.ReadJSON(json); err != nil {
		ctx.JSON(iris.Map{"code": 400, "message": "Invalid JSON"})
		return
	}
	session := Session{}
	query := Session{Sessionid: json.Sessionid}
	err := db.First(&session, &query).Error
	if err == gorm.ErrRecordNotFound {
		ctx.JSON(iris.Map{"code": 404, "message": "Session not found"})
		return
	}
	db.Delete(&session)
	ctx.RemoveCookie("sessionid")
	ctx.JSON(iris.Map{"code": 200, "message": "Logged out successfully"})
}

func CreateUser(ctx iris.Context) {
	type CreateUserRequest struct {
		Password string `json:"password"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	db := database.Database.Db
	json := new(CreateUserRequest)
	if err := ctx.ReadForm(json); err != nil {
		fmt.Println(json)
		ctx.JSON(iris.Map{"code": 400, "message": err.Error()})
		return
	}
	password := hashAndSalt([]byte(json.Password))
	err := checkmail.ValidateFormat(json.Email)
	if err != nil {
		ctx.JSON(iris.Map{"code": 400, "message": "Invalid Email Address"})
		return
	}
	new := User{
		Username: json.Username,
		Password: password,
		Email:    json.Email,
		ID:       guuid.New(),
	}
	found := User{}
	query := User{Username: json.Username}
	err = db.First(&found, &query).Error
	if err != gorm.ErrRecordNotFound {
		ctx.JSON(iris.Map{"code": 400, "message": "Username already exists"})
		return
	}
	db.Create(&new)
	session := Session{UserRefer: new.ID, Sessionid: guuid.New()}
	err = db.Create(&session).Error
	if err != nil {
		ctx.JSON(iris.Map{"code": 500, "message": "Database Error"})
	}
	ctx.SetCookie(&iris.Cookie{
		Name:     "sessionid",
		Path:     "/",
		Expires:  time.Now().Add(5 * 24 * time.Hour),
		Value:    session.Sessionid.String(),
		HttpOnly: true,
	})
	ctx.JSON(iris.Map{"code": 200, "message": "User created successfully"})
}

func GetUserInfo(ctx iris.Context) {
	user, ok := ctx.Values().Get("user").(User)
	if !ok {
		ctx.JSON(iris.Map{"code": 400, "message": "User not found"})
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(user)
}

func DeleteUser(ctx iris.Context) {
	type DeleteUserRequest struct {
		password string
	}
	db := database.Database.Db
	json := new(DeleteUserRequest)
	user, ok := ctx.Values().Get("user").(User)
	if !ok {
		ctx.JSON(iris.Map{"code": 400, "message": "User not found"})
	}

	if err := ctx.ReadJSON(json); err != nil {
		ctx.JSON(iris.Map{"code": 400, "message": "Invalid JSON"})
		return
	}
	if !comparePasswords(user.Password, []byte(json.password)) {
		ctx.JSON(iris.Map{"code": 401, "message": "Invalid Password"})
		return
	}
	if err := db.Model(&user).Association("Sessions").Delete(); err != nil {
		ctx.JSON(iris.Map{"code": 500, "message": "Database Error"})
		return
	}

	db.Delete(&user)
	ctx.RemoveCookie("sessionid")
	ctx.JSON(iris.Map{"code": 200, "message": "User Deleted Successfully"})
}

func ChangePassword(ctx iris.Context) {
	type ChangePasswordRequest struct {
		Password    string `json:"password"`
		NewPassword string `json:"newPassword"`
	}
	db := database.Database.Db
	user, ok := ctx.Values().Get("user").(User)
	if !ok {
		ctx.JSON(iris.Map{"code": 400, "message": "User not found"})
	}
	json := new(ChangePasswordRequest)
	if err := ctx.ReadJSON(json); err != nil {
		ctx.JSON(iris.Map{"code": 400, "message": "Invalid JSON"})
		return
	}
	if !comparePasswords(user.Password, []byte(json.Password)) {
		ctx.JSON(iris.Map{"code": 401, "message": "Invalid Password"})
		return
	}
	user.Password = hashAndSalt([]byte(json.NewPassword))
	db.Save(&user)
	ctx.JSON(iris.Map{"code": 200, "message": "Password Changed Successfully"})
}

func hashAndSalt(pwd []byte) string {
	hash, _ := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	return err == nil
}

// Universal date the Session Will Expire
func SessionExpires() time.Time {
	return time.Now().Add(5 * 24 * time.Hour)
}

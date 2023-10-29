package controllers

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/AYGA2K/photo-editor-web-app/webapp/database"
	"github.com/AYGA2K/photo-editor-web-app/webapp/models"
	"github.com/AYGA2K/photo-editor-web-app/webapp/views"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

// upload image
const maxSize = 8 * iris.MB

func UploadImage(ctx iris.Context) {
	fmt.Println("upload image")
	ctx.SetMaxRequestBodySize(maxSize)
	_, fileHeader, err := ctx.FormFile("file")
	category := ctx.FormValue("category")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	// rename file to unique name and save the name in the database
	sessionId := ctx.GetCookie("sessionid")
	SessionId, err := uuid.Parse(sessionId)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	user, err := GetUser(SessionId)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
	}
	image := models.Image{
		Imageid:   uuid.New(),
		Name:      uuid.New().String() + filepath.Ext(fileHeader.Filename),
		UserRefer: user.ID,
		Category:  category,
	}
	fileHeader.Filename = image.Name
	db := database.Database.Db
	db.Create(&image)
	// Upload the file to specific destination.
	dest := filepath.Join("./uploads", fileHeader.Filename)
	if _, err := ctx.SaveFormFile(fileHeader, dest); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
	}
	ctx.StatusCode(iris.StatusOK)
}

func getImagesByCategory(ctx iris.Context, category string) {
	sessionId := ctx.GetCookie("sessionid")
	parsedSessionId, err := uuid.Parse(sessionId)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	user, err := GetUser(parsedSessionId)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	var images []models.Image
	err = database.Database.Db.Where("user_refer = ? AND category = ?", user.ID, category).Find(&images).Error
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
	}
	var imagesNames []string
	for _, image := range images {
		imagesNames = append(imagesNames, image.Name)
	}
	imagesHTML := views.Images(imagesNames)
	imagesHTML.Render(context.Background(), ctx.ResponseWriter())
}

// get images based on category
func GetImages(ctx iris.Context) {
	category := ctx.Params().Get("category")
	getImagesByCategory(ctx, category)
}

func GetImagesFromSelect(ctx iris.Context) {
	category := ctx.FormValue("category")
	getImagesByCategory(ctx, category)
}

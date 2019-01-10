package handlers

import (
	"encoding/json"
	"github.com/GymWorkoutApp/gwap-files/cache"
	"github.com/GymWorkoutApp/gwap-files/database"
	"github.com/GymWorkoutApp/gwap-files/errors"
	"github.com/GymWorkoutApp/gwap-files/models"
	"github.com/labstack/echo"
	"io"
	"net/http"
	"os"
)

func HandleFilesCreateRequest(c echo.Context) error {

	sourceFile, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := sourceFile.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	source := os.Getenv("MEDIA_URL") + sourceFile.Filename
	dst, err := os.Create(source)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	file := models.File{Source: source, Filename: sourceFile.Filename}

	db := database.GetDB(c.Request().Context())

	err = db.Create(&file).Error
	if err != nil {
		fileJson, err := json.Marshal(file)
		if err != nil {
			return err
		}
		cache.GetRedisClient().HSet(file.GetID(), "file", fileJson)
		return err
	}

	return c.JSON(http.StatusCreated, file)
}

func HandleFilesGetRequest(c echo.Context) error {
	id := c.Param("id")
	result := cache.GetRedisClient().HGet(id, "file")

	var file models.File

	if result.Val() != "" {
		b, err := result.Bytes()
		if err != nil {
			return err
		}
		if err = json.Unmarshal(b, &file); err != nil {
			return err
		}
	} else {
		db := database.GetDB(c.Request().Context())
		err := db.Where("id = ?", id).First(&file).Error

		if err != nil {
			return err
		}

		if file.ID == "" {
			return errors.NewResponseByError(errors.ErrFileNotFound)
		}

		result, err := json.Marshal(file)

		if err != nil {
			return errors.NewResponseByError(errors.ErrJsonMarshal)
		}

		cache.GetRedisClient().HSet(id, "file", result)
	}

	return c.JSON(http.StatusOK, file)
}

func HandleFilesDownloadGetRequest(c echo.Context) error {
	id := c.Param("id")
	result := cache.GetRedisClient().HGet(id, "file")

	var file models.File

	if result.Val() != "" {
		b, err := result.Bytes()
		if err != nil {
			return err
		}
		if err = json.Unmarshal(b, &file); err != nil {
			return err
		}
	} else {
		db := database.GetDB(c.Request().Context())
		err := db.Where("id = ?", id).First(&file).Error

		if err != nil {
			return err
		}

		if file.ID == "" {
			return errors.NewResponseByError(errors.ErrFileNotFound)
		}

		result, err := json.Marshal(file)

		if err != nil {
			return errors.NewResponseByError(errors.ErrJsonMarshal)
		}

		cache.GetRedisClient().HSet(id, "file", result)
	}

	return c.File(file.GetSource())
}

func HandleFilesDeleteRequest(c echo.Context) error {
	id := c.Param("id")

	err := database.GetDB(c.Request().Context()).Where("id = ?", id).Delete(models.File{}).Error

	if err != nil {
		return errors.NewResponse(err, http.StatusBadRequest, "Could not delete this file")
	}

	cache.GetRedisClient().HDel(id, "file")

	return c.JSON(http.StatusOK, nil)
}
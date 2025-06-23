package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"kluisz-object-storage/models"
	"context"
	"fmt"
	"io"
)


// Upload File
// @Summary Upload a file to a given bucket
// @Tags files
// @Accept multipart/form-data
// @Produce plain
// @Param bucket path string true "Bucket name"
// @Param file formData file true "File to upload"
// @Success 200 {object} models.UploadFileResponse
// @Failure 400 {object} models.ErrorResponse400
// @Failure 500 {object} models.ErrorResponse500
// @Router /upload/{bucket} [post]
func UploadFile(c *gin.Context) {
	bucket := c.Param("bucket")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse400{
			Code:  http.StatusBadRequest,
			Error: "File missing" + err.Error(),
		})
		return
	}
	defer file.Close()

	client, err := getMinioClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
			Code:  http.StatusInternalServerError,
			Error: "Internal Server Error - Try again in sometime ",
		})
		return
	}

	uploadInfo, err := client.PutObject(context.Background(), bucket, header.Filename, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
			Code:  http.StatusInternalServerError,
			Error: "Upload Failed" + err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, models.UploadFileResponse{
		Message:   "File uploaded successfully",
		File:      header.Filename,
		Size:      uploadInfo.Size,
		Bucket:    bucket,
		ETag:      uploadInfo.ETag,
	})
}

// Download File
// @Summary Download a file from a bucket
// @Tags files
// @Produce octet-stream
// @Param bucket path string true "Bucket name"
// @Param key path string true "Object key"
// @Success 200 {file} file "File downloaded"
// @Failure 404 {object} models.ErrorResponse404
// @Failure 500 {object} models.ErrorResponse500
// @Router /download/{bucket}/{key} [get]
func DownloadFile(c *gin.Context) {
	bucket := c.Param("bucket")
	file := c.Param("file")

	client, err := getMinioClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
			Code:  http.StatusInternalServerError,
			Error: "Internal Server Error - Try again in sometime ",
		})
		return
	}

	object, err := client.GetObject(context.Background(), bucket, file, minio.GetObjectOptions{})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
			Code:  http.StatusInternalServerError,
			Error: "Failed to get file ",
		})
		return
	}
	defer object.Close()

	stat, err := object.Stat()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound,  models.ErrorResponse404{
			Code:  http.StatusNotFound,
			Error: "File not found ",
		})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file))
	c.Header("Content-Type", stat.ContentType)
	c.Stream(func(w io.Writer) bool {
		io.Copy(w, object)
		return false
	})
}

// ListObjects godoc
// @Summary List objects in a bucket
// @Description Lists all object names in a specified bucket
// @Tags objects
// @Produce json
// @Param bucket path string true "Bucket name"
// @Success 200 {object} models.ListObjectsResponse
// @Failure 500 {object} models.ErrorResponse500
// @Router /objects/{bucket} [get]
func ListObjects(c *gin.Context) {
	bucket := c.Param("bucket")

	client, err := getMinioClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
			Code:  http.StatusInternalServerError,
			Error: "Internal Server Error - Try again in sometime ",
		})
		return
	}

	objectCh := client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{
		Recursive: true,
	})

	var objects []string
	for object := range objectCh {
		if object.Err != nil {
			c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
				Code:  http.StatusInternalServerError,
				Error: "Failed to list objects: " + object.Err.Error(),
			})
			return
		}
		objects = append(objects, object.Key)
	}

	c.IndentedJSON(http.StatusOK, models.ListObjectsResponse{Bucket: bucket, Objects: objects})
}

// DeleteObject godoc
// @Summary Delete a file from a bucket
// @Description Deletes a specified file from a given bucket
// @Tags objects
// @Param bucket path string true "Bucket name"
// @Param file path string true "File name"
// @Success 200 {object} models.DeleteObjectResponse
// @Failure 500 {object} models.ErrorResponse500
// @Router /objects/{bucket}/{file} [delete]
func DeleteObject(c *gin.Context) {
	bucket := c.Param("bucket")
	filename := c.Param("file")

	client, err := getMinioClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
			Code:  http.StatusInternalServerError,
			Error: "Internal Server Error - Try again in sometime ",
		})
		return
	}

	err = client.RemoveObject(context.Background(), bucket, filename, minio.RemoveObjectOptions{})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
			Code:  http.StatusInternalServerError,
			Error: "Object "+ err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, models.DeleteObjectResponse{
		Message: "File deleted",
		Bucket:  bucket,
		File:    filename,
	})
}



package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"kluisz-object-storage/config"
	"kluisz-object-storage/models"
	"context"
)

// func getMinioClient() (*minio.Client, error) {
// 	return minio.New(config.Cfg.S3.Endpoint, &minio.Options{
// 		Creds:  credentials.NewStaticV4(config.Cfg.S3.AccessKey, config.Cfg.S3.SecretKey, ""),
// 		Secure: config.Cfg.S3.UseSSL,
// 		Region: config.Cfg.S3.Region,
// 	})
// }

// Create Bucket
// @Summary Create a new S3 bucket
// @Tags buckets
// @Accept json
// @Param request body models.CreateBucketRequest true "Bucket name payload"
// @Success 200 {object} models.BucketResponseC
// @Failure 400 {object} models.ErrorResponse400
// @Failure 500 {object} models.ErrorResponse500
// @Router /bucket [post]
func CreateBucket(c *gin.Context) {
	var req models.CreateBucketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse400{
			Code:  http.StatusBadRequest,
			Error: "Bad Request- Bucket could not be created" + err.Error(),
		})
		return
	}

	client, err := getMinioClient()
	if err != nil {
		//c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Client init failed: " + err.Error()})
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
			Code:  http.StatusInternalServerError,
			Error: "Internal Server Error - Try again in sometime ",
		})
		return
	}

	err = client.MakeBucket(context.Background(), req.BucketName, minio.MakeBucketOptions{Region: config.Cfg.S3.Region})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
			Code:  http.StatusInternalServerError,
			Error: "Bucket creation failed: " + err.Error(),
		})
		return
	}
	c.IndentedJSON(http.StatusOK, models.BucketResponseC{
		Message: "Bucket created",
		Bucket:  req.BucketName,
	})
}
// Delete Bucket
// @Summary Delete an existing S3 bucket
// @Tags buckets
// @Produce json
// @Param bucket path string true "Bucket name"
// @Success 200 {object} models.BucketResponseD "status message"
// @Failure 500 {object} models.ErrorResponse500 "error message"
// @Router /bucket/{bucket} [delete]
func DeleteBucket(c *gin.Context) {
	bucket := c.Param("name")
	client, err := getMinioClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
			Code:  http.StatusInternalServerError,
			Error: "Internal Server Error - Try again in sometime ",
		})
		return
	}

	err = client.RemoveBucket(context.Background(), bucket)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
			Code:  http.StatusInternalServerError,
			Error: "Bucket deletion failed: " + err.Error(),
		})
		return
	}
	c.IndentedJSON(http.StatusOK, models.BucketResponseD{
		Message: "Bucket deleted",
		Bucket:  bucket,
	})
}

// List Buckets
// @Summary List all available S3 buckets
// @Tags buckets
// @Produce json
// @Success 200 {object} models.ListBucketsResponse
// @Failure 500 {object} models.ErrorResponse500 "error message"
// @Router /buckets [get]
func ListBuckets(c *gin.Context) {
	client, err := getMinioClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
			Code:  http.StatusInternalServerError,
			Error: "Internal Server Error - Try again in sometime ",
		})
		return
	}

	buckets, err := client.ListBuckets(context.Background())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
			Code:  http.StatusInternalServerError,
			Error: "Failed to list buckets: " + err.Error(),
		})
		return
	}

	bucketNames := make([]string, len(buckets))
	for i, bucket := range buckets {
		bucketNames[i] = bucket.Name
	}

	c.IndentedJSON(http.StatusOK, models.ListBucketsResponse{
		Buckets: bucketNames,
    })
}

// // Upload File
// // @Summary Upload a file to a given bucket
// // @Tags files
// // @Accept multipart/form-data
// // @Produce plain
// // @Param bucket path string true "Bucket name"
// // @Param file formData file true "File to upload"
// // @Success 200 {object} UploadFileResponse
// // @Failure 400 {object} models.ErrorResponse400
// // @Failure 500 {object} models.ErrorResponse500
// // @Router /upload/{bucket} [post]
// func UploadFile(c *gin.Context) {
// 	bucket := c.Param("bucket")

// 	file, header, err := c.Request.FormFile("file")
// 	if err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse400{
// 			Code:  http.StatusBadRequest,
// 			Error: "File missing" + err.Error(),
// 		})
// 		return
// 	}
// 	defer file.Close()

// 	client, err := getMinioClient()
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
// 			Code:  http.StatusInternalServerError,
// 			Error: "Internal Server Error - Try again in sometime ",
// 		})
// 		return
// 	}

// 	uploadInfo, err := client.PutObject(context.Background(), bucket, header.Filename, file, header.Size, minio.PutObjectOptions{
// 		ContentType: header.Header.Get("Content-Type"),
// 	})
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
// 			Code:  http.StatusInternalServerError,
// 			Error: "Upload Failed" + err.Error(),
// 		})
// 		return
// 	}

// 	c.IndentedJSON(http.StatusOK, UploadFileResponse{
// 		Message:   "File uploaded successfully",
// 		File:      header.Filename,
// 		Size:      uploadInfo.Size,
// 		Bucket:    bucket,
// 		ETag:      uploadInfo.ETag,
// 	})
// }

// // Download File
// // @Summary Download a file from a bucket
// // @Tags files
// // @Produce octet-stream
// // @Param bucket path string true "Bucket name"
// // @Param key path string true "Object key"
// // @Success 200 {file} file "File downloaded"
// // @Failure 404 {object} models.ErrorResponse404
// // @Failure 500 {object} models.ErrorResponse500
// // @Router /download/{bucket}/{key} [get]
// func DownloadFile(c *gin.Context) {
// 	bucket := c.Param("bucket")
// 	file := c.Param("file")

// 	client, err := getMinioClient()
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
// 			Code:  http.StatusInternalServerError,
// 			Error: "Internal Server Error - Try again in sometime ",
// 		})
// 		return
// 	}

// 	object, err := client.GetObject(context.Background(), bucket, file, minio.GetObjectOptions{})
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
// 			Code:  http.StatusInternalServerError,
// 			Error: "Failed to get file ",
// 		})
// 		return
// 	}
// 	defer object.Close()

// 	stat, err := object.Stat()
// 	if err != nil {
// 		c.IndentedJSON(http.StatusNotFound,  models.ErrorResponse404{
// 			Code:  http.StatusNotFound,
// 			Error: "File not found ",
// 		})
// 		return
// 	}

// 	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file))
// 	c.Header("Content-Type", stat.ContentType)
// 	c.Stream(func(w io.Writer) bool {
// 		io.Copy(w, object)
// 		return false
// 	})
// }

// // ListObjects godoc
// // @Summary List objects in a bucket
// // @Description Lists all object names in a specified bucket
// // @Tags objects
// // @Produce json
// // @Param bucket path string true "Bucket name"
// // @Success 200 {object} ListObjectsResponse
// // @Failure 500 {object} models.ErrorResponse500
// // @Router /objects/{bucket} [get]
// func ListObjects(c *gin.Context) {
// 	bucket := c.Param("bucket")

// 	client, err := getMinioClient()
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
// 			Code:  http.StatusInternalServerError,
// 			Error: "Internal Server Error - Try again in sometime ",
// 		})
// 		return
// 	}

// 	objectCh := client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{
// 		Recursive: true,
// 	})

// 	var objects []string
// 	for object := range objectCh {
// 		if object.Err != nil {
// 			c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
// 				Code:  http.StatusInternalServerError,
// 				Error: "Failed to list objects: " + object.Err.Error(),
// 			})
// 			return
// 		}
// 		objects = append(objects, object.Key)
// 	}

// 	c.IndentedJSON(http.StatusOK, ListObjectsResponse{Bucket: bucket, Objects: objects})
// }

// // DeleteObject godoc
// // @Summary Delete a file from a bucket
// // @Description Deletes a specified file from a given bucket
// // @Tags objects
// // @Param bucket path string true "Bucket name"
// // @Param file path string true "File name"
// // @Success 200 {object} DeleteObjectResponse
// // @Failure 500 {object} models.ErrorResponse500
// // @Router /objects/{bucket}/{file} [delete]
// func DeleteObject(c *gin.Context) {
// 	bucket := c.Param("bucket")
// 	filename := c.Param("file")

// 	client, err := getMinioClient()
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse500{
// 			Code:  http.StatusInternalServerError,
// 			Error: "Internal Server Error - Try again in sometime ",
// 		})
// 		return
// 	}

// 	err = client.RemoveObject(context.Background(), bucket, filename, minio.RemoveObjectOptions{})
// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, ErrorResponse500{
// 			Code:  http.StatusInternalServerError,
// 			Error: "Object "+ err.Error(),
// 		})
// 		return
// 	}

// 	c.IndentedJSON(http.StatusOK, DeleteObjectResponse{
// 		Message: "File deleted",
// 		Bucket:  bucket,
// 		File:    filename,
// 	})
// }



package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"kluisz-object-storage/config"
	"kluisz-object-storage/models"
	"context"
)


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


package models


// ErrorResponse -- error message with code
type ErrorResponse500 struct {
	Code  int    `json:"code" example:"500"`
	Error string `json:"error" example:"Internal Server Error message"`
}
type ErrorResponse404 struct {
	Code  int    `json:"code" example:"404"`
	Error string `json:"error" example:"Not Found Error message"`
}
type ErrorResponse400 struct {
	Code  int    `json:"code" example:"400"`
	Error string `json:"error" example:"Bad request Error message"`
}

type CreateBucketRequest struct {
	BucketName string `json:"bucketName" example:"mybucket"`
}

// response to create bucket
type BucketResponseC struct {
	Message string `json:"message" example:"Bucket created"`
	Bucket  string `json:"bucket" example:"mybucket"`
}

// response to delete bucket
type BucketResponseD struct {
	Message string `json:"message" example:"Bucket deleted"`
	Bucket  string `json:"bucket" example:"mybucket"`
}

type ListBucketsResponse struct {
	Buckets []string `json:"buckets"`
}


type UploadFileResponse struct {
	Message string `json:"message" example:"File uploaded successfully"`
	File    string `json:"file" example:"file.txt"`
	Size    int64  `json:"size" example:"1234"`
	Bucket  string `json:"bucket" example:"mybucket"`
	ETag    string `json:"etag" example:"abcd1234"`
}

type ListObjectsResponse struct {
	Bucket  string   `json:"bucket" example:"mybucket"`
	Objects []string `json:"objects"`
}

type DeleteObjectResponse struct {
	Message string `json:"message" example:"File deleted"`
	Bucket  string `json:"bucket" example:"mybucket"`
	File    string `json:"file" example:"file.txt"`
}

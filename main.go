package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const PORT = 8080
var s3Client *s3.S3


func createToken(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Get the token from the query
	s3Path := query.Get("s3")
	separatorIndex := strings.Index(s3Path, "/")

	// Check if path seems to be valid
	if separatorIndex == -1 || separatorIndex == 0 || separatorIndex == len(s3Path) - 1 {
		log.Println("Invalid S3 path")
		fmt.Fprintf(w, "Invalid S3 path")
		return
	}

	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
        Bucket: aws.String(s3Path[:separatorIndex]),
        Key:    aws.String(s3Path[separatorIndex+1:]),
    })

    urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Println("Failed to sign request", err)
		fmt.Fprintf(w, "Failed to sign request")
		return
	}

	fmt.Fprintf(w, "%s", urlStr)
	log.Printf("Created token for %s", s3Path)
}

func main() {
	// Create S3 service client
	sess, _ := session.NewSession(&aws.Config{
        Region: aws.String("us-east-1"),
	})

    // Create S3 service client
    s3Client = s3.New(sess)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request){
        fmt.Fprintf(w, "pong")
    })

    http.HandleFunc("/token/create", createToken)

	fmt.Printf("Server is running on port %d\n", PORT)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
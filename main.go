package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const PORT = 8080
var s3Client *s3.S3


func createToken(path string, ttlMinutes int) string {
	separatorIndex := strings.Index(path, "/")

	// Check if path seems to be valid
	if separatorIndex == -1 || separatorIndex == 0 || separatorIndex == len(path) - 1 {
		log.Println("Invalid S3 path")
		return "invalid_path"
	}

	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
        Bucket: aws.String(path[:separatorIndex]),
        Key:    aws.String(path[separatorIndex+1:]),
    })

    urlStr, err := req.Presign(time.Duration(ttlMinutes))
	if err != nil {
		log.Println("Failed to sign request", err)
		return "signing_failed"
	}

	return urlStr
}

func createTokenEndpoint(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Get the parameters
	s3Path := query.Get("s3")
	ttlMinutes := query.Get("ttl")

	// Check if parameters are valid
	if s3Path == "" {
		log.Println("Missing path parameter")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "missing_s3_path")
		return
	}

	// Convert ttlMinutes to int
	ttlMinutesInt, err := strconv.Atoi(ttlMinutes)
	if err != nil {
		log.Println("Failed to convert to int:", ttlMinutes)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid_ttl")
		return
	}
	
	response := createToken(s3Path, ttlMinutesInt)

	fmt.Fprintf(w, "%s", response)
	log.Printf("Created token for %s", s3Path)
}

func pingPongEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func main() {
	// Create S3 service client
	sess, _ := session.NewSession(&aws.Config{
        Region: aws.String("us-east-1"),
	})
    s3Client = s3.New(sess)

	// Register endpoints
	http.HandleFunc("/ping", pingPongEndpoint)
    http.HandleFunc("/s3/sign", createTokenEndpoint)

	// Start server
	fmt.Printf("Server is running on port %d\n", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
    log.Fatal(err)
}
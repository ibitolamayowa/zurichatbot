package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"os"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	// "zuri.chat/zccore/auth"/
)

type M map[string]interface{}

// ErrorResponse : This is error model.
type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// DetailedErrorResponse : This is success model.
type DetailedErrorResponse struct {
	StatusCode int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// SuccessResponse : This is success model.
type SuccessResponse struct {
	StatusCode int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// GetError : This is helper function to prepare error model.
func GetError(err error, StatusCode int, w http.ResponseWriter) {
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   StatusCode,
	}

	w.WriteHeader(response.StatusCode)
	w.Header().Set("Content-Type", "application/json<Left>")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

// GetDetailedError: This function provides detailed error information
func GetDetailedError(msg string, StatusCode int, data interface{}, w http.ResponseWriter) {
	var response = DetailedErrorResponse{
		Message:    msg,
		StatusCode: StatusCode,
		Data:       data,
	}

	w.WriteHeader(response.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

// GetSuccess : This is helper function to prepare success model.
func GetSuccess(msg string, data interface{}, w http.ResponseWriter) {
	var response = SuccessResponse{
		Message:    msg,
		StatusCode: http.StatusOK,
		Data:       data,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

// get env vars; return empty string if not found
func Env(key string) string {
	return os.Getenv(key)
}

// check if a file exists, useful in checking for .env
func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// convert map to bson.M for mongoDB docs
func MapToBson(data map[string]interface{}) bson.M {
	return bson.M(data)
}

// StructToMap converts a struct of any type to a map[string]inteface{}
func StructToMap(inStruct interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	inrec, _ := json.Marshal(inStruct)
	json.Unmarshal(inrec, &out)
	return out, nil
}

// ConvertStructure does map to struct conversion and vice versa.
// The input structure will be converted to the output
func ConvertStructure(input interface{}, output interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, output)
}

func ParseJsonFromRequest(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GenUUID() string {
	id := uuid.New()
	return id.String()
}

// Check the validaity of a UUID. Returns a valid UUID from a string input. Returns an error otherwise
func ValidateUUID(s string) (uuid.UUID, error) {
	if len(s) != 36 {
		return uuid.Nil, errors.New("invalid uuid format")
	}

	b, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, err
	}

	return b, nil
}

func ConvertImageTo64(ImgDirectory string) string {
	// Read the entire file into a byte slice
	bytes, err := ioutil.ReadFile(ImgDirectory)
	if err != nil {
		log.Fatal(err)
	}

	var base64Encoding string

	// Determine the content type of the image file
	// mimeType := http.DetectContentType(bytes)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	// switch mimeType {
	// case "image/jpeg":
	// 	base64Encoding += "data:image/jpeg;base64,"
	// case "image/png":
	// 	base64Encoding += "data:image/png;base64,"
	// }

	// Append the base64 encoded output
	base64Encoding += base64.StdEncoding.EncodeToString(bytes)

	// Print the full base64 representation of the image
	return base64Encoding
}

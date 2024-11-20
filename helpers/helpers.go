package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Envelope map[string]interface{}

type Message struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

var MessageLogs = &Message{
    InfoLog: infoLog,
    ErrorLog: errorLog,
}

// ReadJSON - helper method that reads incoming json from http requests
func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
    maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

// WriteJSON - method to write data into a json response 
func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
    out, err := json.MarshalIndent(data, "", "\t")
    if err != nil {
        return err
    }

    if len(headers) > 0 {
        // headers[0] is a map of keys and values
        for key, value := range headers[0] {
            w.Header()[key] = value
        }
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    _, err = w.Write(out)
    if err != nil {
        return err
    }
    return nil
}

// WriteERROR - Method that writes error back to the api in json format
func WriteERROR(w http.ResponseWriter, status int, err error) {
    WriteJSON(w, status, map[string]string{"error": err.Error()})
}

// convertURLToBase64ID - will grab any characters between min(16 str.length) 
// and convert that to base64, a deterministic way to create a uniqueID by link URL
func ConvertURLToBase64ID(url string) string {
    urlLen := len(url)
    mid := Min(urlLen, 16)
    encodedString := base64.StdEncoding.EncodeToString([]byte(url))
    return encodedString[urlLen-mid:urlLen]
}

func Min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

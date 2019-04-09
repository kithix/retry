package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/kithix/retry"
)

func statusWriter(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}))
}

func statusCodeRetrier(err error) bool {
	if strings.Contains(err.Error(), "500") {
		fmt.Println("Error received:", err, "Retrying")
		return true
	}
	fmt.Println("Error received:", err, "Skipping")
	return false
}

func serverRetriever(url string) func() error {
	return func() error {
		fmt.Println("Trying", url)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		// Lazy status code check
		if resp.StatusCode >= 500 {
			return errors.New("Received status code >= 500")
		}
		if resp.StatusCode >= 400 {
			return errors.New("Recevied status code >= 400")
		}
		return nil
	}
}

func main() {
	OKWriter := statusWriter(http.StatusOK)
	ForbiddenWriter := statusWriter(http.StatusForbidden)
	ServerErrorWriter := statusWriter(http.StatusInternalServerError)

	// This will try 5 times, we want to retry server errors.
	fmt.Println("Error:", retry.Do(
		serverRetriever(ServerErrorWriter.URL),
		retry.WithLimit(statusCodeRetrier, 5),
	), "\n")

	// This will exit after the first attempt, we skip forbidden errors.
	fmt.Println("Error:", retry.Do(
		serverRetriever(ForbiddenWriter.URL),
		retry.WithLimit(statusCodeRetrier, 5),
	), "\n")

	// This will exit after the first attempt, success!
	err := (retry.Do(
		serverRetriever(OKWriter.URL),
		retry.WithLimit(statusCodeRetrier, 5),
	))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("No error, we did it!")
}

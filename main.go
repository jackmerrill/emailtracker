package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"os"
	"time"

	"github.com/jackmerrill/emailtracker/database"
)

type TrackingData struct {
	Opened        bool      `json:"opened"`
	FirstOpenedAt time.Time `json:"first_opened_at"`
	OpenedAmount  int       `json:"opened_amount"`
}

type EncodedData struct {
	To      string    `json:"to,omitempty"`
	From    string    `json:"from,omitempty"`
	Subject string    `json:"subject,omitempty"`
	Date    time.Time `json:"date,omitempty"`
}

type ReturnedData struct {
	Tracking    TrackingData `json:"tracking"`
	DecodedData EncodedData  `json:"decoded_data,omitempty"`
}

func main() {
	db, err := database.NewDatabase("db/db.json")

	if err != nil {
		panic(err)
	}

	// ID could be Base64 encoded JSON of the email's to, from, subject, and date
	http.HandleFunc("/i/", func(w http.ResponseWriter, r *http.Request) {
		img := image.NewAlpha(image.Rect(0, 0, 1, 1))

		record := TrackingData{}

		exists := db.Exists(r.URL.Query().Get("id"))

		if exists {
			err := db.Get(r.URL.Query().Get("id"), &record)

			if err != nil {
				panic(err)
			}
		} else {
			record.FirstOpenedAt = time.Now()
		}

		record.Opened = true
		record.OpenedAmount++

		err = db.Set(r.URL.Query().Get("id"), record)

		if err != nil {
			panic(err)
		}

		png.Encode(w, img)
	})

	// Generate a token for the email
	http.HandleFunc("/t", func(w http.ResponseWriter, r *http.Request) {
		// preflight request
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.Write([]byte("OK"))
			return
		}

		data := EncodedData{
			To:      r.URL.Query().Get("to"),
			From:    r.URL.Query().Get("from"),
			Subject: r.URL.Query().Get("subject"),
			Date:    time.Now(),
		}

		jsonData, err := json.Marshal(data)

		if err != nil {
			panic(err)
		}

		encodedData := base64.StdEncoding.EncodeToString(jsonData)

		fmt.Fprintf(w, encodedData)
	})

	http.HandleFunc("/panel", func(w http.ResponseWriter, r *http.Request) {
		// Request username and password from user (use env vars)
		// If correct, show the panel

		username, password, ok := r.BasicAuth()

		if !ok || username != os.Getenv("AUTH_USER") || password != os.Getenv("AUTH_PASS") {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Unauthorized")
			return
		}

		data := map[string]TrackingData{}
		err := db.GetAll(&data)

		if err != nil {
			panic(err)
		}

		returnedData := []ReturnedData{}

		for id, trackingData := range data {
			returnedData = append(returnedData, ReturnedData{
				Tracking: trackingData,
			})

			// Decode the base64 encoded JSON

			base64Data, err := base64.StdEncoding.DecodeString(id)

			if err != nil {
				fmt.Println("not base64, skipping")
				continue
			}

			var decodedData EncodedData

			err = json.Unmarshal(base64Data, &decodedData)

			if err != nil {
				fmt.Println("not json, skipping")
				continue
			} else {
				returnedData[len(returnedData)-1].DecodedData = decodedData
			}
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(returnedData)
	})

	fmt.Println("Listening on port 8080")

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}

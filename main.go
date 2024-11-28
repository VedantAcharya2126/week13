package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type TimeResponse struct {
	CurrentTime string `json:"current_time"`
}

func initDB() error {
	var err error
	db, err = sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/week13Lab")
	if err != nil {
		return err
	}
	return db.Ping()
}

// Get current time in Toronto timezone
func getCurrentTorontoTime() (time.Time, error) {
	location, err := time.LoadLocation("America/Toronto")
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(location), nil
}

// API handler for /current-time
func currentTimeHandler(w http.ResponseWriter, r *http.Request) {
	// Get Toronto time
	torontoTime, err := getCurrentTorontoTime()
	if err != nil {
		http.Error(w, "Error getting Toronto time", http.StatusInternalServerError)
		log.Printf("Timezone error: %v", err)
		return
	}

	// Insert time into MySQL
	_, err = db.Exec("INSERT INTO time_log (timestamp) VALUES (?)", torontoTime)
	if err != nil {
		http.Error(w, "Error logging time to database", http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}

	// Respond with JSON
	response := TimeResponse{CurrentTime: torontoTime.Format(time.RFC3339)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Initialize database
	if err := initDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Set up HTTP server
	http.HandleFunc("/current-time", currentTimeHandler)
	fmt.Println("Server is running on http://localhost:80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

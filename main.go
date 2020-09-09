package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "time"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

/*
 * Models
 */
type Event struct {
  ID        int       `json:"id" gorm:"primaryKey"`
  CreatedAt time.Time `json:"Created_at"`
  Level     uint      `json:"level"`
  Title     string    `json:"title"`
  Traceback string    `json:"traceback"`
  ProjectID int       `json:"project_id"`
}

type Project struct {
  ID    uint    `json:"id" gorm:"primaryKey"`
  Name  string  `json:"name"`
  Event []Event `json:"-" gorm:"foreignkey:ProjectID"`
}

/*
 * Database management
 */
var db *gorm.DB

func initDB() {
  var err error
  db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
  if err != nil {
      fmt.Println(err)
      panic("failed to connect database")
  }

  db.AutoMigrate(&Project{}, &Event{})
}

/*
 * API
 */
func createProject(w http.ResponseWriter, r *http.Request) {
  var project Project
  err := json.NewDecoder(r.Body).Decode(&project)
  if err != nil {
      http.Error(w, "Bad JSON data", http.StatusBadRequest)
      return
  }

  db.Create(&project)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(project)
}

func getProjects(w http.ResponseWriter, r *http.Request) {
  var projects []Project
  db.Find(&projects)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(projects)
}

func getProject(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  projectID := params["projectID"]

  var project Project
  db.Preload("Event").First(&project, projectID)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(project)
}

func createEvent(w http.ResponseWriter, r *http.Request) {
  var event Event
  err := json.NewDecoder(r.Body).Decode(&event)
  if err != nil {
      http.Error(w, "Bad JSON data", http.StatusBadRequest)
      return
  }

  db.Create(&event)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(event)
}

func getEvents(w http.ResponseWriter, r *http.Request) {
  var events []Event
  db.Find(&events)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(events)
}

func getEvent(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  eventID := params["eventID"]

  var event Event
  db.First(&event, eventID)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(event)
}

/*
 * Main Snitch 
 */
func main() {
  // The naming of StrictSlash is bad -- true/false are inverted
  // See: https://github.com/gorilla/mux/issues/145
  r := mux.NewRouter().StrictSlash(true)
  ar := r.PathPrefix("/api").Subrouter()
  ar.HandleFunc("/projects", createProject).Methods("POST")
  ar.HandleFunc("/projects", getProjects).Methods("GET")
  ar.HandleFunc("/projects/{projectID}", getProject).Methods("GET")
  ar.HandleFunc("/events", createEvent).Methods("POST")
  ar.HandleFunc("/events", getEvents).Methods("GET")
  ar.HandleFunc("/events/{eventID}", getEvent).Methods("GET")

  initDB()

  srv := &http.Server{
		Handler: r,
		Addr: "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

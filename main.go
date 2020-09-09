package main

import (
	"fmt"
	"net/http"
  "time"

  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
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
  Event []Event `json:"events,omitempty" gorm:"foreignkey:ProjectID"`
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
func createProject(c echo.Context) (err error) {
  project := new(Project)
  if err = c.Bind(project); err != nil {
    return
  }
  db.Create(&project)
  return c.JSON(http.StatusOK, project)
}

func getProjects(c echo.Context) error {
  var projects []Project
  db.Find(&projects)
  return c.JSON(http.StatusOK, projects)
}

func getProject(c echo.Context) error {
  projectID := c.Param("id")

  var project Project
  db.Preload("Event").First(&project, projectID)
  return c.JSON(http.StatusOK, project)
}

func createEvent(c echo.Context) (err error) {
  event := new(Event)
  if err = c.Bind(event); err != nil {
    return
  }
  db.Create(&event)
  return c.JSON(http.StatusOK, event)
}

func getEvents(c echo.Context) error {
  var events []Event
  db.Find(&events)
  return c.JSON(http.StatusOK, events)
}

func getEvent(c echo.Context) error {
  eventID := c.Param("id")

  var event Event
  db.First(&event, eventID)
  return c.JSON(http.StatusOK, event)
}

/*
 * Main Snitch 
 */
func main() {
  e := echo.New()
  e.Pre(middleware.RemoveTrailingSlash())
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  a := e.Group("/api")
  a.POST("/projects", createProject)
  a.GET("/projects", getProjects)
  a.GET("/projects/:id", getProject)
  a.POST("/events", createEvent)
  a.GET("/events", getEvents)
  a.GET("/events/:id", getEvent)

  initDB()

  e.Logger.Fatal(e.Start(":8000"))
}

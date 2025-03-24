package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nixoncode/gallery_api/internal/database"
	"github.com/nixoncode/gallery_api/internal/handlers"
	"github.com/nixoncode/gallery_api/internal/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Error Handler
	e.HTTPErrorHandler = customHTTPErrorHandler

	baseDir := "uploads"
	os.MkdirAll(baseDir, 0755)

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	maxFileSizeStr := os.Getenv("MAX_FILE_SIZE")
	maxFileSize, err := strconv.ParseInt(maxFileSizeStr, 10, 64)
	if err != nil {
		log.Fatalf("failed to parse MAX_FILE_SIZE: %v", err)
	}

	storage := storage.NewStorage(baseDir, db)
	handlers := handlers.NewHandlers(storage, maxFileSize) // 10MB limit

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, Response{Message: "OK"})
	})
	e.POST("/upload", handlers.UploadImage)
	e.GET("/image", handlers.GetImage)
	e.GET("/image/file", handlers.GetImageFile)
	e.GET("/images", handlers.GetImageDetails)

	e.Logger.Fatal(e.Start("0.0.0.0:6661"))
}

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := "Internal Server Error"
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message.(string)
	}
	if err := c.JSON(code, Response{Error: msg}); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}

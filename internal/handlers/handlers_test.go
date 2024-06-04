package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"service-catalog/internal/db"
	"service-catalog/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testClient *mongo.Client

func setupTestDB() *mongo.Client {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("No .env file found")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		panic("MONGODB_URI environment variable is not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	return client
}

func clearTestCollection() {
	collection := db.GetCollection(testClient, "services")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.Drop(ctx)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	testClient = setupTestDB()
	code := m.Run()
	clearTestCollection()
	os.Exit(code)
}

func TestListServices(t *testing.T) {
	clearTestCollection()
	collection := db.GetCollection(testClient, "services")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection.InsertOne(ctx, models.Service{
		ID:          "1",
		Name:        "Test Service 1",
		Description: "Description 1",
		Versions:    []string{"1.0"},
		CreatedAt:   time.Now(),
	})
	collection.InsertOne(ctx, models.Service{
		ID:          "2",
		Name:        "Test Service 2",
		Description: "Description 2",
		Versions:    []string{"1.0"},
		CreatedAt:   time.Now(),
	})

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", testClient)
	})
	r.GET("/services", ListServices)

	req, _ := http.NewRequest("GET", "/services", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(response["data"].([]interface{})))
}

func TestGetService(t *testing.T) {
	clearTestCollection()
	collection := db.GetCollection(testClient, "services")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection.InsertOne(ctx, models.Service{
		ID:          "1",
		Name:        "Test Service 1",
		Description: "Description 1",
		Versions:    []string{"1.0"},
		CreatedAt:   time.Now(),
	})

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", testClient)
	})
	r.GET("/services/:id", GetService)

	req, _ := http.NewRequest("GET", "/services/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response models.Service
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Test Service 1", response.Name)
}

func TestCreateService(t *testing.T) {
	clearTestCollection()
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", testClient)
	})
	r.POST("/services", CreateService)

	service := models.Service{
		Name:        "New Test Service",
		Description: "New Description",
		Versions:    []string{"1.0"},
	}
	jsonValue, _ := json.Marshal(service)
	req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response models.Service
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "New Test Service", response.Name)
}

func TestUpdateService(t *testing.T) {
	clearTestCollection()
	collection := db.GetCollection(testClient, "services")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection.InsertOne(ctx, models.Service{
		ID:          "1",
		Name:        "Test Service 1",
		Description: "Description 1",
		Versions:    []string{"1.0"},
		CreatedAt:   time.Now(),
	})

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", testClient)
	})
	r.PUT("/services/:id", UpdateService)

	updatedService := models.Service{
		Name:        "Updated Test Service 1",
		Description: "Updated Description 1",
		Versions:    []string{"1.1"},
	}
	jsonValue, _ := json.Marshal(updatedService)
	req, _ := http.NewRequest("PUT", "/services/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Service updated successfully", response["message"])
}

func TestDeleteService(t *testing.T) {
	clearTestCollection()
	collection := db.GetCollection(testClient, "services")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection.InsertOne(ctx, models.Service{
		ID:          "1",
		Name:        "Test Service 1",
		Description: "Description 1",
		Versions:    []string{"1.0"},
		CreatedAt:   time.Now(),
	})

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", testClient)
	})
	r.DELETE("/services/:id", DeleteService)

	req, _ := http.NewRequest("DELETE", "/services/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Service deleted successfully", response["message"])
}

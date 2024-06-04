package handlers

import (
	"context"
	"net/http"
	"service-catalog/internal/db"
	"service-catalog/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListServices(c *gin.Context) {
	client := c.MustGet("db").(*mongo.Client)
	collection := db.GetCollection(client, "services")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	sortField := c.DefaultQuery("sortField", "created_at")
	sortOrder, err := strconv.Atoi(c.DefaultQuery("sortOrder", "1"))
	if err != nil || (sortOrder != 1 && sortOrder != -1) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort order"})
		return
	}

	if sortField != "name" && sortField != "created_at" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort field"})
		return
	}

	searchQuery := c.DefaultQuery("search", "")

	filter := bson.M{}
	if searchQuery != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": searchQuery, "$options": "i"}},
				{"description": bson.M{"$regex": searchQuery, "$options": "i"}},
			},
		}
	}

	skip := (page - 1) * pageSize

	cursor, err := collection.Find(ctx, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.D{{Key: sortField, Value: sortOrder}}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var services []models.Service
	if err := cursor.All(ctx, &services); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"data":     services,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	}

	c.JSON(http.StatusOK, response)
}

func GetService(c *gin.Context) {
	client := c.MustGet("db").(*mongo.Client)
	collection := db.GetCollection(client, "services")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	var service models.Service
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&service)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, service)
}

func CreateService(c *gin.Context) {
	var service models.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service.ID = uuid.New().String()
	service.CreatedAt = time.Now()

	client := c.MustGet("db").(*mongo.Client)
	collection := db.GetCollection(client, "services")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

func DeleteService(c *gin.Context) {
	id := c.Param("id")
	client := c.MustGet("db").(*mongo.Client)
	collection := db.GetCollection(client, "services")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

func UpdateService(c *gin.Context) {
	id := c.Param("id")
	var service models.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := c.MustGet("db").(*mongo.Client)
	collection := db.GetCollection(client, "services")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"name":        service.Name,
			"description": service.Description,
			"versions":    service.Versions,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service updated successfully"})
}

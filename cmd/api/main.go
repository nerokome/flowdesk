package main

import (
    "log"
    "net/http"
    "os"
    "flowdesk/internal/repository"
    "flowdesk/internal/routes" 
    "github.com/gin-gonic/gin" 
    "github.com/joho/godotenv"
)

func main() {
    godotenv.Load()

    db := repository.ConnectDB()
    sqlDB, _ := db.DB()
    if err := sqlDB.Ping(); err != nil {
        log.Fatal("DB Ping Failed")
    }
    log.Println("🚀 Database Connected!")

    r := gin.Default()

    // 2. Load the routes here
    routes.SetupRoutes(r, db)

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Flowdesk Server is running!"})
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("🔥 Server starting on port %s", port)
    r.Run(":" + port)
}
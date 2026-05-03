package main

import (
	"flag"                     
	"flowdesk/internal/models" 
	"flowdesk/internal/repository"
	"flowdesk/internal/routes"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	resetDB := flag.Bool("reset", false, "Reset the database")
	flag.Parse()

	db := repository.ConnectDB()

	if *resetDB {
		log.Println("⚠️  Resetting database...")
		db.Migrator().DropTable(&models.User{})
		db.AutoMigrate(&models.User{})
		log.Println("✅ Database reset complete. Exiting...")
		return // Stop here so it doesn't start the server
	}

	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("DB Ping Failed")
	}
	log.Println("🚀 Database Connected!")

	r := gin.Default()
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

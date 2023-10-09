package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type apiConfig struct {
	DB *gorm.DB
}

func connectToDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// Define a function to generate a random hexadecimal value
func generateRandomHexValue() string {
	// Generate a random byte array
	randomBytes := make([]byte, 20)
	_, _ = rand.Read(randomBytes)

	// Convert the byte array to a hexadecimal string
	return fmt.Sprintf("%x", randomBytes)
}

func main() {
	//? Loading environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	//? Getting port
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	//? Setting up the database
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	db, err := connectToDB(dsn)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	fmt.Println("Database Connected")

	//? Setting up automigrations with GORM
	err = db.AutoMigrate(&Users{}, &Feeds{}, &FeedFollows{}, &Posts{})
	if err != nil {
		log.Fatal("Failed to auto migrate models")
	}

	// Execute the raw SQL migration
	migrator := db.Migrator()

	// Check if the api_key column already exists
	if !migrator.HasColumn(&Users{}, "api_key") {
		// Generate a random hexadecimal string
		randomHexValue := generateRandomHexValue()

		// Define the ALTER TABLE statement with the default value
		sqlQuery := fmt.Sprintf("ALTER TABLE `server1`.`users` ADD COLUMN `api_key` VARCHAR(255) NOT NULL DEFAULT '%v' AFTER `name`, ADD UNIQUE INDEX `api_key_UNIQUE` (`api_key` ASC) VISIBLE;", randomHexValue)

		// Execute the query with the random hexadecimal value
		err = db.Exec(sqlQuery).Error
		if err != nil {
			panic("failed to execute migration")
		}
	}
	// Check if the api_key column already exists
	if !migrator.HasColumn(&Feeds{}, "last_fetched_at") {

		// Define the ALTER TABLE statement with the default value
		sqlQuery := "ALTER TABLE `server1`.`feeds` ADD COLUMN `last_fetched_at` TIMESTAMP"

		// Execute the query with the random hexadecimal value
		err = db.Exec(sqlQuery).Error
		if err != nil {
			panic("failed to execute migration")
		}
	}

	apiCfg := apiConfig{
		DB: db,
	}

	go startScrapping(db, 10, time.Minute)

	//? Routes
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	//? Routers
	v1Router := chi.NewRouter()
	router.Mount("/v1", v1Router)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollows))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollows))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	//? Starting Server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server started on port %s\n", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error in starting the server: ", err)
	}
}

package main

import (
	"database/sql"
	"fmt"
	"go-server/internal"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	logger := internal.Logger{}
	logger.LoggerInit()
	db_url := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", os.Getenv("HOST"), 5432, os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DBNAME"))

	// connect to database
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		logger.LogError(err.Error())
		return
	}

	defer db.Close()

	// ping database
	err = db.Ping()
	if err != nil {
		logger.LogError(err.Error())
		return
	}

	logger.LogInfo("connected to database successfully.")

	// connect to pytorch grpc server
	conn, err := grpc.Dial(os.Getenv("PYTORCH"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.LogError(err.Error())
		return
	}

	defer conn.Close()

	logger.LogInfo("connected to pytorch grpc server successfully.")

	var (
		repository = internal.NewRepository(db)
		service    = internal.NewService(repository, conn, logger)
		handler    = internal.NewHandler(service, logger)
	)

	// implement gin server
	r := gin.Default()

	r.POST("/embed", handler.EmbedTexts)
	r.GET("/search", handler.PerformSemanticSearch)

	r.Run(os.Getenv("PORT"))
}

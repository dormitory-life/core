package main

import (
	"log"
	"os"

	"github.com/dormitory-life/core/internal/config"
	"github.com/dormitory-life/core/internal/database"
	"github.com/dormitory-life/core/internal/logger"
	"github.com/dormitory-life/core/internal/server"
	core "github.com/dormitory-life/core/internal/service"
)

func main() {
	configPath := os.Args[1]
	cfg, err := config.ParseConfig(configPath)
	if err != nil {
		panic(err)
	}

	log.Println("CONFIG: ", cfg)

	logger, err := logger.New(cfg)
	if err != nil {
		panic(err)
	}

	db, err := database.InitDb(cfg.Db)
	if err != nil {
		panic(err)
	}

	repository := database.New(db)

	coreService := core.New(core.CoreServiceConfig{
		Repository: repository,
	})

	s := server.New(server.ServerConfig{
		Config:      cfg.Server,
		CoreService: coreService,
		Logger:      logger,
	})

	panic(s.Start())
}

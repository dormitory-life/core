package main

import (
	"log"
	"os"

	"github.com/dormitory-life/core/internal/auth"
	"github.com/dormitory-life/core/internal/config"
	"github.com/dormitory-life/core/internal/database"
	"github.com/dormitory-life/core/internal/logger"
	"github.com/dormitory-life/core/internal/server"
	core "github.com/dormitory-life/core/internal/service"
	"github.com/dormitory-life/core/internal/storage"
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

	log.Println("Core init db...")
	db, err := database.InitDb(cfg.Db)
	if err != nil {
		panic(err)
	}

	repository := database.New(db)

	authClient, err := auth.New(auth.AuthClientConfig{
		GRPCAuthServerAddress: cfg.Auth.AuthClientConfig.GRPCAuthServerAddress,
		Timeout:               cfg.Auth.AuthClientConfig.Timeout,
		Logger:                *logger,
	})
	if err != nil {
		panic(err)
	}

	s3Client, err := storage.New(storage.S3StorageConfig{
		Type:            cfg.Storage.Type,
		Endpoint:        cfg.Storage.MinIO.Endpoint,
		AccessKeyId:     cfg.Storage.MinIO.AccessKeyId,
		SecretAccessKey: cfg.Storage.MinIO.SecretAccessKey,
		UseSSL:          cfg.Storage.MinIO.UseSSL,
		BucketName:      cfg.Storage.MinIO.BucketName,
		PublicUrl:       cfg.Storage.MinIO.PublicUrl,

		Logger: *logger,
	})
	if err != nil {
		panic(err)
	}

	coreService := core.New(core.CoreServiceConfig{
		Repository: repository,
		AuthClient: authClient,
		Logger:     *logger,
		S3Client:   s3Client,
	})

	s := server.New(server.ServerConfig{
		Config:      cfg.Server,
		CoreService: coreService,
		Logger:      logger,
	})

	panic(s.Start())
}

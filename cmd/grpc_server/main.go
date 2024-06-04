package main

import (
	"context"
	"flag"
	"github.com/dmtrybogdanov/garantex/internal/tracing"
	"log"
	"net"

	ratesAPI "github.com/dmtrybogdanov/garantex/internal/api/rates"
	"github.com/dmtrybogdanov/garantex/internal/config"
	rates2 "github.com/dmtrybogdanov/garantex/internal/repository/rates"
	"github.com/dmtrybogdanov/garantex/internal/service/rates"
	"github.com/dmtrybogdanov/garantex/pkg/rates_v1"
	"github.com/gofiber/contrib/otelfiber"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	_, err := tracing.InitTracer("http://localhost:14268/api/traces", "Rates Service")
	if err != nil {
		log.Fatal("init tracer", err)
	}
	app := fiber.New()
	app.Use(otelfiber.Middleware())

	flag.Parse()
	ctx := context.Background()

	err = config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	ratesRepo := rates2.NewRepository(pool)
	rateSrv := rates.NewService(ratesRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	rates_v1.RegisterRatesV1Server(s, ratesAPI.NewImplementation(rateSrv))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

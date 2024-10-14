package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"muzz-explore-service/pkg/service"
	"muzz-explore-service/pkg/tests"
	"net"
	"os"

	"muzz-explore-service/pkg/protos/generated"
	"muzz-explore-service/pkg/store"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get Redis address from environment variables
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatalf("REDIS_ADDR not set in environment variables")
	}

	// Initialize Redis store
	newStore := store.NewRedisStore(redisAddr)

	// Populate test data
	log.Println("Running Redis test data initialization")
	err = tests.PopulateRedisTestData(newStore.GetClient())
	if err != nil {
		log.Fatalf("Failed to populate Redis test data: %v", err)
	}

	// Set up gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	s := grpc.NewServer()
	exploreService := service.NewExploreService(newStore) // Pass Redis newStore to the service

	// Register the ExploreService with the gRPC server
	generated.RegisterExploreServiceServer(s, exploreService)

	log.Println("gRPC server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

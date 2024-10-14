package service

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"muzz-explore-service/pkg/protos/generated"
	"muzz-explore-service/pkg/store"
	"os"
	"testing"
)

// BenchmarkExploreService_ListLikedYou benchmarks the ListLikedYou method of ExploreService.
func BenchmarkExploreService_ListLikedYou(b *testing.B) {
	benchmarkExploreServiceMethod(b, "ListLikedYou")
}

// BenchmarkExploreService_ListNewLikedYou benchmarks the ListNewLikedYou method of ExploreService.
func BenchmarkExploreService_ListNewLikedYou(b *testing.B) {
	benchmarkExploreServiceMethod(b, "ListNewLikedYou")
}

// BenchmarkExploreService_CountLikedYou benchmarks the CountLikedYou method of ExploreService.
func BenchmarkExploreService_CountLikedYou(b *testing.B) {
	benchmarkExploreServiceMethod(b, "CountLikedYou")
}

// BenchmarkExploreService_PutDecision benchmarks the PutDecision method of ExploreService.
func BenchmarkExploreService_PutDecision(b *testing.B) {
	benchmarkExploreServiceMethod(b, "PutDecision")
}

// benchmarkExploreServiceMethod is a helper function that benchmarks the specified method.
func benchmarkExploreServiceMethod(b *testing.B, method string) {
	// Load environment variables
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Setup Redis client and service
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatalf("REDIS_ADDR not set in environment variables")
	}

	newStore := store.NewRedisStore(redisAddr)
	service := NewExploreService(newStore)

	// Create context
	ctx := context.Background()

	// Run the benchmark for the specified method
	for i := 0; i < b.N; i++ {
		var err error
		switch method {
		case "ListLikedYou":
			req := &generated.ListLikedYouRequest{RecipientUserId: "recipient123"}
			_, err = service.ListLikedYou(ctx, req)
		case "ListNewLikedYou":
			req := &generated.ListLikedYouRequest{RecipientUserId: "recipient123"}
			_, err = service.ListNewLikedYou(ctx, req)
		case "CountLikedYou":
			req := &generated.CountLikedYouRequest{RecipientUserId: "recipient123"}
			_, err = service.CountLikedYou(ctx, req)
		case "PutDecision":
			req := &generated.PutDecisionRequest{
				ActorUserId:     "user1",
				RecipientUserId: "recipient123",
				LikedRecipient:  true,
			}
			_, err = service.PutDecision(ctx, req)
		}
		if err != nil {
			b.Fatalf("failed to call %s: %v", method, err)
		}
	}
}

package store

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

// StoreInterface defines the methods that the store should implement
type StoreInterface interface {
	GetClient() *redis.Client
	AddLike(ctx context.Context, recipientID, actorID string) error
	IsMutualLike(ctx context.Context, actorID, recipientID string) (bool, error)
	AddMutualLike(ctx context.Context, actorID, recipientID string) error
	RecordDecision(ctx context.Context, actorID, recipientID string, likedRecipient bool) error
	GetLikesWithTimestamps(ctx context.Context, recipientID string, offset, limit int64) ([]string, []int64, error)
	CountLikes(ctx context.Context, recipientID string) (int64, error)
}

// Ensure that Store implements StoreInterface
var _ StoreInterface = (*Store)(nil)

// Store defines a struct to hold the Redis client
type Store struct {
	redisClient *redis.Client
}

// NewRedisStore initializes a new Store with a Redis client
func NewRedisStore(redisAddr string) *Store {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &Store{redisClient: rdb}
}

// GetClient returns the underlying Redis client
func (s *Store) GetClient() *redis.Client {
	return s.redisClient
}

// AddLike stores a "like" for a recipient with the current timestamp in a sorted set
func (s *Store) AddLike(ctx context.Context, recipientID string, actorID string) error {
	// Use current Unix time as the score (timestamp)
	currentTimestamp := float64(time.Now().Unix())

	// Store the like in a sorted set with the timestamp as the score
	return s.redisClient.ZAdd(ctx, "likes:"+recipientID, redis.Z{
		Score:  currentTimestamp,
		Member: actorID,
	}).Err()
}

// IsMutualLike checks if mutual like exists by checking if recipient is in actor's likes sorted set
func (s *Store) IsMutualLike(ctx context.Context, actorID, recipientID string) (bool, error) {
	// Check if "likes:<actorID>" sorted set contains recipientID
	score, err := s.redisClient.ZScore(ctx, "likes:"+actorID, recipientID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil // Member does not exist
		}
		return false, err
	}
	// If score exists, then mutual like exists
	return score != 0, nil
}

// AddMutualLike stores a mutual like flag in both directions using sets
func (s *Store) AddMutualLike(ctx context.Context, actorID, recipientID string) error {
	err1 := s.redisClient.SAdd(ctx, "mutual:"+actorID+":"+recipientID, true).Err()
	err2 := s.redisClient.SAdd(ctx, "mutual:"+recipientID+":"+actorID, true).Err()
	// Return the first error (if any)
	if err1 != nil {
		return err1
	} else if err2 != nil {
		return err2
	}
	return nil
}

// RecordDecision records the decision of the actor (like/pass), allowing overwrites
func (s *Store) RecordDecision(ctx context.Context, actorID, recipientID string, likedRecipient bool) error {
	// Optional: Use a TTL (Time-to-Live) to expire old decisions
	ttl := time.Duration(0) // Set TTL to 0 for no expiration, or adjust as needed
	return s.redisClient.Set(ctx, "decision:"+actorID+":"+recipientID, likedRecipient, ttl).Err()
}

// GetLikesWithTimestamps returns a paginated list of likers along with their like timestamps
func (s *Store) GetLikesWithTimestamps(ctx context.Context, recipientUserId string, offset int64, limit int64) ([]string, []int64, error) {
	// Define the range for fetching elements in the sorted set by score (timestamp)
	zRangeBy := &redis.ZRangeBy{
		Min:    "-inf", // Start from the lowest score (timestamp)
		Max:    "+inf", // Up to the highest score (timestamp)
		Offset: offset, // Apply offset for pagination
		Count:  limit,  // Limit the number of elements returned
	}

	// Fetch likers and their timestamps (scores) from Redis sorted set
	likersWithTimestamps, err := s.redisClient.ZRevRangeByScoreWithScores(ctx, "likes:"+recipientUserId, zRangeBy).Result()
	if err != nil {
		return nil, nil, err
	}

	var likers []string
	var timestamps []int64
	for _, likerWithTimestamp := range likersWithTimestamps {
		likers = append(likers, likerWithTimestamp.Member.(string))
		timestamps = append(timestamps, int64(likerWithTimestamp.Score)) // Redis stores scores as float64
	}

	return likers, timestamps, nil
}

// CountLikes returns the count of users who liked the recipient
func (s *Store) CountLikes(ctx context.Context, recipientID string) (int64, error) {
	// Use ZCARD to count the number of elements in the sorted set (which stores likes with timestamps)
	return s.redisClient.ZCard(ctx, "likes:"+recipientID).Result()
}

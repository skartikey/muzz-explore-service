package tests

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
)

const (
	numUsers = 100 // Number of users to generate
)

// PopulateRedisTestData populates Redis with test data. It accepts a Redis client instance as a parameter.
func PopulateRedisTestData(rdb *redis.Client) error {

	ctx := context.Background()

	// Populate Redis with users and likes
	err := populateUsers(ctx, rdb, numUsers)
	if err != nil {
		return fmt.Errorf("failed to populate users: %w", err)
	}

	// Simulate mutual likes
	err = simulateMutualLikes(ctx, rdb, numUsers)
	if err != nil {
		return fmt.Errorf("failed to simulate mutual likes: %w", err)
	}

	fmt.Println("Data insertion completed successfully.")
	return nil
}

// populateUsers generates random likes between users and adds them to Redis with timestamps
func populateUsers(ctx context.Context, rdb *redis.Client, numUsers int) error {
	for i := 1; i <= numUsers; i++ {
		for j := 1; j <= rand.Intn(50); j++ { // Each user will like between 1-50 others
			recipientID := rand.Intn(numUsers) + 1
			if i != recipientID { // Ensure users don’t like themselves
				currentTimestamp := float64(time.Now().Unix()) // Current timestamp as score

				// Add like with timestamp to the sorted set
				err := rdb.ZAdd(ctx, fmt.Sprintf("likes:%d", recipientID), redis.Z{
					Score:  currentTimestamp,
					Member: fmt.Sprintf("%d", i),
				}).Err()
				if err != nil {
					return fmt.Errorf("error adding like with timestamp: %w", err)
				}
			}
		}
	}
	return nil
}

// simulateMutualLikes simulates mutual likes between random users and adds them with timestamps
func simulateMutualLikes(ctx context.Context, rdb *redis.Client, numUsers int) error {
	for i := 1; i <= numUsers/2; i++ {
		actorID := rand.Intn(numUsers) + 1
		recipientID := rand.Intn(numUsers) + 1

		if actorID != recipientID {
			currentTimestamp := float64(time.Now().Unix()) // Current timestamp as score

			// Ensure both users like each other, with timestamps
			err1 := rdb.ZAdd(ctx, fmt.Sprintf("likes:%d", actorID), redis.Z{
				Score:  currentTimestamp,
				Member: fmt.Sprintf("%d", recipientID),
			}).Err()

			err2 := rdb.ZAdd(ctx, fmt.Sprintf("likes:%d", recipientID), redis.Z{
				Score:  currentTimestamp,
				Member: fmt.Sprintf("%d", actorID),
			}).Err()

			if err1 != nil || err2 != nil {
				return fmt.Errorf("error adding mutual like with timestamp: %w", err1)
			}

			// Add to mutual likes sets (binary flag, just to indicate mutual like)
			errMutual1 := rdb.SAdd(ctx, fmt.Sprintf("mutual:%d:%d", actorID, recipientID), true).Err()
			errMutual2 := rdb.SAdd(ctx, fmt.Sprintf("mutual:%d:%d", recipientID, actorID), true).Err()

			if errMutual1 != nil || errMutual2 != nil {
				return fmt.Errorf("error adding mutual like flag: %w", errMutual1)
			}
		}
	}
	return nil
}

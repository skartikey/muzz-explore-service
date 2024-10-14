package service

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"muzz-explore-service/pkg/protos/generated"
	"testing"
)

// MockStore is a mock implementation of the store.StoreInterface
type MockStore struct {
	mock.Mock
}

func (m *MockStore) GetClient() *redis.Client {
	args := m.Called()
	return args.Get(0).(*redis.Client)
}

func (m *MockStore) AddLike(ctx context.Context, recipientID, actorID string) error {
	args := m.Called(ctx, recipientID, actorID)
	return args.Error(0)
}

func (m *MockStore) GetLikesWithTimestamps(ctx context.Context, userID string, offset, limit int64) ([]string, []int64, error) {
	args := m.Called(ctx, userID, offset, limit)
	return args.Get(0).([]string), args.Get(1).([]int64), args.Error(2)
}

func (m *MockStore) IsMutualLike(ctx context.Context, likerID, recipientID string) (bool, error) {
	args := m.Called(ctx, likerID, recipientID)
	return args.Bool(0), args.Error(1)
}

func (m *MockStore) CountLikes(ctx context.Context, userID string) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStore) RecordDecision(ctx context.Context, actorID, recipientID string, liked bool) error {
	args := m.Called(ctx, actorID, recipientID, liked)
	return args.Error(0)
}

func (m *MockStore) AddMutualLike(ctx context.Context, actorID, recipientID string) error {
	args := m.Called(ctx, actorID, recipientID)
	return args.Error(0)
}

// TestExploreService tests all methods of the ExploreService
func TestExploreService(t *testing.T) {
	mockStore := new(MockStore)
	service := NewExploreService(mockStore)

	ctx := context.Background()

	// Test ListLikedYou
	req := &generated.ListLikedYouRequest{RecipientUserId: "recipient123"}
	mockStore.On("GetLikesWithTimestamps", ctx, "recipient123", int64(0), int64(10)).
		Return([]string{"user1", "user2"}, []int64{1633036800, 1633040400}, nil)

	resp, err := service.ListLikedYou(ctx, req)
	assert.NoError(t, err)
	assert.Len(t, resp.Likers, 2)
	assert.Equal(t, "user1", resp.Likers[0].ActorId)
	assert.Equal(t, uint64(1633036800), resp.Likers[0].UnixTimestamp)

	// Test ListNewLikedYou
	mockStore.On("GetLikesWithTimestamps", ctx, "recipient123", int64(0), int64(10)).
		Return([]string{"user1", "user2"}, []int64{1633036800, 1633040400}, nil)
	mockStore.On("IsMutualLike", ctx, "user1", "recipient123").Return(false, nil)
	mockStore.On("IsMutualLike", ctx, "user2", "recipient123").Return(true, nil)

	resp, err = service.ListNewLikedYou(ctx, req)
	assert.NoError(t, err)
	assert.Len(t, resp.Likers, 1)
	assert.Equal(t, "user1", resp.Likers[0].ActorId)

	// Test CountLikedYou
	countReq := &generated.CountLikedYouRequest{RecipientUserId: "recipient123"}
	mockStore.On("CountLikes", ctx, "recipient123").Return(int64(5), nil)

	countResp, err := service.CountLikedYou(ctx, countReq)
	assert.NoError(t, err)
	assert.Equal(t, uint64(5), countResp.Count)

	// Test PutDecision
	decisionReq := &generated.PutDecisionRequest{
		ActorUserId:     "user1",
		RecipientUserId: "recipient123",
		LikedRecipient:  true,
	}
	mockStore.On("RecordDecision", ctx, "user1", "recipient123", true).Return(nil)
	mockStore.On("IsMutualLike", ctx, "recipient123", "user1").Return(true, nil)
	mockStore.On("AddMutualLike", ctx, "user1", "recipient123").Return(nil)

	decisionResp, err := service.PutDecision(ctx, decisionReq)
	assert.NoError(t, err)
	assert.True(t, decisionResp.MutualLikes)

	// Ensure all expectations are met
	mockStore.AssertExpectations(t)
}

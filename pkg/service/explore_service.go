package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"muzz-explore-service/pkg/protos/generated"
	"muzz-explore-service/pkg/store"
)

// ExploreService defines the ExploreService using the Redis store
type ExploreService struct {
	store store.StoreInterface
	generated.UnimplementedExploreServiceServer
}

// NewExploreService initializes a new ExploreService
func NewExploreService(store store.StoreInterface) *ExploreService {
	return &ExploreService{store: store}
}

// ListLikedYou returns all users who liked the recipient with real timestamps from Redis
func (s *ExploreService) ListLikedYou(ctx context.Context, req *generated.ListLikedYouRequest) (*generated.ListLikedYouResponse, error) {
	// TODO: Use pagination parameters from the request
	//offset := req.GetOffset()
	//limit := req.GetLimit()
	offset, limit := 0, 10

	likers, timestamps, err := s.store.GetLikesWithTimestamps(ctx, req.GetRecipientUserId(), int64(offset), int64(limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get likes with timestamps: %v", err)
	}

	var likerList []*generated.ListLikedYouResponse_Liker
	for i, liker := range likers {
		likerList = append(likerList, &generated.ListLikedYouResponse_Liker{
			ActorId:       liker,
			UnixTimestamp: uint64(timestamps[i]),
		})
	}

	return &generated.ListLikedYouResponse{
		Likers: likerList,
	}, nil
}

// ListNewLikedYou returns users who liked the recipient but were not liked back, with real timestamps
func (s *ExploreService) ListNewLikedYou(ctx context.Context, req *generated.ListLikedYouRequest) (*generated.ListLikedYouResponse, error) {
	// TODO: Use pagination parameters from the request
	//offset := req.GetOffset()
	//limit := req.GetLimit()
	offset, limit := 0, 10

	likers, timestamps, err := s.store.GetLikesWithTimestamps(ctx, req.GetRecipientUserId(), int64(offset), int64(limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get new likes with timestamps: %v", err)
	}

	var newLikers []*generated.ListLikedYouResponse_Liker
	for i, liker := range likers {
		isMutual, err := s.store.IsMutualLike(ctx, liker, req.GetRecipientUserId())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to check mutual like: %v", err)
		}

		if !isMutual {
			newLikers = append(newLikers, &generated.ListLikedYouResponse_Liker{
				ActorId:       liker,
				UnixTimestamp: uint64(timestamps[i]),
			})
		}
	}

	return &generated.ListLikedYouResponse{
		Likers: newLikers,
	}, nil
}

// CountLikedYou returns the count of users who liked the recipient
func (s *ExploreService) CountLikedYou(ctx context.Context, req *generated.CountLikedYouRequest) (*generated.CountLikedYouResponse, error) {
	likeCount, err := s.store.CountLikes(ctx, req.GetRecipientUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count likes for recipient %s: %v", req.GetRecipientUserId(), err)
	}

	return &generated.CountLikedYouResponse{
		Count: uint64(likeCount),
	}, nil
}

// PutDecision records whether the actor liked or passed on the recipient, handling mutual likes
func (s *ExploreService) PutDecision(ctx context.Context, req *generated.PutDecisionRequest) (*generated.PutDecisionResponse, error) {
	// Record the actor's decision
	err := s.store.RecordDecision(ctx, req.GetActorUserId(), req.GetRecipientUserId(), req.GetLikedRecipient())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to record decision for actor %s: %v", req.GetActorUserId(), err)
	}

	// Check if it is a mutual like
	mutual := false
	if req.GetLikedRecipient() {
		isMutual, err := s.store.IsMutualLike(ctx, req.GetRecipientUserId(), req.GetActorUserId())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to check mutual like for actor %s and recipient %s: %v", req.GetActorUserId(), req.GetRecipientUserId(), err)
		}

		if isMutual {
			// Record mutual like in both directions
			err := s.store.AddMutualLike(ctx, req.GetActorUserId(), req.GetRecipientUserId())
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to add mutual like for actor %s and recipient %s: %v", req.GetActorUserId(), req.GetRecipientUserId(), err)
			}
			mutual = true
		}
	}

	return &generated.PutDecisionResponse{
		MutualLikes: mutual,
	}, nil
}

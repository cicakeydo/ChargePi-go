package grpc

import (
	"context"

	"github.com/ChargePi/ChargePi-go/internal/auth"
	"github.com/ChargePi/ChargePi-go/pkg/proto/v1/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TagAuthService struct {
	grpc.UnimplementedTagServer
	tagManager auth.Service
}

func NewTagAuthService(service auth.Service) *TagAuthService {
	return &TagAuthService{
		tagManager: service,
	}
}

func (s *TagAuthService) GetAuthorizedCards(ctx context.Context, empty *empty.Empty) (*grpc.GetAuthorizedCardsResponse, error) {
	response := &grpc.GetAuthorizedCardsResponse{
		AuthorizedCards: []*grpc.AuthorizedCard{},
	}

	// Get all tags from the database
	tags, err := s.tagManager.GetTags()
	if err != nil {
		return response, nil
	}

	// Convert the tags to the gRPC response
	for _, tag := range tags {
		var timestamp *timestamppb.Timestamp
		if tag.IdTagInfo.ExpiryDate != nil {
			timestamp = timestamppb.New(tag.IdTagInfo.ExpiryDate.Time)
		}

		card := &grpc.AuthorizedCard{
			TagId:      tag.IdTag,
			Status:     string(tag.IdTagInfo.Status),
			ExpiryDate: timestamp,
		}
		response.AuthorizedCards = append(response.AuthorizedCards, card)
	}

	return response, nil
}

func (s *TagAuthService) AddAuthorizedCards(ctx context.Context, request *grpc.AddAuthorizedCardsRequest) (*grpc.AddAuthorizedCardsResponse, error) {
	response := &grpc.AddAuthorizedCardsResponse{Status: []string{}}

	for _, tag := range request.GetAuthorizedCards() {
		err := s.tagManager.CacheTag(tag.TagId, types.NewIdTagInfo(types.AuthorizationStatus(tag.Status)))
		if err != nil {
			response.Status = append(response.Status, "Failed")
			continue
		}

		response.Status = append(response.Status, "Success")
	}

	return response, nil
}

func (s *TagAuthService) RemoveAuthorizedCard(ctx context.Context, request *grpc.RemoveCardRequest) (*grpc.RemoveCardResponse, error) {
	response := &grpc.RemoveCardResponse{
		Status: grpc.ResponseStatus_Error,
	}

	// Remove the tag from the database
	err := s.tagManager.RemoveTag(request.GetTagId())
	if err != nil {
		return response, nil
	}

	response.Status = grpc.ResponseStatus_Success
	return response, nil
}

func (s *UserService) mustEmbedUnimplementedTagServer() {
}

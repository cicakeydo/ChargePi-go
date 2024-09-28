package auth

import (
	"errors"
	"testing"

	"github.com/ChargePi/ChargePi-go/internal/auth/mocks"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type tagManagerTestSuite struct {
	suite.Suite
	mockLocalAuthListRepository *mocks.MockLocalAuthListRepository
	mockTagRepository           *mocks.MockTagRepository
	manager                     *ManagerV1
}

func (s *tagManagerTestSuite) SetupTest() {
	s.mockLocalAuthListRepository = mocks.NewMockLocalAuthListRepository(s.T())
	s.mockTagRepository = mocks.NewMockTagRepository(s.T())
	s.manager = NewManager(s.mockLocalAuthListRepository, s.mockTagRepository)
	s.manager.SetMaxTags(3)
}

func (s *tagManagerTestSuite) TearDownTest() {}

func (s *tagManagerTestSuite) TestCacheTag() {
	tests := []struct {
		name         string
		tag          *types.IdTagInfo
		tagId        string
		cacheEnabled bool
		wantErr      bool
	}{
		{
			name:         "Cache tag",
			cacheEnabled: true,
			tagId:        "1",
			tag:          okTag,
		},
		{
			name:         "Cache disabled",
			tagId:        "1",
			tag:          okTag,
			cacheEnabled: false,
			wantErr:      true,
		},
		{
			name:         "Cache full",
			tagId:        "1",
			tag:          okTag,
			cacheEnabled: true,
			wantErr:      true,
		},
		{
			name:         "Repository error",
			tagId:        "1",
			tag:          okTag,
			cacheEnabled: true,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			switch tt.name {
			case "Cache full":
				s.mockTagRepository.EXPECT().GetTags().Return(nil, nil).Once()
			case "Cache disabled":
			default:
				s.mockTagRepository.EXPECT().GetTags().Return([]*types.IdTagInfo{
					{
						Status: "Accepted",
					},
				}, nil).Once()
				s.mockTagRepository.EXPECT().AddTag(tt.tagId, tt.tag).Return(nil).Once()
			}

			s.manager.ToggleAuthCache(tt.cacheEnabled)

			err := s.manager.CacheTag(tt.tagId, tt.tag)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *tagManagerTestSuite) TestGetTag() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagManagerTestSuite) TestGetTags() {
	tests := []struct {
		name            string
		authListEnabled bool
		wantErr         bool
		wantTags        []localauth.AuthorizationData
	}{
		{
			name:            "Auth list enabled",
			authListEnabled: true,
			wantErr:         false,
			wantTags: []localauth.AuthorizationData{
				{IdTag: "1", IdTagInfo: &types.IdTagInfo{}},
			},
		},
		{
			name:            "Auth list disabled",
			authListEnabled: false,
			wantErr:         true,
		},
		{
			name:            "Local auth list error",
			authListEnabled: true,
			wantErr:         true,
			wantTags: []localauth.AuthorizationData{
				{IdTag: "1", IdTagInfo: &types.IdTagInfo{}},
			},
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			if test.name == "Local auth list error" {
				s.mockLocalAuthListRepository.EXPECT().GetLocalAuthListTags().Return(nil, errors.New("error")).Once()
			} else if test.name != "Auth list disabled" {
				s.mockLocalAuthListRepository.EXPECT().GetLocalAuthListTags().Return(test.wantTags, nil).Once()
			}

			s.manager.ToggleLocalAuthList(test.authListEnabled)

			tags, err := s.manager.GetTags()
			if test.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)
			s.NotNil(tags)
			s.ElementsMatch(test.wantTags, tags)
		})
	}
}

func (s *tagManagerTestSuite) TestRemoveTag() {
	tests := []struct {
		name            string
		tagId           string
		authListEnabled bool
		wantErr         bool
	}{
		{
			name:            "Auth list enabled",
			tagId:           "1",
			authListEnabled: true,
			wantErr:         false,
		},
		{
			name:            "Auth list disabled",
			tagId:           "1",
			authListEnabled: false,
			wantErr:         true,
		},
		{
			name:            "Local auth list error",
			tagId:           "1",
			authListEnabled: true,
			wantErr:         true,
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			if test.name == "Local auth list error" {
				s.mockLocalAuthListRepository.EXPECT().RemoveAuthListTag(test.tagId).Return(errors.New("error")).Once()
			} else if test.name != "Auth list disabled" {
				s.mockLocalAuthListRepository.EXPECT().RemoveAuthListTag(test.tagId).Return(nil).Once()
			}

			s.manager.ToggleLocalAuthList(test.authListEnabled)

			err := s.manager.RemoveTag(test.tagId)
			if test.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *tagManagerTestSuite) TestClearCache() {
	tests := []struct {
		name            string
		tagCacheEnabled bool
		wantErr         bool
	}{
		{
			name:            "Auth cache enabled",
			tagCacheEnabled: true,
			wantErr:         false,
		},
		{
			name:            "Auth cache disabled",
			tagCacheEnabled: false,
			wantErr:         true,
		},
		{
			name:            "Database error",
			tagCacheEnabled: true,
			wantErr:         true,
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			if test.name == "Database error" {
				s.mockTagRepository.EXPECT().RemoveTags().Return(errors.New("database error")).Once()
			} else if test.name == "Auth cache enabled" {
				s.mockTagRepository.EXPECT().RemoveTags().Return(nil).Once()
			}

			err := s.manager.ClearCache()
			if test.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *tagManagerTestSuite) TestUpdateLocalAuthList() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *tagManagerTestSuite) TestSetMaxTags() {
	tests := []struct {
		name            string
		tagLimit        int
		authListEnabled bool
		tagCacheEnabled bool
		want            int
	}{
		{
			name:            "Auth list enabled",
			tagLimit:        1,
			authListEnabled: true,
			tagCacheEnabled: false,
			want:            1,
		},
		{
			name:            "Auth cache enabled",
			tagLimit:        1,
			authListEnabled: false,
			tagCacheEnabled: true,
			want:            1,
		},
		{
			name:            "Auth list and cache enabled",
			tagLimit:        3,
			authListEnabled: true,
			tagCacheEnabled: true,
			want:            3,
		},
		{
			name:            "Both disabled",
			tagLimit:        1,
			authListEnabled: false,
			tagCacheEnabled: false,
			want:            1,
		},
	}

	for _, test := range tests {

		s.T().Run(test.name, func(t *testing.T) {
			s.manager.SetMaxTags(test.tagLimit)
			// Add some tags
			//

			tags, err := s.manager.authList.GetTags()
			s.NoError(err)
			s.Less(len(tags), test.want)
		})
	}
}

func (s *tagManagerTestSuite) TestToggleAuthCache() {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Auth list enabled",
			want: true,
		},
		{
			name: "Auth list disabled",
			want: false,
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			s.manager.ToggleAuthCache(test.want)
			s.Equal(test.want, s.manager.authCacheEnabled)
		})
	}
}

func (s *tagManagerTestSuite) TestToggleAuthList() {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Auth list enabled",
			want: true,
		},
		{
			name: "Auth list disabled",
			want: false,
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			s.manager.ToggleLocalAuthList(test.want)
			s.Equal(test.want, s.manager.localAuthListEnabled)
		})
	}
}

func (s *tagManagerTestSuite) TestGetAuthListVersion() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestTagManager(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	suite.Run(t, new(tagManagerTestSuite))
}

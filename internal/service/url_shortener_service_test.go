package service_test

import (
	"context"
	"testing"
	"time"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/repository/dto"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/mocks"
	"url-shortener/internal/repository/reperrors"
	"url-shortener/internal/service"
	"url-shortener/internal/service/serverrors"
)

var (
	now          = time.Now().UTC()
	past         = now.Add(-time.Hour)
	future       = now.Add(time.Hour)
	url          = "https://example.com"
	code         = "D402D"
	id     int64 = 7
)

type shortenTestCase struct {
	desc      string
	sourceURL string
	expiresAt *time.Time
	mockRepo  func(repo *mocks.MockLinkRepository)
	mockHash  func(hash *mocks.MockHashing)
	expected  string
	expectErr error
}

type resolveTestCase struct {
	desc      string
	shortCode string
	mockRepo  func(repo *mocks.MockLinkRepository)
	mockHash  func(hash *mocks.MockHashing)
	expected  string
	expectErr error
}

func TestShorten(t *testing.T) {
	t.Parallel()

	testCases := []shortenTestCase{
		{
			desc:      "new link",
			sourceURL: url,
			mockRepo: func(repo *mocks.MockLinkRepository) {
				repo.EXPECT().GetLinkBySourceUrl(gomock.Any(), gomock.Eq(url)).Return(nil, reperrors.ErrLinkNotFound)
				repo.EXPECT().CreateLink(gomock.Any(), gomock.Eq(url), gomock.Nil()).Return(int64(id), nil)
			},
			mockHash: func(hash *mocks.MockHashing) {
				hash.EXPECT().EncodeInt64(gomock.Eq([]int64{id})).Return(code, nil)
			},
			expected: code,
		},
		{
			desc:      "create expired link",
			sourceURL: url,
			expiresAt: &past,
			mockRepo: func(repo *mocks.MockLinkRepository) {
				repo.EXPECT().GetLinkBySourceUrl(gomock.Any(), gomock.Eq(url)).Return(nil, reperrors.ErrLinkNotFound)
				repo.EXPECT().UpdateExpires(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
			expectErr: serverrors.ErrExpired,
		},
		{
			desc:      "URL already exists with different expiresAt",
			sourceURL: url,
			expiresAt: &future,
			mockRepo: func(repo *mocks.MockLinkRepository) {
				repo.EXPECT().GetLinkBySourceUrl(gomock.Any(), gomock.Eq(url)).Return(&dto.LinkDTO{ID: id, SourceUrl: url}, nil)
				repo.EXPECT().UpdateExpires(gomock.Any(), gomock.Eq(int64(id)), gomock.Eq(&future)).Return(nil).Times(1)
			},
			mockHash: func(hash *mocks.MockHashing) {
				hash.EXPECT().EncodeInt64(gomock.Eq([]int64{id})).Return(code, nil)
			},
			expected: code,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			hashing, repo, svc, ctx, cancel := setupTest(t)
			defer cancel()

			if tc.mockRepo != nil {
				tc.mockRepo(repo)
			}
			if tc.mockHash != nil {
				tc.mockHash(hashing)
			}

			result, err := svc.Shorten(ctx, tc.sourceURL, tc.expiresAt)
			if tc.expectErr != nil {
				require.ErrorIs(t, err, tc.expectErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	t.Parallel()

	testCases := []resolveTestCase{
		{
			desc:      "valid short link",
			shortCode: code,
			mockHash: func(hash *mocks.MockHashing) {
				hash.EXPECT().DecodeInt64WithError(gomock.Eq(code)).Return([]int64{id}, nil)
			},
			mockRepo: func(repo *mocks.MockLinkRepository) {
				repo.EXPECT().GetLinkById(gomock.Any(), gomock.Eq(int64(id))).Return(&dto.LinkDTO{ID: id, SourceUrl: url}, nil)
				repo.EXPECT().UpdateLastRequested(gomock.Any(), gomock.Eq(int64(id))).Return(nil)
			},
			expected: url,
		},
		{
			desc:      "expired link",
			shortCode: code,
			mockHash: func(hash *mocks.MockHashing) {
				hash.EXPECT().DecodeInt64WithError(gomock.Eq(code)).Return([]int64{id}, nil)
			},
			mockRepo: func(repo *mocks.MockLinkRepository) {
				repo.EXPECT().GetLinkById(gomock.Any(), gomock.Eq(int64(id))).Return(&dto.LinkDTO{ID: id, SourceUrl: url, ExpiresAt: &past}, nil)
				repo.EXPECT().UpdateLastRequested(gomock.Any(), gomock.Eq(int64(id))).Return(nil)
			},
			expectErr: serverrors.ErrURLExpired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			hashing, repo, svc, ctx, cancel := setupTest(t)
			defer cancel()

			if tc.mockRepo != nil {
				tc.mockRepo(repo)
			}
			if tc.mockHash != nil {
				tc.mockHash(hashing)
			}

			result, err := svc.Resolve(ctx, tc.shortCode)
			if tc.expectErr != nil {
				require.ErrorIs(t, err, tc.expectErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func setupTest(t *testing.T) (*mocks.MockHashing, *mocks.MockLinkRepository, *service.UrlShortenerService, context.Context, context.CancelFunc) {
	t.Helper()

	ctrl := gomock.NewController(t)
	hashing := mocks.NewMockHashing(ctrl)
	repo := mocks.NewMockLinkRepository(ctrl)
	log := slogdiscard.NewDiscardLogger()
	svc := service.NewUrlShortenerService(log, hashing, repo)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return hashing, repo, svc, ctx, cancel
}

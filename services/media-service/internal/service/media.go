package service

import (
	"bytes"
	"context"
	"errors"
	"log/slog"

	"github.com/diyorbeknematov/minitwitter/gen/go/media"
	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/models"
	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/repository"
	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/storage"
	"github.com/diyorbeknematov/minitwitter/services/media-service/pkg/apperror"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type mediaService struct {
	repo    *repository.Repository
	storage storage.ObjectStorage
	logger  *slog.Logger
	cfg     *config.Config

	media.UnimplementedMediaServiceServer
}

func NewMediaService(
	repo *repository.Repository,
	store storage.ObjectStorage,
	logger *slog.Logger,
	cfg *config.Config,
) *mediaService {
	return &mediaService{
		repo:    repo,
		storage: store,
		logger:  logger,
		cfg:     cfg,
	}
}

func (s *mediaService) UploadMedia(ctx context.Context, req *media.UploadMediaRequest) (*media.UploadMediaResponse, error) {
	if len(req.File) == 0 {
		return nil, apperror.Wrap("service", "UploadMedia", "failed to upload media", errors.New("file is empty"))
	}

	if req.Filename == "" {
		return nil, apperror.Wrap("service", "UploadMedia", "failed to upload media", errors.New("filename is empty"))
	}

	mimeType := detectMimeType(req.File)

	// resolve prefix based on media category and mime type
	prefix, err := resolvePrefix(req.Category, mimeType)
	if err != nil {
		return nil, apperror.Wrap("service", "UploadMedia", "failed to upload media", err)
	}

	// generate object key based on prefix and mime type
	objectKey, err := generateObjectKey(prefix, mimeType)
	if err != nil {
		return nil, apperror.Wrap("service", "UploadMedia", "failed to upload media", err)
	}

	ownerId, err := uuid.Parse(req.OwnerId)
	if err != nil {
		return nil, apperror.Wrap("service", "UploadMedia", "failed to parse owner ID", err)
	}

	// upload media to storage
	_, err = s.storage.Upload(
		ctx,
		objectKey,
		bytes.NewReader(req.File),
		int64(len(req.File)),
		mimeType,
	)
	if err != nil {
		return nil, apperror.Wrap("service", "UploadMedia", "failed to upload media", err)
	}

	mediaModel := &models.Media{
		OwnerID:         ownerId,
		ObjectKey:       objectKey,
		OriginalName:    req.Filename,
		MimeType:        mimeType,
		Size:            int64(len(req.File)),
		StorageProvider: "minio",
	}

	// create media in repository
	if err := s.repo.Media.Create(ctx, mediaModel); err != nil {
		// if repo creation fails, delete the media from storage
		deleteErr := s.storage.Delete(ctx, objectKey)
		if deleteErr != nil {
			s.logger.Error(
				"failed to delete media from storage after failed repo creation",
				"object_key", objectKey,
				"error", deleteErr,
			)
		}

		return nil, apperror.Wrap("service", "UploadMedia", "failed to create media repo", err)
	}

	protoMedia, err := s.toProtoMedia(ctx, mediaModel)
	if err != nil {
		return nil, apperror.Wrap("service", "UploadMedia", "failed to convert media to proto", err)
	}

	return &media.UploadMediaResponse{
		Media: protoMedia,
	}, nil
}

func (s *mediaService) GetMedia(ctx context.Context, req *media.GetMediaRequest) (*media.Media, error) {
	// get media from repository
	m, err := s.repo.Media.GetByID(ctx, req.MediaId)
	if err != nil {
		return nil, apperror.Wrap("service", "GetMedia", "failed to get media", err)
	}

	protoMedia, err := s.toProtoMedia(ctx, m)
	if err != nil {
		return nil, apperror.Wrap("service", "GetMedia", "failed to convert media to proto", err)
	}

	return protoMedia, nil
}

func (s *mediaService) GetMedias(ctx context.Context, req *media.GetMediasRequest) (*media.GetMediasResponse, error) {
	medias, err := s.repo.Media.GetMany(ctx, req.MediaIds)
	if err != nil {
		return nil, apperror.Wrap("service", "GetMedias", "failed to get medias", err)
	}

	mediaList := make([]*media.Media, 0, len(medias))
	for _, m := range medias {

		protoMedia, err := s.toProtoMedia(ctx, &m)
		if err != nil {
			return nil, apperror.Wrap("service", "GetMedias", "failed to convert media to proto", err)
		}
		mediaList = append(mediaList, protoMedia)
	}

	return &media.GetMediasResponse{
		Medias: mediaList,
	}, nil
}

func (s *mediaService) DeleteMedia(ctx context.Context, req *media.DeleteMediaRequest) (*emptypb.Empty, error) {
	media, err := s.repo.Media.GetByID(ctx, req.MediaId)
	if err != nil {
		return nil, apperror.Wrap("service", "DeleteMedia", "failed to get media", err)
	}

	// delete media from storage first
	if err := s.storage.Delete(ctx, media.ObjectKey); err != nil {
		return nil, apperror.Wrap("service", "DeleteMedia", "failed to delete media from storage", err)
	}

	// delete media from repository
	if err := s.repo.Media.Delete(ctx, req.MediaId); err != nil {
		return nil, apperror.Wrap("service", "DeleteMedia", "failed to delete media from repository", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *mediaService) toProtoMedia(
	ctx context.Context,
	m *models.Media,
) (*media.Media, error) {

	url, err := s.storage.PresignedURL(ctx, m.ObjectKey, s.cfg.MinIO.PresignedExpiry)
	if err != nil {
		return nil, err
	}

	return &media.Media{
		Id:              m.ID.String(),
		OwnerId:         m.OwnerID.String(),
		ObjectKey:       m.ObjectKey,
		Url:             url,
		OriginalName:    m.OriginalName,
		MimeType:        m.MimeType,
		Size:            m.Size,
		StorageProvider: m.StorageProvider,
		CreatedAt:       timestamppb.New(m.CreatedAt),
	}, nil
}

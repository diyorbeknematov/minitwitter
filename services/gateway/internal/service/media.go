package service

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/gen/go/media"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/grpcclient"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/mapper"
	"github.com/diyorbeknematov/minitwitter/services/gateway/pkg/apperror"
	"github.com/google/uuid"
)

type mediaService struct {
	client *grpcclient.Client
}

func NewMediaService(c *grpcclient.Client) *mediaService {
	return &mediaService{
		client: c,
	}
}

func (s *mediaService) UploadMedia(
	ctx context.Context,
	userID uuid.UUID,
	file []byte,
	filename string,
	req dto.UploadMediaReq,
) (dto.Media, error) {

	resp, err := s.client.Media.UploadMedia(
		ctx,
		mapper.ToUploadMediaReq(userID, file, filename, req),
	)
	if err != nil {
		return dto.Media{}, apperror.Wrap(
			"service",
			"UploadMedia",
			"failed to upload media",
			err,
		)
	}

	media, err := mapper.ToMedia(resp.Media)
	if err != nil {
		return dto.Media{}, apperror.Wrap(
			"service",
			"UploadMedia",
			"failed to map media",
			err,
		)
	}

	return media, nil
}

func (s *mediaService) GetMedia(
	ctx context.Context,
	mediaID uuid.UUID,
) (dto.Media, error) {

	resp, err := s.client.Media.GetMedia(
		ctx,
		&media.GetMediaRequest{
			MediaId: mediaID.String(),
		},
	)
	if err != nil {
		return dto.Media{}, apperror.Wrap(
			"service",
			"GetMedia",
			"failed to get media",
			err,
		)
	}

	return mapper.ToMedia(resp)
}

func (s *mediaService) GetMedias(
	ctx context.Context,
	mediaIDs []uuid.UUID,
) (dto.GetMediasResp, error) {

	resp, err := s.client.Media.GetMedias(
		ctx,
		mapper.ToGetMediasReq(mediaIDs),
	)
	if err != nil {
		return dto.GetMediasResp{}, apperror.Wrap(
			"service",
			"GetMedias",
			"failed to get medias",
			err,
		)
	}

	return mapper.ToGetMediasResp(resp)
}

func (s *mediaService) DeleteMedia(
	ctx context.Context,
	mediaID uuid.UUID,
) error {

	_, err := s.client.Media.DeleteMedia(
		ctx,
		&media.DeleteMediaRequest{
			MediaId: mediaID.String(),
		},
	)
	if err != nil {
		return apperror.Wrap(
			"service",
			"DeleteMedia",
			"failed to delete media",
			err,
		)
	}

	return nil
}

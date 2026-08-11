package mapper

import (
	"github.com/diyorbeknematov/minitwitter/gen/go/media"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/google/uuid"
)

func ToMedia(pb *media.Media) (dto.Media, error) {
	id, err := uuid.Parse(pb.Id)
	if err != nil {
		return dto.Media{}, nil
	}

	return dto.Media{
		ID:              id,
		ObjectKey:       pb.ObjectKey,
		Url:             pb.Url,
		OriginalName:    pb.OriginalName,
		MimeType:        pb.MimeType,
		Size:            pb.Size,
		StorageProvider: pb.StorageProvider,
		CreatedAt:       pb.CreatedAt.AsTime(),
	}, nil
}

func ToUploadMediaReq(
	userID uuid.UUID,
	file []byte,
	filename string,
	req dto.UploadMediaReq,
) (*media.UploadMediaRequest, error) {

	var category media.MediaCategory

	switch req.Category {
	case dto.MediaCategoryTweet:
		category = media.MediaCategory_TWEET

	case dto.MediaCategoryAvatar:
		category = media.MediaCategory_AVATAR

	default:
		category = media.MediaCategory_MEDIA_CATEGORY_UNSPECIFIED
	}

	return &media.UploadMediaRequest{
		OwnerId:  userID.String(),
		File:     file,
		Filename: filename,
		Category: category,
	}, nil
}

func ToGetMediasResponse(
	pb *media.GetMediasResponse,
) (dto.GetMediasResp, error) {

	medias := make([]dto.Media, len(pb.Medias))

	for i, m := range pb.Medias {
		mediaDTO, err := ToMedia(m)
		if err != nil {
			return dto.GetMediasResp{}, err
		}

		medias[i] = mediaDTO
	}

	return dto.GetMediasResp{
		Medias: medias,
	}, nil
}
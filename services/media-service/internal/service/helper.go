package service

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/diyorbeknematov/minitwitter/gen/go/media"
	"github.com/diyorbeknematov/minitwitter/services/media-service/pkg/apperror"
	"github.com/google/uuid"
)

const (
	imagesFolder = "images"
	videosFolder = "videos"

	tweetsFolder  = "tweets"
	avatarsFolder = "avatars"
)

func detectMimeType(file []byte) string {
	return http.DetectContentType(file[:512])
}

func resolvePrefix(category media.MediaCategory, mimeType string) (string, error) {
	var base string

	switch category {
	case media.MediaCategory_TWEET:
		base = tweetsFolder

	case media.MediaCategory_AVATAR:
		base = avatarsFolder

	default:
		return "", apperror.Wrap("service", "resolvePrefix", "failed to resolve prefix", errors.New("invalid media category"))
	}

	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return path.Join(base, imagesFolder), nil

	case strings.HasPrefix(mimeType, "video/"):
		if category == media.MediaCategory_AVATAR {
			return "", apperror.Wrap("service", "resolvePrefix", "failed to resolve prefix", errors.New("avatar must be an image"))
		}

		return path.Join(base, videosFolder), nil

	default:
		return "", apperror.Wrap("service", "resolvePrefix", "failed to resolve prefix", errors.New("unsupported media type"))
	}
}

func generateObjectKey(prefix, mimeType string) (string, error) {
	id := uuid.NewString()

	extentions, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return "", apperror.Wrap("service", "generateObjectKey", "failed to generate object key", err)
	}

	if len(extentions) == 0 {
		return "", apperror.Wrap("service", "generateObjectKey", "failed to generate object key", fmt.Errorf("unsupported mime type: %s", mimeType))
	}

	return path.Join(prefix, id+extentions[0]), nil
}

package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cloudinary *cloudinary.Cloudinary
	folder     string
}

func NewCloudinaryService(cloudName, apiKey, apiSecret, folder string) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	return &CloudinaryService{
		cloudinary: cld,
		folder:     folder,
	}, nil
}

func (s *CloudinaryService) UploadImage(ctx context.Context, file multipart.File, filename string) (string, error) {
	uploadParams := uploader.UploadParams{
		PublicID:     filename,
		Folder:       s.folder,
		ResourceType: "image",
		Timestamp:    time.Now().Unix(),
	}

	result, err := s.cloudinary.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	return result.SecureURL, nil
}

func (s *CloudinaryService) UploadAvatar(ctx context.Context, file multipart.File, userID string) (string, error) {
	filename := fmt.Sprintf("avatar_%s_%d", userID, time.Now().Unix())
	return s.UploadImage(ctx, file, filename)
}

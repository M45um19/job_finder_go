package cloudinary

import (
	"context"
	"io"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Service struct {
	cld *cloudinary.Cloudinary
}

func New(cloudName, apiKey, apiSecret string) (*Service, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}
	return &Service{cld: cld}, nil
}

func (s *Service) UploadFile(ctx context.Context, file io.Reader, folder string) (string, error) {
	resp, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}

// DeleteFile deletes an asset from Cloudinary using its public ID
func (s *Service) DeleteFile(ctx context.Context, publicID string) error {
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

// ExtractPublicID parses the Cloudinary URL and returns the asset's public ID.
// E.g., "https://res.cloudinary.com/demo/image/upload/v12345/user_photos/profile.jpg" -> "user_photos/profile"
func ExtractPublicID(url string) string {
	idx := strings.Index(url, "/upload/")
	if idx == -1 {
		return ""
	}

	part := url[idx+len("/upload/"):]
	segments := strings.Split(part, "/")
	if len(segments) == 0 {
		return ""
	}

	// Skip the version segment (e.g. "v1234567") if it exists
	startIndex := 0
	if strings.HasPrefix(segments[0], "v") {
		isVersion := true
		for _, char := range segments[0][1:] {
			if char < '0' || char > '9' {
				isVersion = false
				break
			}
		}
		if isVersion {
			startIndex = 1
		}
	}

	publicIDWithExt := strings.Join(segments[startIndex:], "/")
	lastDot := strings.LastIndex(publicIDWithExt, ".")
	if lastDot != -1 {
		return publicIDWithExt[:lastDot]
	}
	return publicIDWithExt
}

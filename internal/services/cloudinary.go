package services

import (
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryInstance struct {
	Client *cloudinary.Cloudinary
}

func InitCloudinary() *CloudinaryInstance {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}
	return &CloudinaryInstance{Client: cld}
}

// UploadMedia handles Images, Videos, and Files (PDF/Docs)
func (c *CloudinaryInstance) UploadMedia(ctx context.Context, file interface{}, folder string) (string, error) {
	resp, err := c.Client.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:          folder,
		// "auto" tells Cloudinary to detect if it's an image, video, or raw file
		ResourceType:    "auto", 
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}
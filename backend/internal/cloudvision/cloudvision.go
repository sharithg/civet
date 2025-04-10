package cloudvision

import (
	"bytes"
	"context"
	"fmt"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
)

type CloudVision struct {
	client   *vision.ImageAnnotatorClient
	cacheDir string
}

func NewCloudVision(ctx context.Context, cacheDir string) (*CloudVision, error) {
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return nil, err
	}

	return &CloudVision{
		client:   client,
		cacheDir: cacheDir,
	}, nil
}

func (cv *CloudVision) DetectText(ctx context.Context, content []byte) ([]*visionpb.EntityAnnotation, error) {

	img, err := vision.NewImageFromReader(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %w", err)
	}

	annotations, err := cv.client.DetectTexts(ctx, img, nil, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to detect texts: %w", err)
	}

	return annotations, nil
}

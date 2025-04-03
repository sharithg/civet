package cloudvision

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

func (cv *CloudVision) DetectText(ctx context.Context, content []byte, imageHash string) ([]*visionpb.EntityAnnotation, error) {
	cachePath := filepath.Join(cv.cacheDir, imageHash+".bin")

	if data, err := os.ReadFile(cachePath); err == nil {
		fmt.Printf("Loaded cached response for hash %s\n", imageHash)
		return cv.deserializeResponse(data)
	}

	fmt.Printf("Calling API for hash %s\n", imageHash)
	img, err := vision.NewImageFromReader(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %w", err)
	}

	annotations, err := cv.client.DetectTexts(ctx, img, nil, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to detect texts: %w", err)
	}

	if err := cv.serializeResponse(cachePath, annotations); err != nil {
		log.Printf("Failed to cache result: %v", err)
	}

	return annotations, nil
}

func (cv *CloudVision) serializeResponse(path string, annotations []*visionpb.EntityAnnotation) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(annotations)
	if err != nil {
		return err
	}

	return os.WriteFile(path, buf.Bytes(), 0644)
}

func (cv *CloudVision) deserializeResponse(data []byte) ([]*visionpb.EntityAnnotation, error) {
	var annotations []*visionpb.EntityAnnotation
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(&annotations)
	if err != nil {
		return nil, err
	}

	return annotations, nil
}

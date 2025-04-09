package receipt

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/sharithg/civet/internal/cloudvision"
	"github.com/sharithg/civet/internal/genai"
	"github.com/sharithg/civet/internal/storage"
)

const prompt = "Convert the given text of a receipt into a structured output format"

type Extract struct {
	ImageBytes   []byte
	FileName     string
	ImageHash    string
	FileExt      string
	visionClient *cloudvision.CloudVision
	openaiClient genai.OpenAi
	storage      storage.Storage
}

func NewExtract(ctx context.Context, storage storage.Storage, openai genai.OpenAi, imageBytes []byte, fname string) (*Extract, error) {
	hash := sha256.Sum256(imageBytes)
	imageHash := hex.EncodeToString(hash[:])
	ext := strings.TrimPrefix(filepath.Ext(fname), ".")

	visionClient, err := cloudvision.NewCloudVision(ctx, "cache/cloud_vision")
	if err != nil {
		return nil, err
	}

	// openaiClient := genai.NewOpenAiClient()

	return &Extract{
		ImageBytes:   imageBytes,
		FileName:     fname,
		ImageHash:    imageHash,
		FileExt:      ext,
		visionClient: visionClient,
		openaiClient: openai,
		storage:      storage,
	}, nil
}

func (e *Extract) Upload(ctx context.Context) (bucket, key string, err error) {
	objectName := fmt.Sprintf("%s.%s", e.ImageHash, e.FileExt)
	bucket = "receipts"
	_, err = e.storage.UploadImageBytes(ctx, bucket, objectName, e.ImageBytes, "image/"+e.FileExt)
	return bucket, objectName, err
}

func (e *Extract) ExtractText(ctx context.Context) (string, error) {
	annotations, err := e.visionClient.DetectText(ctx, e.ImageBytes, e.ImageHash)
	if err != nil {
		return "", err
	}
	lines := GroupTextByLines(annotations, 10)
	return strings.Join(lines, "\n"), nil
}

func (e *Extract) StructuredOutput(ctx context.Context, input string) (Receipt, error) {
	var Schema = GenerateSchema[Receipt]()

	return genai.JsonChat[Receipt](ctx, &e.openaiClient, prompt, input, "receipt_info", Schema)
}

func (e *Extract) Run(ctx context.Context) (ParsedReceipt, string, string, string, error) {
	bucket, key, err := e.Upload(ctx)
	if err != nil {
		return ParsedReceipt{}, "", "", "", err
	}

	text, err := e.ExtractText(ctx)
	if err != nil {
		return ParsedReceipt{}, "", "", "", err
	}

	out, err := e.StructuredOutput(ctx, text)
	if err != nil {
		return ParsedReceipt{}, "", "", "", err
	}

	model, err := e.ToModel(out)
	if err != nil {
		return ParsedReceipt{}, "", "", "", err
	}

	return model, text, bucket, key, nil
}

func (e *Extract) ToModel(output Receipt) (ParsedReceipt, error) {
	dateStr := output.Opened
	var parsed time.Time

	if dateStr != "" {
		t, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04:05", dateStr)
		}
		if err != nil {
			log.Printf("[WARN] Unable to parse date string: %s, err: %v", dateStr, err)
		} else {
			parsed = t
		}
	}

	return ParsedReceipt{
		Receipt: Receipt{
			Restaurant:  output.Restaurant,
			Address:     output.Address,
			Opened:      output.Opened,
			OrderNumber: output.OrderNumber,
			OrderType:   output.OrderType,
			Table:       output.Table,
			Server:      output.Server,
			Items:       output.Items,
			Subtotal:    output.Subtotal,
			SalesTax:    output.SalesTax,
			Total:       output.Total,
			Payment:     output.Payment,
			Copy:        output.Copy,
			OtherFees:   output.OtherFees,
		},
		Opened: parsed,
	}, nil
}

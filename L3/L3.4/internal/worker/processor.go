package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	stdimage "image"
	"image/color"
	"image/draw"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/infrastructure"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/models"
	imagesvc "github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/service/image"
)

// Processor отвечает за фоновые задачи обработки изображений
type Processor struct {
	storage infrastructure.Storage
	queue   infrastructure.MessageQueue
	service imagesvc.Service
}

// NewProcessor конструктор
func NewProcessor(storage infrastructure.Storage, queue infrastructure.MessageQueue, service imagesvc.Service) *Processor {
	return &Processor{
		storage: storage,
		queue:   queue,
		service: service,
	}
}

// Run запускает обработку
func (p *Processor) Run(ctx context.Context, topic, groupID string) error {
	if topic == "" || groupID == "" {
		return fmt.Errorf("topic and groupID are required")
	}

	log.Printf("[worker] subscribing to topic %s (group %s)", topic, groupID)
	if err := p.queue.Subscribe(topic, groupID, p.handleMessage); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	<-ctx.Done()
	return nil
}

func (p *Processor) handleMessage(msg []byte) error {
	var task models.ImageTaskMessage
	if err := json.Unmarshal(msg, &task); err != nil {
		return fmt.Errorf("decode task: %w", err)
	}

	ctx := context.Background()
	log.Printf("[worker] start processing image=%s type=%s", task.ImageID, task.ProcessingType)

	if _, err := p.service.MarkProcessing(ctx, task.ImageID); err != nil {
		return fmt.Errorf("mark processing: %w", err)
	}

	originalBucket, originalObject := ensureBucketAndKey(task.OriginalBucket, task.OriginalPath)
	processedBucket, processedObject := ensureBucketAndKey(task.ProcessedBucket, task.ProcessedPath)
	if originalObject == "" || processedObject == "" {
		return fmt.Errorf("object keys must be provided")
	}

	reader, err := p.storage.Download(originalBucket, originalObject)
	if err != nil {
		p.service.MarkFailed(ctx, task.ImageID, err)
		return fmt.Errorf("download original: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		p.service.MarkFailed(ctx, task.ImageID, err)
		return fmt.Errorf("read original: %w", err)
	}

	processedBytes, contentType, err := processImage(data, task.ProcessingType, processedObject)
	if err != nil {
		p.service.MarkFailed(ctx, task.ImageID, err)
		return fmt.Errorf("process image: %w", err)
	}

	if err := p.storage.Upload(processedBucket, processedObject, bytes.NewReader(processedBytes), int64(len(processedBytes)), contentType); err != nil {
		p.service.MarkFailed(ctx, task.ImageID, err)
		return fmt.Errorf("upload processed: %w", err)
	}

	if _, err := p.service.MarkReady(ctx, task.ImageID, processedObject); err != nil {
		return fmt.Errorf("mark ready: %w", err)
	}

	log.Printf("[worker] image=%s processed successfully", task.ImageID)
	return nil
}

func processImage(payload []byte, processingType models.ProcessingType, targetPath string) ([]byte, string, error) {
	img, _, err := stdimage.Decode(bytes.NewReader(payload))
	if err != nil {
		return nil, "", fmt.Errorf("decode image: %w", err)
	}

	var result stdimage.Image

	switch processingType {
	case models.ProcessingResize:
		result = imaging.Resize(img, 1280, 0, imaging.Lanczos)
	case models.ProcessingThumbnail:
		result = imaging.Thumbnail(img, 320, 320, imaging.Lanczos)
	case models.ProcessingWatermark:
		resized := imaging.Resize(img, 1280, 0, imaging.Lanczos)
		result = applyWatermark(resized, "ImageProcessor")
	default:
		return nil, "", fmt.Errorf("unsupported processing type: %s", processingType)
	}

	targetFormat, contentType := pickFormat(targetPath)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, result, targetFormat); err != nil {
		return nil, "", fmt.Errorf("encode image: %w", err)
	}

	return buf.Bytes(), contentType, nil
}

func applyWatermark(img stdimage.Image, text string) stdimage.Image {
	bounds := img.Bounds()
	canvas := imaging.Clone(img)

	overlay := stdimage.NewUniform(color.RGBA{0, 0, 0, 80})
	draw.Draw(canvas, stdimage.Rect(bounds.Max.X-220, bounds.Max.Y-80, bounds.Max.X, bounds.Max.Y), overlay, stdimage.Point{}, draw.Over)

	d := &font.Drawer{
		Dst:  canvas,
		Src:  stdimage.NewUniform(color.RGBA{255, 255, 255, 160}),
		Face: basicfont.Face7x13,
	}

	padding := 20
	x := bounds.Max.X - padding - len(text)*7
	y := bounds.Max.Y - padding
	if x < padding {
		x = padding
	}
	d.Dot = fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}

	d.DrawString(text)

	return canvas
}

func pickFormat(path string) (imaging.Format, string) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
		return imaging.PNG, "image/png"
	case ".gif":
		return imaging.GIF, "image/gif"
	case ".jpg", ".jpeg":
		return imaging.JPEG, "image/jpeg"
	default:
		return imaging.JPEG, "image/jpeg"
	}
}

func ensureBucketAndKey(bucket, path string) (string, string) {
	if bucket != "" {
		return bucket, strings.TrimPrefix(path, "/")
	}

	path = strings.TrimPrefix(path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return bucket, path
}

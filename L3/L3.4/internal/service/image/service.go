package image

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/infrastructure"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/models"
)

var (
	errNilStorage        = errors.New("storage dependency is required")
	errNilQueue          = errors.New("message queue dependency is required")
	errInvalidConfig     = errors.New("invalid image service config")
	errInvalidInputFile  = errors.New("file payload is required")
	errUnsupportedFormat = errors.New("unsupported processing type")
)

var supportedProcessing = map[models.ProcessingType]struct{}{
	models.ProcessingResize:    {},
	models.ProcessingThumbnail: {},
	models.ProcessingWatermark: {},
}

type service struct {
	storage infrastructure.Storage
	queue   infrastructure.MessageQueue
	cfg     Config
}

// NewService конструктор ImageService
func NewService(storage infrastructure.Storage, queue infrastructure.MessageQueue, cfg Config) (Service, error) {
	if storage == nil {
		return nil, errNilStorage
	}
	if queue == nil {
		return nil, errNilQueue
	}
	if cfg.Buckets.Original == "" || cfg.Buckets.Processed == "" || cfg.TopicImageTasks == "" {
		return nil, errInvalidConfig
	}
	if cfg.DefaultProcessing == "" {
		cfg.DefaultProcessing = models.ProcessingResize
	}
	if _, ok := supportedProcessing[cfg.DefaultProcessing]; !ok {
		return nil, fmt.Errorf("%w: %s", errUnsupportedFormat, cfg.DefaultProcessing)
	}

	return &service{
		storage: storage,
		queue:   queue,
		cfg:     cfg,
	}, nil
}

func (s *service) Upload(ctx context.Context, file io.Reader, size int64, contentType, filename string, processingType models.ProcessingType) (*models.ImageInfo, error) {
	ctx = normalizeContext(ctx)
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if file == nil || size <= 0 {
		return nil, errInvalidInputFile
	}

	pt, err := s.resolveProcessingType(processingType)
	if err != nil {
		return nil, err
	}

	imageID := uuid.NewString()
	ext := detectExtension(filename, contentType)
	originalObject := buildOriginalObjectKey(imageID, ext)
	processedObject := buildProcessedObjectKey(imageID, pt, ext)

	if err := s.storage.Upload(s.cfg.Buckets.Original, originalObject, file, size, contentType); err != nil {
		return nil, fmt.Errorf("upload original image: %w", err)
	}

	metadata := &models.ImageInfo{
		ImageID:     imageID,
		OriginalURL: originalObject,
		Status:      models.StatusPending,
	}

	if err := s.storage.SaveMetadata(imageID, metadata); err != nil {
		_ = s.storage.Delete(s.cfg.Buckets.Original, originalObject)
		return nil, fmt.Errorf("save metadata: %w", err)
	}

	task := models.ImageTaskMessage{
		ImageID:         imageID,
		OriginalBucket:  s.cfg.Buckets.Original,
		OriginalPath:    originalObject,
		ProcessedBucket: s.cfg.Buckets.Processed,
		ProcessedPath:   processedObject,
		ProcessingType:  pt,
	}

	payload, err := json.Marshal(task)
	if err != nil {
		metadata.Status = models.StatusFailed
		metadata.ErrorMessage = "failed to serialize task payload"
		_ = s.storage.UpdateMetadata(imageID, metadata)
		return nil, fmt.Errorf("marshal task: %w", err)
	}

	if err := s.queue.Publish(s.cfg.TopicImageTasks, payload); err != nil {
		metadata.Status = models.StatusFailed
		metadata.ErrorMessage = fmt.Sprintf("queue publish: %v", err)
		_ = s.storage.UpdateMetadata(imageID, metadata)
		return nil, fmt.Errorf("publish task: %w", err)
	}

	return s.decorate(metadata), nil
}

func (s *service) Get(ctx context.Context, imageID string) (*models.ImageInfo, error) {
	ctx = normalizeContext(ctx)
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	info, err := s.storage.GetMetadata(imageID)
	if err != nil {
		return nil, fmt.Errorf("get metadata: %w", err)
	}

	return s.decorate(info), nil
}

func (s *service) List(ctx context.Context, limit int) ([]*models.ImageInfo, error) {
	ctx = normalizeContext(ctx)
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	items, err := s.storage.ListMetadata(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("list metadata: %w", err)
	}

	return s.decorateMany(items), nil
}

func (s *service) Delete(ctx context.Context, imageID string) error {
	ctx = normalizeContext(ctx)
	if err := ctx.Err(); err != nil {
		return err
	}

	info, err := s.storage.GetMetadata(imageID)
	if err != nil {
		return fmt.Errorf("get metadata: %w", err)
	}

	if key := strings.TrimSpace(info.OriginalURL); key != "" {
		if err := s.storage.Delete(s.cfg.Buckets.Original, key); err != nil {
			return fmt.Errorf("delete original object: %w", err)
		}
	}

	if key := strings.TrimSpace(info.ProcessedURL); key != "" {
		if err := s.storage.Delete(s.cfg.Buckets.Processed, key); err != nil {
			return fmt.Errorf("delete processed object: %w", err)
		}
	}

	if err := s.storage.DeleteMetadata(imageID); err != nil {
		return fmt.Errorf("delete metadata: %w", err)
	}

	return nil
}

func (s *service) MarkProcessing(ctx context.Context, imageID string) (*models.ImageInfo, error) {
	ctx = normalizeContext(ctx)
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	info, err := s.storage.GetMetadata(imageID)
	if err != nil {
		return nil, fmt.Errorf("get metadata: %w", err)
	}

	info.Status = models.StatusProcessing
	info.ErrorMessage = ""

	if err := s.storage.UpdateMetadata(imageID, info); err != nil {
		return nil, fmt.Errorf("update metadata: %w", err)
	}

	return s.decorate(info), nil
}

func (s *service) MarkReady(ctx context.Context, imageID string, processedObjectKey string) (*models.ImageInfo, error) {
	ctx = normalizeContext(ctx)
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if processedObjectKey == "" {
		return nil, errors.New("processed object key is required")
	}

	info, err := s.storage.GetMetadata(imageID)
	if err != nil {
		return nil, fmt.Errorf("get metadata: %w", err)
	}

	info.Status = models.StatusReady
	info.ProcessedURL = processedObjectKey
	info.ErrorMessage = ""

	if err := s.storage.UpdateMetadata(imageID, info); err != nil {
		return nil, fmt.Errorf("update metadata: %w", err)
	}

	return s.decorate(info), nil
}

func (s *service) MarkFailed(ctx context.Context, imageID string, processingErr error) (*models.ImageInfo, error) {
	ctx = normalizeContext(ctx)
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	info, err := s.storage.GetMetadata(imageID)
	if err != nil {
		return nil, fmt.Errorf("get metadata: %w", err)
	}

	info.Status = models.StatusFailed
	if processingErr != nil {
		info.ErrorMessage = processingErr.Error()
	}

	if err := s.storage.UpdateMetadata(imageID, info); err != nil {
		return nil, fmt.Errorf("update metadata: %w", err)
	}

	return s.decorate(info), nil
}

func (s *service) resolveProcessingType(pt models.ProcessingType) (models.ProcessingType, error) {
	if pt == "" {
		return s.cfg.DefaultProcessing, nil
	}
	if _, ok := supportedProcessing[pt]; !ok {
		return "", fmt.Errorf("%w: %s", errUnsupportedFormat, pt)
	}
	return pt, nil
}

func (s *service) decorate(info *models.ImageInfo) *models.ImageInfo {
	if info == nil {
		return nil
	}

	clone := *info
	if info.OriginalURL != "" {
		clone.OriginalURL = s.storage.GetURL(s.cfg.Buckets.Original, info.OriginalURL)
	}
	if info.ProcessedURL != "" {
		clone.ProcessedURL = s.storage.GetURL(s.cfg.Buckets.Processed, info.ProcessedURL)
	}

	return &clone
}

func (s *service) decorateMany(items []*models.ImageInfo) []*models.ImageInfo {
	result := make([]*models.ImageInfo, 0, len(items))
	for _, item := range items {
		result = append(result, s.decorate(item))
	}
	return result
}

func detectExtension(filename, contentType string) string {
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(filename)))
	if ext != "" {
		return ensureExtFormat(ext)
	}

	if contentType != "" {
		if exts, err := mime.ExtensionsByType(contentType); err == nil && len(exts) > 0 {
			return ensureExtFormat(exts[0])
		}
	}

	return ".bin"
}

func ensureExtFormat(ext string) string {
	if ext == "" {
		return ".bin"
	}
	if !strings.HasPrefix(ext, ".") {
		return strings.ToLower("." + ext)
	}
	return strings.ToLower(ext)
}

func buildOriginalObjectKey(imageID, ext string) string {
	return fmt.Sprintf("%s/original%s", imageID, ensureExtFormat(ext))
}

func buildProcessedObjectKey(imageID string, pt models.ProcessingType, ext string) string {
	return fmt.Sprintf("%s/%s%s", imageID, string(pt), ensureExtFormat(ext))
}

func normalizeContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

package service

import (
	"context"
	"errors"
	"fmt"

	"uuid"

	"velo-server/internal/model"
	"velo-server/internal/repository"
	"velo-server/internal/storage"
)

type MusicService struct {
	tracks *repository.TrackRepo
	r2     *storage.R2Storage
}

var (
	ErrTrackNotFound = errors.New("track not found")
)

func NewMusicService(
	tracks *repository.TrackRepo,
	r2 *storage.R2Storage,
) *MusicService {
	return &MusicService{
		tracks: tracks,
		r2:     r2,
	}
}

// RegisterTrack saves track metadata in PostgreSQL.
func (s *MusicService) RegisterTrack(
	ctx context.Context,
	track *model.Track,
) (*model.Track, error) {

	if track.Title == "" {
		return nil, fmt.Errorf("track title is required")
	}

	if track.Artist == "" {
		return nil, fmt.Errorf("track artist is required")
	}

	if track.ObjectKey == "" {
		return nil, fmt.Errorf("object key is required")
	}

	if track.MIMEType == "" {
		return nil, fmt.Errorf("mime type is required")
	}

	if track.Duration < 0 || track.FileSize < 0 {
		return nil, fmt.Errorf("invalid duration or file size")
	}

	return s.tracks.Create(ctx, track)
}

// GetTrack retrieves a track by its ID.
func (s *MusicService) GetTrack(
	ctx context.Context,
	id uuid.UUID,
) (*model.Track, error) {

	track, err := s.tracks.GetByID(ctx, id)

	if track == nil {
		return nil, ErrTrackNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get track: %w", err)
	}

	return track, nil
}

// GetTracks retrieves all tracks.
func (s *MusicService) GetTracks(
	ctx context.Context,
	limit int,
	offset int,
) ([]model.Track, error) {

	return s.tracks.GetAll(ctx, limit, offset)
}

// GetStreamURL generates a temporary URL for a track.
func (s *MusicService) GetStreamURL(
	ctx context.Context,
	id uuid.UUID,
) (string, error) {

	track, err := s.GetTrack(ctx, id)
	if err != nil {
		return "", err
	}

	url, err := s.r2.GeneratePresignedURL(
		ctx,
		track.ObjectKey,
	)

	if err != nil {
		return "", fmt.Errorf("generate stream URL: %w", err)
	}

	return url, nil
}

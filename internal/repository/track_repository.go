package repository

import (
	"context"
	"database/sql"
	"fmt"
	"velo-server/internal/model"

	"uuid"
)

type TrackRepo struct {
	db *sql.DB
}

func NewTrackRepository(db *sql.DB) *TrackRepo {
	return &TrackRepo{db: db}
}

// create entry for new track
func (r *TrackRepo) Create(
	ctx context.Context,
	track *model.Track,
) (*model.Track, error) {

	query := `
        INSERT INTO tracks (
            title,
            artist,
            album,
            object_key,
            duration,
            mime_type,
            file_size,
            cover_art_url
        )
        VALUES ($1, $2, NULLIF($3, ''), $4, $5, $6, $7, NULLIF($8, ''))
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRowContext(
		ctx,
		query,
		track.Title,
		track.Artist,
		track.Album,
		track.ObjectKey,
		track.Duration,
		track.MIMEType,
		track.FileSize,
		track.CoverArtURL,
	).Scan(
		&track.ID,
		&track.CreatedAt,
		&track.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create track: %w", err)
	}

	return track, nil
}

func (r *TrackRepo) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.Track, error) {

	query := `
        SELECT
            id,
            title,
            artist,
            COALESCE(album, ''),
            object_key,
            duration,
            mime_type,
            file_size,
            COALESCE(cover_art_url, ''),
            created_at,
            updated_at
        FROM tracks
        WHERE id = $1
    `

	var track model.Track

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&track.ID,
		&track.Title,
		&track.Artist,
		&track.Album,
		&track.ObjectKey,
		&track.Duration,
		&track.MIMEType,
		&track.FileSize,
		&track.CoverArtURL,
		&track.CreatedAt,
		&track.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("get track: %w", err)
	}

	return &track, nil
}

func (r *TrackRepo) GetAll(
	ctx context.Context,
	limit int,
	offset int,
) ([]model.Track, error) {

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	if offset < 0 {
		offset = 0
	}

	query := `
        SELECT
            id,
            title,
            artist,
            COALESCE(album, ''),
            object_key,
            duration,
            mime_type,
            file_size,
            COALESCE(cover_art_url, ''),
            created_at,
            updated_at
        FROM tracks
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := r.db.QueryContext(
		ctx,
		query,
		limit,
		offset,
	)

	if err != nil {
		return nil, fmt.Errorf("list tracks: %w", err)
	}

	defer rows.Close()

	tracks := make([]model.Track, 0)

	for rows.Next() {
		var track model.Track

		err := rows.Scan(
			&track.ID,
			&track.Title,
			&track.Artist,
			&track.Album,
			&track.ObjectKey,
			&track.Duration,
			&track.MIMEType,
			&track.FileSize,
			&track.CoverArtURL,
			&track.CreatedAt,
			&track.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan track: %w", err)
		}

		tracks = append(tracks, track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tracks: %w", err)
	}

	return tracks, nil
}

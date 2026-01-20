package stories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

var (
	ErrFailedToGetRecord    = errors.New("failed to get record")
	ErrFailedToListRecords  = errors.New("failed to list records")
	ErrFailedToInsertRecord = errors.New("failed to insert record")
	ErrFailedToUpdateRecord = errors.New("failed to update record")
	ErrFailedToRemoveRecord = errors.New("failed to remove record")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrStoryNotFound        = errors.New("story not found")
	ErrInvalidSlugPrefix    = errors.New("slug must start with YYYYMMDD of publish date")
)

// validateSlugDatePrefix validates that the slug starts with YYYYMMDD format
// matching the provided publish date.
func validateSlugDatePrefix(slug string, publishDate time.Time) error {
	expectedPrefix := publishDate.Format("20060102")

	// Check if slug starts with the expected date prefix
	if !strings.HasPrefix(slug, expectedPrefix) {
		return fmt.Errorf("%w: expected prefix %s", ErrInvalidSlugPrefix, expectedPrefix)
	}

	return nil
}

type Repository interface {
	GetProfileIDBySlug(ctx context.Context, slug string) (string, error)
	GetProfileByID(ctx context.Context, localeCode string, id string) (*profiles.Profile, error)
	GetStoryIDBySlug(ctx context.Context, slug string) (string, error)
	GetStoryByID(
		ctx context.Context,
		localeCode string,
		id string,
		authorProfileID *string,
	) (*StoryWithChildren, error)
	ListStoriesOfPublication(
		ctx context.Context,
		localeCode string,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*StoryWithChildren], error)
	// CRUD methods
	InsertStory(
		ctx context.Context,
		id string,
		authorProfileID string,
		slug string,
		kind string,
		status string,
		isFeatured bool,
		storyPictureURI *string,
		properties map[string]any,
		publishedAt *time.Time,
	) (*Story, error)
	InsertStoryTx(
		ctx context.Context,
		storyID string,
		localeCode string,
		title string,
		summary string,
		content string,
	) error
	InsertStoryPublication(
		ctx context.Context,
		id string,
		storyID string,
		profileID string,
		kind string,
		properties map[string]any,
	) error
	UpdateStory(
		ctx context.Context,
		id string,
		slug string,
		status string,
		isFeatured bool,
		storyPictureURI *string,
		publishedAt *time.Time,
	) error
	UpdateStoryTx(
		ctx context.Context,
		storyID string,
		localeCode string,
		title string,
		summary string,
		content string,
	) error
	UpsertStoryTx(
		ctx context.Context,
		storyID string,
		localeCode string,
		title string,
		summary string,
		content string,
	) error
	RemoveStory(ctx context.Context, id string) error
	GetStoryForEdit(ctx context.Context, localeCode string, id string) (*StoryForEdit, error)
	GetStoryOwnershipForUser(
		ctx context.Context,
		userID string,
		storyID string,
	) (*StoryOwnership, error)
}

type Service struct {
	logger      *logfx.Logger
	repo        Repository
	idGenerator RecordIDGenerator
}

func NewService(logger *logfx.Logger, repo Repository) *Service {
	return &Service{logger: logger, repo: repo, idGenerator: DefaultIDGenerator}
}

func (s *Service) GetByID(
	ctx context.Context,
	localeCode string,
	id string,
) (*StoryWithChildren, error) {
	record, err := s.repo.GetStoryByID(ctx, localeCode, id, nil)
	if err != nil {
		return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, id, err)
	}

	return record, nil
}

func (s *Service) GetBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
) (*StoryWithChildren, error) {
	storyID, err := s.repo.GetStoryIDBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	record, err := s.repo.GetStoryByID(ctx, localeCode, storyID, nil)
	if err != nil {
		return nil, fmt.Errorf("%w(story_id: %s): %w", ErrFailedToGetRecord, storyID, err)
	}

	return record, nil
}

func (s *Service) List(
	ctx context.Context,
	localeCode string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*StoryWithChildren], error) {
	records, err := s.repo.ListStoriesOfPublication(ctx, localeCode, cursor)
	if err != nil {
		return cursors.Cursored[[]*StoryWithChildren]{}, fmt.Errorf(
			"%w: %w",
			ErrFailedToListRecords,
			err,
		)
	}

	return records, nil
}

func (s *Service) ListByPublicationProfileSlug(
	ctx context.Context,
	localeCode string,
	publicationProfileSlug string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*StoryWithChildren], error) {
	publicationProfileID, err := s.repo.GetProfileIDBySlug(ctx, publicationProfileSlug)
	if err != nil {
		return cursors.Cursored[[]*StoryWithChildren]{}, fmt.Errorf(
			"%w(slug: %s): %w",
			ErrFailedToGetRecord,
			publicationProfileSlug,
			err,
		)
	}

	cursor.Filters["publication_profile_id"] = publicationProfileID

	records, err := s.repo.ListStoriesOfPublication(
		ctx,
		localeCode,
		cursor,
	)
	if err != nil {
		return cursors.Cursored[[]*StoryWithChildren]{}, fmt.Errorf(
			"%w: %w",
			ErrFailedToListRecords,
			err,
		)
	}

	return records, nil
}

// CanUserEditStory checks if a user has permission to edit a story.
func (s *Service) CanUserEditStory(
	ctx context.Context,
	userID string,
	storyID string,
) (bool, error) {
	ownership, err := s.repo.GetStoryOwnershipForUser(ctx, userID, storyID)
	if err != nil {
		return false, fmt.Errorf(
			"%w(userID: %s, storyID: %s): %w",
			ErrFailedToGetRecord,
			userID,
			storyID,
			err,
		)
	}

	if ownership == nil {
		return false, nil
	}

	return ownership.CanEdit, nil
}

// GetForEdit retrieves a story for editing (raw content, no compilation).
func (s *Service) GetForEdit(
	ctx context.Context,
	localeCode string,
	storyID string,
) (*StoryForEdit, error) {
	story, err := s.repo.GetStoryForEdit(ctx, localeCode, storyID)
	if err != nil {
		return nil, fmt.Errorf("%w(storyID: %s): %w", ErrFailedToGetRecord, storyID, err)
	}

	return story, nil
}

// Create creates a new story with its translation and publication.
func (s *Service) Create(
	ctx context.Context,
	userID string,
	authorProfileSlug string,
	publicationProfileSlug string,
	localeCode string,
	slug string,
	kind string,
	status string,
	title string,
	summary string,
	content string,
	storyPictureURI *string,
	isFeatured bool,
	publishedAt *time.Time,
) (*Story, error) {
	// Validate slug date prefix for published stories
	if status == "published" && publishedAt != nil {
		err := validateSlugDatePrefix(slug, *publishedAt)
		if err != nil {
			return nil, err
		}
	}

	// Get author profile ID
	authorProfileID, err := s.repo.GetProfileIDBySlug(ctx, authorProfileSlug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, authorProfileSlug, err)
	}

	// Get publication profile ID
	publicationProfileID, err := s.repo.GetProfileIDBySlug(ctx, publicationProfileSlug)
	if err != nil {
		return nil, fmt.Errorf(
			"%w(slug: %s): %w",
			ErrFailedToGetRecord,
			publicationProfileSlug,
			err,
		)
	}

	// Generate new story ID
	storyID := s.idGenerator()

	// Create the main story record
	story, err := s.repo.InsertStory(
		ctx,
		string(storyID),
		authorProfileID,
		slug,
		kind,
		status,
		isFeatured,
		storyPictureURI,
		nil, // No additional properties for now
		publishedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToInsertRecord, err)
	}

	// Create the localized story data
	err = s.repo.InsertStoryTx(
		ctx,
		string(storyID),
		localeCode,
		title,
		summary,
		content,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: translation: %w", ErrFailedToInsertRecord, err)
	}

	// Create story publication (link story to profile)
	publicationID := s.idGenerator()

	err = s.repo.InsertStoryPublication(
		ctx,
		string(publicationID),
		string(storyID),
		publicationProfileID,
		"original", // publication kind
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: publication: %w", ErrFailedToInsertRecord, err)
	}

	return story, nil
}

// Update updates an existing story.
func (s *Service) Update(
	ctx context.Context,
	userID string,
	storyID string,
	slug string,
	status string,
	isFeatured bool,
	storyPictureURI *string,
	publishedAt *time.Time,
) (*StoryForEdit, error) {
	// Check authorization
	canEdit, err := s.CanUserEditStory(ctx, userID, storyID)
	if err != nil {
		return nil, err
	}

	if !canEdit {
		return nil, fmt.Errorf(
			"%w: user %s cannot edit story %s",
			ErrUnauthorized,
			userID,
			storyID,
		)
	}

	// Validate slug date prefix for published stories
	if status == "published" && publishedAt != nil {
		err := validateSlugDatePrefix(slug, *publishedAt)
		if err != nil {
			return nil, err
		}
	}

	// Update the story
	err = s.repo.UpdateStory(ctx, storyID, slug, status, isFeatured, storyPictureURI, publishedAt)
	if err != nil {
		return nil, fmt.Errorf("%w(storyID: %s): %w", ErrFailedToUpdateRecord, storyID, err)
	}

	// Return updated story
	story, err := s.repo.GetStoryForEdit(ctx, "en", storyID)
	if err != nil {
		return nil, fmt.Errorf("%w(storyID: %s): %w", ErrFailedToGetRecord, storyID, err)
	}

	return story, nil
}

// UpdateTranslation updates story translation fields.
func (s *Service) UpdateTranslation(
	ctx context.Context,
	userID string,
	storyID string,
	localeCode string,
	title string,
	summary string,
	content string,
) error {
	// Check authorization
	canEdit, err := s.CanUserEditStory(ctx, userID, storyID)
	if err != nil {
		return err
	}

	if !canEdit {
		return fmt.Errorf(
			"%w: user %s cannot edit story %s",
			ErrUnauthorized,
			userID,
			storyID,
		)
	}

	// Update the translation (use upsert to handle new locales)
	err = s.repo.UpsertStoryTx(ctx, storyID, localeCode, title, summary, content)
	if err != nil {
		return fmt.Errorf(
			"%w(storyID: %s, locale: %s): %w",
			ErrFailedToUpdateRecord,
			storyID,
			localeCode,
			err,
		)
	}

	return nil
}

// Delete soft-deletes a story.
func (s *Service) Delete(
	ctx context.Context,
	userID string,
	storyID string,
) error {
	// Check authorization
	canEdit, err := s.CanUserEditStory(ctx, userID, storyID)
	if err != nil {
		return err
	}

	if !canEdit {
		return fmt.Errorf(
			"%w: user %s cannot delete story %s",
			ErrUnauthorized,
			userID,
			storyID,
		)
	}

	// Delete the story
	err = s.repo.RemoveStory(ctx, storyID)
	if err != nil {
		return fmt.Errorf("%w(storyID: %s): %w", ErrFailedToRemoveRecord, storyID, err)
	}

	return nil
}

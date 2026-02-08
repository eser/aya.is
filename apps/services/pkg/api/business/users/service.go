package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

var (
	ErrFailedToGetRecord    = errors.New("failed to get record")
	ErrFailedToListRecords  = errors.New("failed to list records")
	ErrFailedToCreateRecord = errors.New("failed to create record")
	ErrFailedToUpdateRecord = errors.New("failed to update record")
)

type Repository interface {
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByGitHubRemoteID(ctx context.Context, githubRemoteID string) (*User, error)
	ListUsers(
		ctx context.Context,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*User], error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	SetUserIndividualProfileID(ctx context.Context, userID string, profileID string) error

	CreateSession(ctx context.Context, session *Session) error
	GetSessionByID(ctx context.Context, id string) (*Session, error)
	UpdateSessionLoggedInAt(ctx context.Context, id string, loggedInAt time.Time) error
	UpdateSessionStatus(ctx context.Context, id string, status string) error
	CopySessionPreferences(ctx context.Context, oldSessionID, newSessionID string) error
	ListSessionsByUserID(ctx context.Context, userID string) ([]*Session, error)
	UpdateSessionActivity(ctx context.Context, id string, userAgent *string) error
	TerminateSession(ctx context.Context, sessionID, userID string) error
}

type Service struct {
	logger       *logfx.Logger
	repo         Repository
	auditService *events.AuditService
	idGenerator  RecordIDGenerator
}

func NewService(
	logger *logfx.Logger,
	repo Repository,
	auditService *events.AuditService,
) *Service {
	return &Service{
		logger:       logger,
		repo:         repo,
		auditService: auditService,
		idGenerator:  DefaultIDGenerator,
	}
}

func (s *Service) GenerateID() RecordID {
	return s.idGenerator()
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	record, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, id, err)
	}

	return record, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	record, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%w(email: %s): %w", ErrFailedToGetRecord, email, err)
	}

	return record, nil
}

func (s *Service) List(
	ctx context.Context,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*User], error) {
	records, err := s.repo.ListUsers(ctx, cursor)
	if err != nil {
		return cursors.Cursored[[]*User]{}, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return records, nil
}

func (s *Service) Create(ctx context.Context, user *User) error {
	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.UserCreated,
		EntityType: "user",
		EntityID:   user.ID,
		ActorKind:  events.ActorSystem,
	})

	return nil
}

func (s *Service) GetSessionByID(ctx context.Context, id string) (*Session, error) {
	session, err := s.repo.GetSessionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, id, err)
	}

	return session, nil
}

func (s *Service) CreateSession(ctx context.Context, session *Session) error {
	err := s.repo.CreateSession(ctx, session)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	return nil
}

func (s *Service) UpdateSessionLoggedInAt(
	ctx context.Context,
	id string,
	loggedInAt time.Time,
) error {
	err := s.repo.UpdateSessionLoggedInAt(ctx, id, loggedInAt)
	if err != nil {
		return fmt.Errorf("%w(id: %s): %w", ErrFailedToUpdateRecord, id, err)
	}

	return nil
}

func (s *Service) InvalidateSession(ctx context.Context, id string) error {
	err := s.repo.UpdateSessionStatus(ctx, id, SessionStatusLoggedOut.String())
	if err != nil {
		return fmt.Errorf("%w(id: %s): %w", ErrFailedToUpdateRecord, id, err)
	}

	return nil
}

func (s *Service) ListSessionsByUserID(ctx context.Context, userID string) ([]*Session, error) {
	sessions, err := s.repo.ListSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w(user_id: %s): %w", ErrFailedToListRecords, userID, err)
	}

	return sessions, nil
}

func (s *Service) UpdateSessionActivity(ctx context.Context, id string, userAgent *string) error {
	err := s.repo.UpdateSessionActivity(ctx, id, userAgent)
	if err != nil {
		return fmt.Errorf("%w(id: %s): %w", ErrFailedToUpdateRecord, id, err)
	}

	return nil
}

func (s *Service) TerminateSession(ctx context.Context, sessionID, userID string) error {
	err := s.repo.TerminateSession(ctx, sessionID, userID)
	if err != nil {
		return fmt.Errorf(
			"%w(session_id: %s, user_id: %s): %w",
			ErrFailedToUpdateRecord,
			sessionID,
			userID,
			err,
		)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.SessionTerminated,
		EntityType: "session",
		EntityID:   sessionID,
		ActorID:    &userID,
		ActorKind:  events.ActorUser,
	})

	return nil
}

func (s *Service) UpsertGitHubUser( //nolint:funlen
	ctx context.Context,
	githubRemoteID string,
	email string,
	name string,
	githubHandle string,
) (*User, error) {
	// First, try to find user by GitHub Remote ID
	existingUser, err := s.repo.GetUserByGitHubRemoteID(ctx, githubRemoteID)
	if err != nil {
		return nil, fmt.Errorf(
			"%w(github_remote_id: %s): %w",
			ErrFailedToGetRecord,
			githubRemoteID,
			err,
		)
	}

	if existingUser != nil {
		// User exists, update their GitHub information
		existingUser.Name = name
		existingUser.Email = &email
		existingUser.GithubHandle = &githubHandle
		existingUser.GithubRemoteID = &githubRemoteID

		now := time.Now()
		existingUser.UpdatedAt = &now

		err = s.repo.UpdateUser(ctx, existingUser)
		if err != nil {
			return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToUpdateRecord, existingUser.ID, err)
		}

		s.auditService.Record(ctx, events.AuditParams{
			EventType:  events.UserUpdated,
			EntityType: "user",
			EntityID:   existingUser.ID,
			ActorKind:  events.ActorSystem,
			Payload:    map[string]any{"github_remote_id": githubRemoteID},
		})

		return existingUser, nil
	}

	// User not found by GitHub ID, try to find by email
	if email != "" {
		existingUser, err = s.repo.GetUserByEmail(ctx, email)
		if err != nil {
			return nil, fmt.Errorf("%w(email: %s): %w", ErrFailedToGetRecord, email, err)
		}

		if existingUser != nil {
			// User exists with same email, update their GitHub information
			existingUser.Name = name
			existingUser.GithubHandle = &githubHandle
			existingUser.GithubRemoteID = &githubRemoteID

			now := time.Now()
			existingUser.UpdatedAt = &now

			err = s.repo.UpdateUser(ctx, existingUser)
			if err != nil {
				return nil, fmt.Errorf(
					"%w(id: %s): %w",
					ErrFailedToUpdateRecord,
					existingUser.ID,
					err,
				)
			}

			s.auditService.Record(ctx, events.AuditParams{
				EventType:  events.UserUpdated,
				EntityType: "user",
				EntityID:   existingUser.ID,
				ActorKind:  events.ActorSystem,
				Payload:    map[string]any{"github_remote_id": githubRemoteID},
			})

			return existingUser, nil
		}
	}

	// User doesn't exist, create new one
	newUser := &User{
		ID:                  string(s.idGenerator()),
		Kind:                "regular",
		Name:                name,
		Email:               &email,
		GithubHandle:        &githubHandle,
		GithubRemoteID:      &githubRemoteID,
		BskyHandle:          nil,
		XHandle:             nil,
		IndividualProfileID: nil,
		CreatedAt:           time.Now(),
		UpdatedAt:           nil,
		DeletedAt:           nil,
	}

	err = s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
	}

	s.auditService.Record(ctx, events.AuditParams{
		EventType:  events.UserCreated,
		EntityType: "user",
		EntityID:   newUser.ID,
		ActorKind:  events.ActorSystem,
		Payload:    map[string]any{"github_remote_id": githubRemoteID},
	})

	return newUser, nil
}

func (s *Service) SetIndividualProfileID(
	ctx context.Context,
	userID string,
	profileID string,
) error {
	err := s.repo.SetUserIndividualProfileID(ctx, userID, profileID)
	if err != nil {
		return fmt.Errorf(
			"%w(user_id: %s, profile_id: %s): %w",
			ErrFailedToUpdateRecord,
			userID,
			profileID,
			err,
		)
	}

	return nil
}

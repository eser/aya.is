package privacy

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrNoAccount            = errors.New("account not found")
	ErrAccountDeleted       = errors.New("account has been deleted")
	ErrExportUnavailable    = errors.New("personal data export is not available on this server")
	ErrExportInProgress     = errors.New("a personal data export is already in progress")
	ErrExportNotFound       = errors.New("export not found")
	ErrExportNotReady       = errors.New("export is not ready")
	ErrExportExpired        = errors.New("export has expired")
	ErrConfirmationMismatch = errors.New("confirmation does not match the username")
	ErrExportTooLarge       = errors.New("export archive exceeds the size limit")

	ErrFailedToGetAccount = errors.New("failed to get account")
	ErrFailedToHide       = errors.New("failed to hide profile")
	ErrFailedToUnhide     = errors.New("failed to restore profile")
	ErrFailedToExport     = errors.New("failed to build personal data export")
	ErrFailedToDelete     = errors.New("failed to delete account")
	ErrFailedToEnqueue    = errors.New("failed to schedule export")
)

// CooldownError reports that a new export may not be requested yet.
type CooldownError struct {
	NextRequestAt time.Time
}

func (e *CooldownError) Error() string {
	return "a new export can be requested after " + e.NextRequestAt.Format(time.RFC3339)
}

// OwnedProfilesError reports that deleting the account would orphan profiles it solely owns.
type OwnedProfilesError struct {
	Profiles []OwnedProfile
}

func (e *OwnedProfilesError) Error() string {
	slugs := make([]string, 0, len(e.Profiles))
	for _, profile := range e.Profiles {
		slugs = append(slugs, profile.Slug)
	}

	return "transfer ownership before deleting the account: " + strings.Join(slugs, ", ")
}

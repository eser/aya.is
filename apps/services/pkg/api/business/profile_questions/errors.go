package profile_questions

import "errors"

// Sentinel errors.
var (
	ErrQANotEnabled            = errors.New("Q&A is not enabled for this profile")
	ErrQuestionNotFound        = errors.New("question not found")
	ErrContentTooShort         = errors.New("question content is too short")
	ErrContentTooLong          = errors.New("question content is too long")
	ErrInsufficientPermission  = errors.New("insufficient permission for this action")
	ErrQuestionAlreadyAnswered = errors.New("question has already been answered")
	ErrQuestionNotAnswered     = errors.New("question has not been answered yet")

	ErrFailedToGetRecord    = errors.New("failed to get record")
	ErrFailedToInsertRecord = errors.New("failed to insert record")
	ErrFailedToUpdateRecord = errors.New("failed to update record")
	ErrFailedToListRecords  = errors.New("failed to list records")
)

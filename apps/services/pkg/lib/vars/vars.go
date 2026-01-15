package vars

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/sqlc-dev/pqtype"
)

var (
	ErrMustBePointer  = errors.New("must be a pointer")
	ErrCannotAssign   = errors.New("cannot assign")
	ErrCannotSetValue = errors.New("cannot set a value")
)

func ToStringPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}

	return nil
}

func ToSQLNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}

	return sql.NullString{
		String: "",
		Valid:  false,
	}
}

func ToTimePtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}

	return nil
}

func ToSQLNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Time: *t, Valid: true}
	}

	return sql.NullTime{
		Time:  time.Time{},
		Valid: false,
	}
}

func ToRawMessage(m pqtype.NullRawMessage) []byte {
	if m.Valid {
		return m.RawMessage
	}

	return nil
}

func ToObject(m pqtype.NullRawMessage) any {
	if m.Valid {
		var obj any

		err := json.Unmarshal(m.RawMessage, &obj)
		if err != nil {
			return nil
		}

		return obj
	}

	return nil
}

func MapValueToNullString(m map[string]string, key string) sql.NullString {
	if v, ok := m[key]; ok {
		return sql.NullString{String: v, Valid: true}
	}

	return sql.NullString{
		String: "",
		Valid:  false,
	}
}

func ToSQLNullRawMessage(m map[string]any) pqtype.NullRawMessage {
	if m != nil {
		bytes, err := json.Marshal(m)
		if err != nil {
			return pqtype.NullRawMessage{Valid: false}
		}

		return pqtype.NullRawMessage{RawMessage: bytes, Valid: true}
	}

	return pqtype.NullRawMessage{Valid: false}
}

func SetValue(dest any, src any) error {
	dv := reflect.ValueOf(dest)
	if dv.Kind() != reflect.Pointer {
		return fmt.Errorf("%w: %s", ErrMustBePointer, dv.Kind())
	}

	ev := dv.Elem() //nolint:varnamelen

	if !ev.CanSet() {
		return fmt.Errorf("%w: %s", ErrCannotSetValue, ev.Type())
	}

	if src == nil {
		ev.Set(reflect.Zero(ev.Type()))

		return nil
	}

	sv := reflect.ValueOf(src) //nolint:varnamelen

	if !sv.Type().AssignableTo(ev.Type()) {
		return fmt.Errorf("%w: %s to %s", ErrCannotAssign, sv.Type(), ev.Type())
	}

	ev.Set(sv)

	return nil
}

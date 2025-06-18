package helper

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToPGText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  s != "",
	}
}

func NullString(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  s != "",
	}
}

func ToPGInt32(s string, defaultVal int32) int32 {
	if s == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return int32(i)
}

func ToPGUUID(id uuid.UUID) pgtype.UUID {
	var pgID pgtype.UUID
	_ = pgID.Scan(id.String())
	return pgID
}

func ToPGUUIDFromString(s string) (pgtype.UUID, error) {
	var pgID pgtype.UUID
	err := pgID.Scan(s)
	return pgID, err
}

func NullUUIDFromString(s string) pgtype.UUID {
	var pgID pgtype.UUID
	if s == "" {
		return pgtype.UUID{Valid: false}
	}
	_ = pgID.Scan(s)
	return pgID
}

func ToPGBool(b bool) pgtype.Bool {
	return pgtype.Bool{
		Bool:  b,
		Valid: true,
	}
}

func ToPGTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

// parsing data type
func ParseInt32(s string, defaultVal int) int32 {
	if s == "" {
		return int32(defaultVal)
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return int32(defaultVal)
	}
	return int32(i)
}

func ParseInt64(s string, defaultVal int) int64 {
	if s == "" {
		return int64(defaultVal)
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return int64(defaultVal)
	}
	return int64(i)
}

func ParseFloat64(s string, defaultVal float64) float64 {
	if s == "" {
		return defaultVal
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultVal
	}
	return f
}

func IsEmptyString(s string) bool {
	return s == ""
}

func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

package pagination

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Cursorable interface {
	GetCursorFields() (time.Time, uint64)
}

type Cursor struct {
	Time      time.Time
	Direction string
	ID        uint64
}

func EncodeCursor(t time.Time, id uint64, direction string) string {
	raw := fmt.Sprintf("%s_%d:%s", t.UTC().Format(time.RFC3339Nano), id, direction)
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func DecodeCursor(s string) (*Cursor, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cursor: %w", err)
	}

	raw := string(data)
	lastColon := strings.LastIndexByte(raw, ':')
	if lastColon == -1 {
		return nil, errors.New("invalid cursor: missing direction")
	}

	timeAndID := raw[:lastColon]
	direction := raw[lastColon+1:]

	lastUnderscore := strings.LastIndexByte(timeAndID, '_')
	if lastUnderscore == -1 {
		return nil, errors.New("invalid cursor: missing id separator")
	}

	timeStr := timeAndID[:lastUnderscore]
	idStr := timeAndID[lastUnderscore+1:]

	t, err1 := time.Parse(time.RFC3339Nano, timeStr)
	id, err2 := strconv.ParseUint(idStr, 10, 64)
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("invalid time or id: %w / %w", err1, err2)
	}

	return &Cursor{
		Time:      t,
		ID:        id,
		Direction: direction,
	}, nil
}

func ApplyOrderBy(db *gorm.DB, direction string) *gorm.DB {
	if direction == "asc" {
		return db.Order("created_at ASC").Order("id ASC")
	}
	return db.Order("created_at DESC").Order("id DESC")
}

func ApplyCursorFilter[T Cursorable](db *gorm.DB, c *Cursor) *gorm.DB {
	if c == nil {
		return db
	}
	if c.Direction == "asc" {
		return db.Where("(created_at > ?) OR (created_at = ? AND id > ?)", c.Time, c.Time, c.ID)
	}
	return db.Where("(created_at < ?) OR (created_at = ? AND id < ?)", c.Time, c.Time, c.ID)
}

func MakeCursorPair[T Cursorable](items []T, limit int, direction string) (next, prev string, hasMore bool) {
	hasMore = len(items) > limit
	if hasMore {
		items = items[:limit]
	}

	if len(items) == 0 {
		return "", "", false
	}

	first := items[0]
	last := items[len(items)-1]

	firstTime, firstID := first.GetCursorFields()
	lastTime, lastID := last.GetCursorFields()

	next = EncodeCursor(lastTime, lastID, direction)
	prev = EncodeCursor(firstTime, firstID, invertDirection(direction))

	return next, prev, hasMore
}

func invertDirection(d string) string {
	if d == "asc" {
		return "desc"
	}
	return "asc"
}

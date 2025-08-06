package pagination

import (
	"testing"
	"time"
)

type TestItem struct {
	CreatedAt time.Time
	ID        uint64
}

func (t TestItem) GetCursorFields() (time.Time, uint64) {
	return t.CreatedAt, t.ID
}

func TestCursorEncodingAndDecoding(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Microsecond)
	id := uint64(42)
	direction := "desc"

	encoded := EncodeCursor(now, id, direction)
	cursor, err := DecodeCursor(encoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !cursor.Time.Equal(now) {
		t.Errorf("Expected time %v, got %v", now, cursor.Time)
	}
	if cursor.ID != id {
		t.Errorf("Expected id %d, got %d", id, cursor.ID)
	}
	if cursor.Direction != direction {
		t.Errorf("Expected direction %s, got %s", direction, cursor.Direction)
	}
}

func TestMakeCursorPair(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Microsecond)
	items := []TestItem{
		{CreatedAt: now.Add(-2 * time.Minute), ID: 1},
		{CreatedAt: now.Add(-1 * time.Minute), ID: 2},
		{CreatedAt: now, ID: 3},
	}

	next, prev, hasMore := MakeCursorPair(items, 2, "desc")

	if next == "" || prev == "" {
		t.Fatal("Cursors should not be nil")
	}
	if !hasMore {
		t.Error("Expected hasMore to be true")
	}
}

package pagination

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func CursorFilter[T Cursorable](cursor string, db *gorm.DB, direction string) *gorm.DB {
	ApplyOrderBy(db, direction)

	if cursor == "" {
		return db
	}

	parts := strings.Split(cursor, "_")
	if len(parts) != 2 {
		return db
	}

	t, err1 := time.Parse(time.RFC3339Nano, parts[0])
	if err1 != nil {
		log.Println("Invalid time in cursor:", err1)
		return db
	}

	id, err2 := strconv.ParseUint(parts[1], 10, 64)
	if err2 != nil {
		log.Println("Invalid ID in cursor:", err2)
		return db
	}

	ApplyCursorFilter[T](db, &Cursor{
		Time:      t,
		ID:        id,
		Direction: direction,
	})

	return db
}

func NextCursor[T Cursorable](items []T, limit int) (nextCursor string, hasMore bool, result []T) {
	if len(items) > limit {
		items = items[:limit]
		hasMore = true
		last := items[len(items)-1]

		createdAt, id := last.GetCursorFields()
		cursorVal := fmt.Sprintf("%s_%d", createdAt.UTC().Format(time.RFC3339Nano), id)
		nextCursor = cursorVal
	}
	return nextCursor, hasMore, items
}

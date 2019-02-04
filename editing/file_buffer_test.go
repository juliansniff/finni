package editing

import (
	"testing"
)

func TestMoveCursorUp(t *testing.T) {
	tests := map[int]*FileBuffer{
		3: &FileBuffer{
			File:   []byte("test\ntest"),
			Cursor: 3,
		},
		4: &FileBuffer{
			File:   []byte("t\ntest"),
			Cursor: 4,
		},
		0: &FileBuffer{
			File:   []byte("test\ntest"),
			Cursor: 5,
		},
		1: &FileBuffer{
			File:   []byte("testing\ntest"),
			Cursor: 9,
		},
		13: &FileBuffer{
			File:   []byte("t\ntest\ntesting"),
			Cursor: 13,
		},
		5: &FileBuffer{
			File:   []byte("t\ntest\ntesting"),
			Cursor: 11,
		},
	}
	for expected, fb := range tests {
		_ = fb.MoveCursorUp()
		if expected != fb.Cursor {
			t.Errorf("Expected: %d, got: %d", expected, fb.Cursor)
		}
	}

}

func TestMoveCursorRight(t *testing.T) {
	tests := map[int]*FileBuffer{
		3: &FileBuffer{
			File:   []byte("test"),
			Cursor: 3,
		},
		1: &FileBuffer{
			File:   []byte("test"),
			Cursor: 0,
		},
	}
	for expected, fb := range tests {
		_ = fb.MoveCursorRight()
		if expected != fb.Cursor {
			t.Errorf("Expected: %d, got: %d", expected, fb.Cursor)
		}
	}
}

func TestMoveCursorDown(t *testing.T) {
	tests := map[int]*FileBuffer{
		4: &FileBuffer{
			File:   []byte("test\nt"),
			Cursor: 4,
		},
		7: &FileBuffer{
			File:   []byte("test\ntesting"),
			Cursor: 3,
		},
		8: &FileBuffer{
			File:   []byte("test\ntest\nt"),
			Cursor: 8,
		},
		13: &FileBuffer{
			File:   []byte("test\ntest\ntest"),
			Cursor: 8,
		},
	}
	for expected, fb := range tests {
		_ = fb.MoveCursorDown()
		if expected != fb.Cursor {
			t.Errorf("Expected: %d, got: %d", expected, fb.Cursor)
		}
	}
}

func TestMoveCursorLeft(t *testing.T) {
	tests := map[int]*FileBuffer{
		0: &FileBuffer{
			File:   []byte("test"),
			Cursor: 0,
		},
		1: &FileBuffer{
			File:   []byte("test"),
			Cursor: 2,
		},
	}
	for expected, fb := range tests {
		_ = fb.MoveCursorLeft()
		if expected != fb.Cursor {
			t.Errorf("Expected: %d, got: %d", expected, fb.Cursor)
		}
	}
}

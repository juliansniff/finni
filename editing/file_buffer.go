package editing

import "fmt"

type FileBuffer struct {
	File   []byte
	Cursor int
}

func NewFileBuffer() *FileBuffer {
	return &FileBuffer{}
}

func (fb *FileBuffer) MoveCursorUp() error {
	newLine := byte('\n')
	previousLineLength, before := 0, 0
	for fb.File[fb.Cursor-before] != newLine && fb.Cursor-before > 0 {
		before++
	}
	for (previousLineLength == 0 || fb.File[fb.Cursor-before-previousLineLength] != newLine) && fb.Cursor-before-previousLineLength > 0 {
		previousLineLength++
	}
	fmt.Printf("%d, %d", before, previousLineLength)
	if previousLineLength < before {
		return fmt.Errorf("Cannot move cursor up")
	}
	position := fb.Cursor - before - (previousLineLength - before + 1)
	if position < 0 {
		return fmt.Errorf("Cannot move cursor up")
	}
	fb.Cursor = position
	return nil
}

func (fb *FileBuffer) MoveCursorRight() error {
	if fb.Cursor >= len(fb.File)-1 {
		return fmt.Errorf("Cannot move cursor to the right")
	}
	fb.Cursor++
	return nil
}

func (fb *FileBuffer) MoveCursorDown() error {
	newLine := byte('\n')
	after := 0
	for fb.File[fb.Cursor+after] != newLine && fb.Cursor+after < len(fb.File)-1 {
		after++
	}
	before := 0
	for fb.File[fb.Cursor-before] != newLine && fb.Cursor-before > 0 {
		before++
	}
	position := fb.Cursor + after + before
	if position >= len(fb.File) {
		return fmt.Errorf("Cannot move cursor down")
	}
	fb.Cursor = position
	return nil
}

func (fb *FileBuffer) MoveCursorLeft() error {
	if fb.Cursor <= 0 {
		return fmt.Errorf("Cannot move cursor to the left")
	}
	fb.Cursor--
	return nil
}

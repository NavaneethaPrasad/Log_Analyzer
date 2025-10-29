package segmenter

import (
	"os"
	"path/filepath"
	"testing"
)

func createTempLogFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp log file : %v\n", err)
	}
	return path
}
func TestParseLogSegments(t *testing.T) {
	tempDir := t.TempDir()

	content := `2025-10-23 15:04:10.001 | INFO | auth | host=db01 | request_id=req-001 | msg="User login successful"
2025-10-27 10:01:00.456 | ERROR | database | host=server2 | request_id=req2 | msg="Failed to connect"`

	createTempLogFile(t, tempDir, "sample.log", content)
	logStore, err := ParseLogSegments(tempDir)
	if err != nil {
		t.Fatalf("ParseLogSegments failed: %v", err)
	}

	if len(logStore.Segment) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(logStore.Segment))
	}

	segment := logStore.Segment[0]

	if segment.FileName != "sample.log" {
		t.Errorf("Expected file name 'sample.log', got '%s'", segment.FileName)
	}

	if len(segment.LogEntries) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(segment.LogEntries))
	}

	if segment.StartTime.IsZero() || segment.EndTime.IsZero() {
		t.Errorf("Start or End time is not set properly")
	}

	if segment.StartTime.After(segment.EndTime) {
		t.Errorf("StartTime %v should not be after EndTime %v", segment.StartTime, segment.EndTime)
	}
}

func TestParseLogSegments_BadDir(t *testing.T) {
	_, err := ParseLogSegments("nonexistent_dir")
	if err == nil {
		t.Errorf("Expected error for nonexistent directory, got nil")
	}
}
func TestParseLogSegments_SkipSubdirectory(t *testing.T) {
	tmpDir := t.TempDir()
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)
	_, err := ParseLogSegments(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
func TestParseLogSegments_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "empty.log"), []byte("not a log line"), 0644)

	logstore, err := ParseLogSegments(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(logstore.Segment) != 0 {
		t.Errorf("Expected 0 segments, got %d", len(logstore.Segment))
	}
}

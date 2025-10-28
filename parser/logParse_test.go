package parser

import (
	"strings"
	"testing"
	"time"
)

func TestLogParseEntry(t *testing.T) {
	logentry := `2025-10-23 15:04:10.001 | DEBUG | auth | host=db01 | request_id=req-hyx6sa-8587 | msg="2FA verification completed"`
	entry, err := LogParseEntry(logentry)
	if err != nil {
		t.Errorf("Log Parsing Failed!")
	}
	expectedTime, _ := time.Parse("2006-01-02 15:04:05.000", "2025-10-23 15:04:10.001")
	if entry.Raw != logentry {
		t.Errorf("Expected raw to be %q but got %q", logentry, entry.Raw)
	}
	if !entry.Time.Equal(expectedTime) {
		t.Errorf("Expected time %v but got %v", expectedTime, entry.Time)
	}
	if entry.Level != "DEBUG" {
		t.Errorf("Expected DEBUG but got %s", entry.Level)
	}

	if entry.Component != "auth" {
		t.Errorf("Expected auth but got %s.\n", entry.Component)
	}
	if entry.Host != "db01" {
		t.Errorf("Expected db01 but got %s.\n", entry.Host)
	}
	if entry.Requestid != "req-hyx6sa-8587" {
		t.Errorf("Expected req-hyx6sa-8587 but got %s.\n", entry.Requestid)
	}
	if entry.Message != "2FA verification completed" {
		t.Errorf("Expected '2FA verification completed' but got %s.\n", entry.Message)
	}
}

func TestParseInvalidLogEntry(t *testing.T) {
	invalidLine := `invalid log line`
	_, err := LogParseEntry(invalidLine)
	if err == nil {
		t.Errorf("Expected error for invalid format but got none")
	}
}

func TestParseLogEntryBadTime(t *testing.T) {

	badTimeLine := `2025-10-23 15:17:08.636000 | WARN | api-server | host=worker01 | request_id=req-4leuyy-5910 | msg="Cache cleared"`
	_, err := LogParseEntry(badTimeLine)
	if err == nil || !strings.Contains(err.Error(), "failed to parse time") {
		t.Errorf("Expected time parsing error, got %v", err)
	}
}

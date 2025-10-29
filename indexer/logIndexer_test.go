package indexer

import (
	"loggenerator/model"
	"reflect"
	"testing"
	"time"
)

func TestBuildSegmentIndex(t *testing.T) {
	entries := []model.LogEntry{
		{
			Time:      time.Now(),
			Level:     model.INFO,
			Component: "api-server",
			Host:      "worker01",
			Requestid: "req-1",
			Message:   "started service",
		},
		{
			Time:      time.Now(),
			Level:     model.ERROR,
			Component: "db",
			Host:      "worker02",
			Requestid: "req-2",
			Message:   "database connection failed",
		},
		{
			Time:      time.Now(),
			Level:     model.INFO,
			Component: "api-server",
			Host:      "worker01",
			Requestid: "req-3",
			Message:   "processed request",
		},
	}

	index := BuildSegmentIndex(entries)

	// ByLevel
	if got, want := index.ByLevel["INFO"], []int{0, 2}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByLevel[INFO] = %v, want %v", got, want)
	}
	if got, want := index.ByLevel["ERROR"], []int{1}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByLevel[ERROR] = %v, want %v", got, want)
	}

	//ByComponent
	if got, want := index.ByComponent["api-server"], []int{0, 2}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByComponent[api-server] = %v, want %v", got, want)
	}
	if got, want := index.ByComponent["db"], []int{1}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByComponent[db] = %v, want %v", got, want)
	}

	//ByHost
	if got, want := index.ByHost["worker01"], []int{0, 2}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByHost[worker01] = %v, want %v", got, want)
	}
	if got, want := index.ByHost["worker02"], []int{1}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByHost[worker02] = %v, want %v", got, want)
	}

	//ByReqID
	if got, want := index.ByReqId["req-1"], []int{0}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByReqID[req-1] = %v, want %v", got, want)
	}
	if got, want := index.ByReqId["req-3"], []int{2}; !reflect.DeepEqual(got, want) {
		t.Errorf("ByReqID[req-3] = %v, want %v", got, want)
	}
}

func TestBuildSegmentIndexEmptyInput(t *testing.T) {
	index := BuildSegmentIndex([]model.LogEntry{})

	if len(index.ByLevel) != 0 {
		t.Errorf("Expected empty ByLevel index, got %v", index.ByLevel)
	}
	if len(index.ByComponent) != 0 {
		t.Errorf("Expected empty ByComponent index, got %v", index.ByComponent)
	}
	if len(index.ByHost) != 0 {
		t.Errorf("Expected empty ByHost index, got %v", index.ByHost)
	}
	if len(index.ByReqId) != 0 {
		t.Errorf("Expected empty ByReqID index, got %v", index.ByReqId)
	}
}

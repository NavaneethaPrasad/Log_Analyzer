package segmenter

import (
	"bufio"
	"fmt"
	"loggenerator/indexer"
	"loggenerator/model"
	"loggenerator/parser"
	"os"
	"path/filepath"
)

func ParseLogSegments(s string) (model.LogStore, error) {
	LogStore := model.LogStore{
		Segment: []model.Segment{},
	}
	files, err := os.ReadDir(s)
	if err != nil {
		return LogStore, fmt.Errorf("failed to read directory : %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filepath := filepath.Join(s, file.Name())
		f, err := os.Open(filepath)
		if err != nil {
			fmt.Printf("Skipping file %s due to error: %v", filepath, err)
			continue
		}
		var LogEntries []model.LogEntry
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024) // allow 10MB lines

		for scanner.Scan() {
			line := scanner.Text()
			entry, err := parser.LogParseEntry(line)
			if err == nil {
				LogEntries = append(LogEntries, *entry)
			}
		}
		f.Close()
		if len(LogEntries) == 0 {
			continue
		}
		index := indexer.BuildSegmentIndex(LogEntries)
		segment := model.Segment{
			FileName:   file.Name(),
			LogEntries: LogEntries,
			StartTime:  LogEntries[0].Time,
			EndTime:    LogEntries[len(LogEntries)-1].Time,
			Index:      index,
		}
		LogStore.Segment = append(LogStore.Segment, segment)
	}
	return LogStore, nil
}

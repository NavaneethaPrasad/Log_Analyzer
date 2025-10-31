package main

import (
	"flag"
	"fmt"
	"log/slog"
	"loggenerator/filter"
	"loggenerator/segmenter"
	"strings"
	"time"
)

func main() {
	level := flag.String("level", "", "Filter by log level")
	component := flag.String("component", "", "Filter by component")
	host := flag.String("host", "", "Filter by host")
	reqID := flag.String("reqID", "", "Filter by requestID")
	startTimeStr := flag.String("start", "", "Filter logs from this start time (YYYY-MM-DD HH:MM:SS)")
	endTimestr := flag.String("end", "", "Filter logs up to this end time (YYYY-MM-DD HH:MM:SS)")

	flag.Parse()

	var startTime, endTime time.Time
	var err error

	if *startTimeStr != "" {
		startTime, err = time.Parse("2006-01-02 15:04:05", *startTimeStr)
		if err != nil {
			slog.Error("Error parsing time: ", "error", err)
		}
	}

	if *endTimestr != "" {
		endTime, err = time.Parse("2006-01-02 15:04:05", *endTimestr)
		if err != nil {
			slog.Error("Error parsing time: ", "error", err)
		}
	}

	logStore, err := segmenter.ParseLogSegments("../logs")
	if err != nil {
		slog.Error("Failed to parse logs\n")
	}
	split := func(s string) []string {
		if s == "" {
			return nil
		}
		parts := strings.Split(s, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	levels := split(*level)
	components := split(*component)
	hosts := split(*host)
	reqIDs := split(*reqID)

	filteredLogs := filter.FilterLogs(logStore, levels, components, hosts, reqIDs, startTime, endTime)
	fmt.Printf("Found %d matching entries\n", len(filteredLogs))
	for _, entry := range filteredLogs {
		fmt.Println(entry.Raw)
	}
}

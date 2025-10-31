package filter

import (
	"loggenerator/model"
	"time"
)

func FilterLogs(
	store model.LogStore,
	levels, components, hosts, reqIDs []string, startTime time.Time, endTime time.Time,
) []model.LogEntry {
	var result []model.LogEntry

	for _, segment := range store.Segment {
		totalFilters := 0

		if !startTime.IsZero() && segment.EndTime.Before(startTime) {
			continue
		}
		if !endTime.IsZero() && segment.StartTime.After(endTime) {
			continue
		}

		matchedIndex := make(map[int]bool)

		if len(levels) > 0 {
			totalFilters++
			for _, level := range levels {
				for _, idx := range segment.Index.ByLevel[level] {
					matchedIndex[idx] = true
				}
			}
		}

		if len(components) > 0 {
			totalFilters++
			componentFilter := make(map[int]bool)
			for _, component := range components {
				for _, idx := range segment.Index.ByComponent[component] {
					if matchedIndex[idx] || len(matchedIndex) == 0 {
						componentFilter[idx] = true
					}
				}
			}
			matchedIndex = componentFilter
		}

		if len(hosts) > 0 {
			totalFilters++
			hostFilter := make(map[int]bool)
			for _, host := range hosts {
				for _, idx := range segment.Index.ByHost[host] {
					if matchedIndex[idx] || len(matchedIndex) == 0 {
						hostFilter[idx] = true
					}
				}
			}
			matchedIndex = hostFilter
		}

		if len(reqIDs) > 0 {
			totalFilters++
			requestFilter := make(map[int]bool)
			for _, reqID := range reqIDs {
				for _, idx := range segment.Index.ByReqId[reqID] {
					if matchedIndex[idx] || len(matchedIndex) == 0 {
						requestFilter[idx] = true
					}
				}
			}
			matchedIndex = requestFilter
		}

		if totalFilters == 0 {
			for _, entry := range segment.LogEntries {
				if !isWithinTimeRange(entry.Time, startTime, endTime) {
					continue
				}
				result = append(result, entry)
			}
			continue
		}

		for idx := range matchedIndex {
			entry := segment.LogEntries[idx]
			if !isWithinTimeRange(entry.Time, startTime, endTime) {
				continue
			}
			result = append(result, entry)
		}
	}

	return result
}

// --- Time Filtering Helper ---
func isWithinTimeRange(t, start, end time.Time) bool {
	// No start or end given
	if start.IsZero() && end.IsZero() {
		return true
	}

	// Only start time
	if !start.IsZero() && end.IsZero() {
		return !t.Before(start)
	}

	// Only end time
	if start.IsZero() && !end.IsZero() {
		return !t.After(end)
	}

	// Both start and end
	return !t.Before(start) && !t.After(end)
}

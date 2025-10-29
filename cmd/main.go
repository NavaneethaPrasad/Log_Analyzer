package main

import (
	"fmt"
	"log"
	"loggenerator/segmenter"
)

// func main() {
// 	logLine := `2025-10-23 15:04:10.001 | DEBUG | auth | host=db01 | request_id=req-hyx6sa-8587 | msg="2FA verification completed"`

// 	entry, err := parser.LogParseEntry(logLine)
// 	if err != nil {
// 		log.Fatal("Error:", err)
// 	}
// 	fmt.Println("Time:", entry.Time.Format("2006-01-02 15:04:05.000"))
// 	fmt.Println("Level:", entry.Level)
// 	fmt.Println("Component:", entry.Component)
// 	fmt.Println("Host:", entry.Host)
// 	fmt.Println("Request ID:", entry.Requestid)
// 	fmt.Println("Message:", entry.Message)

// }

// func main() {
// 	entries, _ := parser.LogParseFiles("../logs")
// 	for _, entry := range entries {
// 		fmt.Println(entry)
// 	}
// 	fmt.Println(len(entries))
// }

// func main() {

// 	logStore, _ := segmenter.ParseLogSegments("../logs")
// 	segment := logStore.Segment[0]
// 	fmt.Printf("File Name: %s\n", segment.FileName)
// 	fmt.Printf("Start Time: %v\n", segment.StartTime)
// 	fmt.Printf("End Time: %v\n", segment.EndTime)
// 	fmt.Printf("Number of Log Entries: %d\n", len(segment.LogEntries))

// 	fmt.Println("\n--- Log Entries ---")
// 	for _, entry := range segment.LogEntries {
// 		fmt.Printf("[%s] | %s | %s| %s | %s\n", entry.Level, entry.Component, entry.Host, entry.Requestid, entry.Message)
// 	}
// }

func main() {
	logStore, err := segmenter.ParseLogSegments("../logs")
	if err != nil {
		log.Fatalf("Error parsing log segments: %v", err)
	}

	if len(logStore.Segment) == 0 {
		fmt.Println("No log segments found.")
		return
	}
	segment := logStore.Segment[0]

	fmt.Println("\n========= Index Summary =========")

	fmt.Println("\nBy Level:")
	for level, indices := range segment.Index.ByLevel {
		fmt.Printf("  %s → %v\n", level, indices)
	}

	fmt.Println("\nBy Component:")
	for comp, indices := range segment.Index.ByComponent {
		fmt.Printf("  %s → %v\n", comp, indices)
	}

	fmt.Println("\nBy Host:")
	for host, indices := range segment.Index.ByHost {
		fmt.Printf("  %s → %v\n", host, indices)
	}

	fmt.Println("\nBy Request ID:")
	for req, indices := range segment.Index.ByReqId {
		fmt.Printf("  %s → %v\n", req, indices)
	}
}

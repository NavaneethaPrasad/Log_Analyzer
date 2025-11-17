package main

import (
	"fmt"
	"log/slog"
	databasemodel "loggenerator/pkg/database_model"
	ginhandler "loggenerator/pkg/gin"
	"loggenerator/pkg/parser"
	"os"

	"github.com/gin-gonic/gin"
)

// const dbUrl = "postgresql:///log_Analyzer?host=/var/run/postgresql/"
const dbUrl = "postgresql:///log_analyzer_db?host=/var/run/postgresql/"

func handleCommand(args []string) error {
	db, err := databasemodel.CreateDB(dbUrl)
	if err != nil {
		return err
	}
	switch args[0] {
	case "init":
		err := databasemodel.InitDb(db)
		if err != nil {
			return err
		}
	case "add":
		dirpath := args[1] // TBD : Handle case where no filename specified
		if dirpath == "" {
			slog.Error("Specify directory!")
		}
		entries, err := parser.LogParseFiles(dirpath)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			databasemodel.AddEntry(db, entry)
		}
		return nil
	case "query":
		query := args[1:]
		fmt.Println(query)
		// query := strings.Join(args[1:], " ")
		entries, err := databasemodel.Query(db, query)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			fmt.Println(entry)
		}
		slog.Info("Filtering successful!", "no. of entries:", len(entries))
		return nil

	case "web":
		r := gin.Default()
		r.LoadHTMLGlob("pkg/gin/templates/*")
		ginhandler.DBRef = db
		ginhandler.SetupRoutes(r)
		r.Run(":8080")
	default:
		return fmt.Errorf("unknown command: %s (expected: init | add | query)", args[0])

	}
	return nil

}

func main() {
	err := handleCommand(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in invocation %v", err)
		os.Exit(-1)
	}

}

package databasemodel

import (
	"fmt"
	"log"
	"loggenerator/model"
	"os"
	"regexp"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type queryComponent struct {
	key      string
	value    []string
	operator string
}

type LogLevel struct {
	ID    uint   `gorm:"primaryKey"`
	Level string `gorm:"size:5;unique;not null"`
}

type LogComponent struct {
	ID        uint   `gorm:"primaryKey"`
	Component string `gorm:"size:12;unique;not null"`
}

type LogHost struct {
	ID   uint   `gorm:"primaryKey"`
	Host string `gorm:"size:10;unique;not null"`
}

// type Entry struct {
// 	gorm.Model
// 	TimeStamp time.Time
// 	Level     string
// 	Component string
// 	Host      string
// 	RequestId string
// 	Message   string
// }

type Entry struct {
	gorm.Model
	TimeStamp time.Time

	LevelID uint
	Level   LogLevel `gorm:"foreignKey:LevelID"`

	ComponentID uint
	Component   LogComponent `gorm:"foreignKey:ComponentID"`

	HostID uint
	Host   LogHost `gorm:"foreignKey:HostID"`

	RequestId string
	Message   string
}

func (l Entry) String() string {
	if l.TimeStamp.IsZero() {
		return "Empty"
	}
	return fmt.Sprintf("%s : %s : %s : %s : %s : %s",
		l.TimeStamp.Format(time.RFC3339),
		l.Level.Level,
		l.Component.Component,
		l.Host.Host,
		l.RequestId,
		l.Message,
	)
}

func CreateDB(dbUrl string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error
			Colorful:                  true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, fmt.Errorf("Couldn't open database %v", err)
	}
	return db, nil
}

// func InitDb(db *gorm.DB) error {
// 	db.AutoMigrate(&Entry{})
// 	return nil
// }

func InitDb(db *gorm.DB) error {
	return db.AutoMigrate(
		&LogLevel{},
		&LogComponent{},
		&LogHost{},
		&Entry{},
	)
}

// func AddEntry(db *gorm.DB, e model.LogEntry) error {
// 	logs := Entry{
// 		TimeStamp: e.Time,
// 		Level:     string(e.Level),
// 		Component: e.Component,
// 		Host:      e.Host,
// 		RequestId: e.Requestid,
// 		Message:   e.Message,
// 	}
// 	ctx := context.Background()
// 	err := gorm.G[Entry](db).Create(ctx, &logs)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func getOrCreateLevel(db *gorm.DB, level string) uint {
	var l LogLevel
	db.Where("level = ?", level).First(&l)
	if l.ID == 0 {
		l.Level = level
		db.Create(&l)
	}
	return l.ID
}

func getOrCreateComponent(db *gorm.DB, comp string) uint {
	var c LogComponent
	db.Where("component = ?", comp).First(&c)
	if c.ID == 0 {
		c.Component = comp
		db.Create(&c)
	}
	return c.ID
}

func getOrCreateHost(db *gorm.DB, host string) uint {
	var h LogHost
	db.Where("host = ?", host).First(&h)
	if h.ID == 0 {
		h.Host = host
		db.Create(&h)
	}
	return h.ID
}

func AddEntry(db *gorm.DB, e model.LogEntry) error {

	levelID := getOrCreateLevel(db, string(e.Level))
	componentID := getOrCreateComponent(db, e.Component)
	hostID := getOrCreateHost(db, e.Host)

	logs := Entry{
		TimeStamp:   e.Time,
		LevelID:     levelID,
		ComponentID: componentID,
		HostID:      hostID,
		RequestId:   e.Requestid,
		Message:     e.Message,
	}

	return db.Create(&logs).Error
}

// func parseQuery(parts []string) ([]queryComponent, error) {
// 	var ret []queryComponent
// 	// pattern := `^(?P<key>\S+)=(?P<value>\S+)$`
// 	//pattern := `^(?P<key>[^\s=!<>]+)\s*(?P<operator>=|!=|>=|<=|>|<)\s*(?P<value>[^,\s]+(?:,[^,\s]+)*)$`
// 	pattern := `^(?P<key>[^\s=!<>]+)\s*(?P<operator>=|!=|>=|<=|>|<)\s*(?P<value>.+)$`
// 	r, _ := regexp.Compile(pattern)
// 	for _, part := range parts {
// 		matches := r.FindStringSubmatch(part)

// 		if matches == nil {
// 			return nil, fmt.Errorf("Invalid condition %s", part)
// 		}
// 		// Allow INFO|ERROR
// 		rawValue := matches[r.SubexpIndex("value")]
// 		rawValue = strings.ReplaceAll(rawValue, "|", ",")

// 		vals := strings.Split(rawValue, ",")

// 		ret = append(ret, queryComponent{
// 			key:      matches[r.SubexpIndex("key")],
// 			operator: matches[r.SubexpIndex("operator")],
// 			value:    vals,
// 		})
// 	}
// 	return ret, nil

// }

// func Query(db *gorm.DB, query []string) ([]Entry, error) {
// 	// parse the Query
// 	var ret []Entry
// 	parsed, err := parseQuery(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println("Conditions", parsed)
// 	q := db
// 	for _, c := range parsed {
// 		if len(c.value) == 1 {
// 			// single value
// 			fmt.Printf("Applying condition: %s %s %s\n", c.key, c.operator, c.value[0])
// 			q = q.Where(fmt.Sprintf("%s %s ?", c.key, c.operator), c.value[0])
// 		} else {
// 			//multiple values and operator is !=
// 			if c.operator == "!=" {
// 				fmt.Printf("Applying NOT IN condition: %s IN %v\n", c.key, c.value)
// 				q = q.Where(fmt.Sprintf("%s NOT IN ?", c.key), c.value)
// 			} else {
// 				// multi value and operator is =
// 				fmt.Printf("Applying IN condition: %s IN %v\n", c.key, c.value)
// 				q = q.Where(fmt.Sprintf("%s IN ?", c.key), c.value)
// 			}

// 		}
// 	}

// 	// Execute final query
// 	if err := q.Find(&ret).Error; err != nil {
// 		return nil, err
// 	}

// 	return ret, nil
// }

func parseQuery(parts []string) ([]queryComponent, error) {
	var ret []queryComponent

	pattern := `^(?P<key>[^\s=!<>]+)\s*(?P<operator>=|!=|>=|<=|>|<)\s*(?P<value>.+)$`
	r := regexp.MustCompile(pattern)

	for _, part := range parts {
		part = strings.TrimSpace(part)

		matches := r.FindStringSubmatch(part)
		if matches == nil {
			return nil, fmt.Errorf("invalid condition: %s", part)
		}

		// Allow INFO|ERROR
		rawValue := matches[r.SubexpIndex("value")]
		rawValue = strings.ReplaceAll(rawValue, "|", ",")

		vals := strings.Split(rawValue, ",")

		ret = append(ret, queryComponent{
			key:      matches[r.SubexpIndex("key")],
			operator: matches[r.SubexpIndex("operator")],
			value:    vals,
		})
	}

	return ret, nil
}

func Query(db *gorm.DB, query []string) ([]Entry, error) {
	var ret []Entry

	// Parse the query string
	parsed, err := parseQuery(query)
	if err != nil {
		return nil, err
	}

	fmt.Println("Parsed conditions:", parsed)

	q := db

	for _, c := range parsed {

		key := strings.ToLower(c.key)

		// Translate logical columns to foreign key columns
		switch key {

		case "level":
			// Convert values INFO → levelID
			var ids []uint
			for _, v := range c.value {
				var lvl LogLevel
				if err := db.First(&lvl, "level = ?", v).Error; err != nil {
					return nil, fmt.Errorf("unknown level '%s'", v)
				}
				ids = append(ids, lvl.ID)
			}
			c.key = "level_id"
			c.value = toStringSlice(ids)

		case "component":
			var ids []uint
			for _, v := range c.value {
				var comp LogComponent
				if err := db.First(&comp, "component = ?", v).Error; err != nil {
					return nil, fmt.Errorf("unknown component '%s'", v)
				}
				ids = append(ids, comp.ID)
			}
			c.key = "component_id"
			c.value = toStringSlice(ids)

		case "host":
			var ids []uint
			for _, v := range c.value {
				var h LogHost
				if err := db.First(&h, "host = ?", v).Error; err != nil {
					return nil, fmt.Errorf("unknown host '%s'", v)
				}
				ids = append(ids, h.ID)
			}
			c.key = "host_id"
			c.value = toStringSlice(ids)
		}

		// Apply WHERE condition
		if len(c.value) == 1 {
			q = q.Where(fmt.Sprintf("%s %s ?", c.key, c.operator), c.value[0])
		} else {
			if c.operator == "!=" {
				q = q.Where(fmt.Sprintf("%s NOT IN ?", c.key), c.value)
			} else {
				q = q.Where(fmt.Sprintf("%s IN ?", c.key), c.value)
			}
		}
	}
	q = q.
		Preload("Level").
		Preload("Component").
		Preload("Host")
	if err := q.Find(&ret).Error; err != nil {
		return nil, err
	}

	return ret, nil
}

// convert []uint → []string
func toStringSlice(nums []uint) []string {
	s := make([]string, len(nums))
	for i, n := range nums {
		s[i] = fmt.Sprint(n)
	}
	return s
}

func SplitUserFilter(input string) []string {
	var parts []string
	current := ""
	tokens := strings.Fields(input)

	for _, tok := range tokens {
		// If token contains an operator, then new condition
		if strings.Contains(tok, "=") ||
			strings.Contains(tok, ">=") ||
			strings.Contains(tok, "<=") ||
			strings.Contains(tok, ">") ||
			strings.Contains(tok, "<") {

			// Save previous condition
			if current != "" {
				parts = append(parts, current)
			}
			current = tok
		} else {
			// continuation (timestamps)
			current += " " + tok
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}
func GetAllLogs(db *gorm.DB) ([]Entry, error) {
	var result []Entry
	err := db.Preload("Level").
		Preload("Component").
		Preload("Host").
		Find(&result).Error
	return result, err
}

// func FilterLogs(db *gorm.DB, levels, components, hosts []string, requestID, timestampCond string) ([]Entry, error) {
// 	var queries []string

// 	// Levels
// 	if len(levels) > 0 && len(levels) < 4 { // skip if all selected
// 		queries = append(queries, "level="+strings.Join(levels, "|"))
// 	}

// 	// Components
// 	if len(components) > 0 && len(components) < 5 {
// 		queries = append(queries, "component="+strings.Join(components, "|"))
// 	}

// 	// Hosts
// 	if len(hosts) > 0 && len(hosts) < 5 {
// 		queries = append(queries, "host="+strings.Join(hosts, "|"))
// 	}

// 	// RequestID filter (direct equality)
// 	if strings.TrimSpace(requestID) != "" {
// 		queries = append(queries, "request_id="+requestID)
// 	}

// 	// Timestamp filter (operator + value)
// 	if strings.TrimSpace(timestampCond) != "" {
// 		queries = append(queries, "time_stamp "+timestampCond)
// 	}

// 	// If no filters, return all logs
// 	if len(queries) == 0 {
// 		return GetAllLogs(db)
// 	}

// 	// Call existing QueryDB which takes []string
// 	return Query(db, queries)
// }

func FilterLogs(
	DBRef *gorm.DB,
	levels []string,
	components []string,
	hosts []string,
	requestId string,
	startTime string,
	endTime string,
) ([]Entry, error) {

	query := DBRef.Model(&Entry{})

	if len(levels) > 0 {
		var levelIDs []uint
		DBRef.Model(&LogLevel{}).Where("level IN ?", levels).Pluck("id", &levelIDs)
		if len(levelIDs) > 0 {
			query = query.Where("level_id IN ?", levelIDs)
		}
	}

	if len(components) > 0 {
		var compIDs []uint
		DBRef.Model(&LogComponent{}).Where("component IN ?", components).Pluck("id", &compIDs)
		if len(compIDs) > 0 {
			query = query.Where("component_id IN ?", compIDs)
		}
	}
	if len(hosts) > 0 {
		var hostIDs []uint
		DBRef.Model(&LogHost{}).Where("host IN ?", hosts).Pluck("id", &hostIDs)
		if len(hostIDs) > 0 {
			query = query.Where("host_id IN ?", hostIDs)
		}
	}

	if requestId != "" {
		query = query.Where("request_id = ?", requestId)
	}

	if startTime != "" && endTime != "" {
		query = query.Where("time_stamp BETWEEN ? AND ?", startTime, endTime)
	} else if startTime != "" {
		query = query.Where("time_stamp >= ?", startTime)
	} else if endTime != "" {
		query = query.Where("time_stamp <= ?", endTime)
	}
	var entries []Entry
	result := query.
		Preload("Level").
		Preload("Component").
		Preload("Host").
		Order("time_stamp").
		Find(&entries)

	return entries, result.Error
}

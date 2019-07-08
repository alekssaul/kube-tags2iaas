package appsettings

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	// MaximumNumberofTags - Set to two because some cloud providers like Azure may have a limit on number of tags on an object
	MaximumNumberofTags = 2
)

// AppSettings - is application settings struct
type AppSettings struct {
	InfrastructureTags []InfrastructureTags `json:"tags"`
}

// InfrastructureTags - hold tags to assign to infrastructure
type InfrastructureTags struct {
	Key   string
	Value string
}

// GetAppSettings - Returns Application settings
func GetAppSettings() (Settings AppSettings, err error) {
	TagsFromEnv := os.Getenv("INFRASTRUCTURETAGS")
	log.Printf("INFRASTRUCTURETAGS set to: %s", TagsFromEnv)
	b := []byte(TagsFromEnv)
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return Settings, err
	}

	for k, v := range data {
		var i = InfrastructureTags{}
		i.Key = k
		i.Value = fmt.Sprintf("%s", v)
		Settings.InfrastructureTags = append(Settings.InfrastructureTags, i)
	}

	if len(Settings.InfrastructureTags) > MaximumNumberofTags {
		return Settings, fmt.Errorf("Found %v tags, expected less than %v", len(Settings.InfrastructureTags), MaximumNumberofTags)
	}

	return Settings, nil
}

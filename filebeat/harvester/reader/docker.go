package reader

import (
	"encoding/json"
	"time"

	"github.com/elastic/beats/libbeat/common"
)

type dockerLogEntry struct {
	Log    string    `json:"log"`
	Stream string    `json:"stream"`
	Time   time.Time `json:"time,string"`
}

type dockerJSONLog struct {
	reader Reader
	entry  dockerLogEntry
}

// NewDockerJSONLog creates a new reader that can decode Docker JSON logs.
func NewDockerJSONLog(r Reader) Reader {
	return &dockerJSONLog{
		reader: r,
	}
}

// Next decodes a Docker JSON log entry and returns the filled Line object.
func (r *dockerJSONLog) Next() (Message, error) {
	message, err := r.reader.Next()
	if err != nil {
		return message, err
	}

	if err := json.Unmarshal(message.Content, &r.entry); err != nil {
		message.AddFields(common.MapStr{
			JsonErrorKey: err.Error(),
		})
		return message, nil
	}

	message.Content = []byte(r.entry.Log)
	message.AddFields(common.MapStr{
		"@timestamp": r.entry.Time.Format(time.RFC3339Nano),
	})

	return message, nil
}

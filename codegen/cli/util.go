package cli

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/PapaCharlie/go-restli/codegen/models"
	"github.com/PapaCharlie/go-restli/codegen/schema"
	"github.com/pkg/errors"
)

type Snapshot struct {
	Models models.SnapshotModels `json:"models"`
	Schema *schema.Resource      `json:"schema"`
}

func LoadSnapshotFromFile(filename string) (*Snapshot, error) {
	snapshot := new(Snapshot)
	err := ReadJSONFromFile(filename, snapshot)
	if err != nil {
		return nil, err
	} else {
		return snapshot, nil
	}
}

func ReadJSONFromFile(filename string, s interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panicln("Could not read", filename, err)
	}

	err = json.Unmarshal(data, s)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

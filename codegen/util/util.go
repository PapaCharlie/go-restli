package util

import (
	"github.com/PapaCharlie/go-restli/codegen/models"
	"github.com/PapaCharlie/go-restli/codegen/schema"
)

type Snapshot struct {
	Models []*models.Model
	Schema *schema.Resource
}

func ReadFullSnapshot(filename string) *Snapshot {

}

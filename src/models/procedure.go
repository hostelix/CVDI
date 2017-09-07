package models

import (
	arangoDB "github.com/hostelix/aranGO"

	"github.com/merakiVE/CVDI/core/types"
	"github.com/merakiVE/CVDI/core/validator"
	"github.com/merakiVE/CVDI/core/tags"
)

type Lane struct {
	//ID            string `json:"id"`
	Name          string `json:"name"`
	PoolRef       string `json:"pool_ref"`
	ActivitiesRef []string `json:"activities_ref"`
}

type Activity struct {
	ID         string `json:"id" validate:"required" on_create:"set,auto_uuid"`
	Name       string `json:"name"`
	NeuronKey  string `json:"neuron_key"`
	ActionID   string `json:"action_id"`
	Sequence   int  `json:"sequence"`
	Type       string `json:"type"` //manually - automatic
	InputData  map[string]interface{} `json:"input_data"`
	OutputData map[string]interface{} `json:"output_data"`
}

type ProcedureModel struct {
	arangoDB.Document

	ID         string        `json:"id" validate:"required" on_create:"set,auto_uuid"`
	Owner      string        `json:"owner,omitempty" validate:"required"`
	Pool       string        `json:"pool"`
	Lanes      []Lane        `json:"lanes"`
	Activities []Activity    `json:"activities,omitempty" validate:"required"`

	types.Timestamps
	ErrorsValidation []map[string]string `json:"errors_validation,omitempty"`
}

func (this ProcedureModel) GetKey() string {
	return this.Key
}

func (this ProcedureModel) GetCollection() string {
	return "procedures"
}

func (this ProcedureModel) GetError() (string, bool) {
	return this.Message, this.Error
}

func (this ProcedureModel) GetValidationErrors() ([]map[string]string) {
	return this.ErrorsValidation
}

func (this *ProcedureModel) PreSave(c *arangoDB.Context) {

	v := validator.New()

	v.Validate(this)

	if v.IsValid() {

		ptag := tags.New()

		//Tag Process for model
		ptag.ProcessTags(this)

		//Process tag to struct slice
		for i := range this.Lanes {
			ptag.ProcessTags(&this.Lanes[i])
		}

		for i := range this.Activities {
			ptag.ProcessTags(&this.Activities[i])
		}

	} else {

		c.Err["error_validation"] = "Error validating model"
		this.ErrorsValidation = v.GetMessagesValidation()
	}

	return
}

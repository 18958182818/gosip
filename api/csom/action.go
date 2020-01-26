package csom

import (
	"bytes"
	"text/template"
)

// Action ...
type Action interface {
	String() string
	SetID(id int)
	GetID() int
	SetObjectID(objectID int)
	GetObjectID() int
}

type action struct {
	template string
	id       int
	objectID int
	err      error
}

// NewAction ...
func NewAction(template string) Action {
	a := &action{}
	a.template = template
	return a
}

func (a *action) String() string {
	a.err = nil

	template, err := template.New("action").Parse(a.template)
	if err != nil {
		a.err = err
		return a.template
	}

	data := &struct {
		ID       int
		ObjectID int
	}{
		ID:       a.GetID(),
		ObjectID: a.GetObjectID(),
	}

	var tpl bytes.Buffer
	if err := template.Execute(&tpl, data); err != nil {
		a.err = err
		return a.template
	}

	return trimMultiline(tpl.String())
}

func (a *action) SetID(id int) {
	a.id = id
}

func (a *action) GetID() int {
	return a.id
}

func (a *action) SetObjectID(objectID int) {
	a.objectID = objectID
}

func (a *action) GetObjectID() int {
	return a.objectID
}

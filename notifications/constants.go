package notifications

import (
	"assignment-2/corona"
)

const (
	// RootPath of notifications endpoints.
	RootPath string = corona.RootPath + "/notifications"

	// IdPattern is the path of any endpoint that takes one id parameter and otherwise is defined by it's http method.
	IdPattern string = "/{id:[a-zA-Z]+}"
)

// Triggers that a webhook waits for.
const (
	// TriggerOnTimeout is triggered when the timeout of a webhook is expires.
	TriggerOnTimeout string = "ON_TIMEOUT"
	// TriggerOnChange is triggered when the latest data on the field the webhook is interested in changes.
	TriggerOnChange string = "ON_CHANGE"
)

// Fields that a webhook cares about.
const (
	FieldStringency string = "stringency"
	FieldConfirmed  string = "confirmed"
)

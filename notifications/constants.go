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

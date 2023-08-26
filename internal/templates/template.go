package templates

import (
	"html/template"
)

var Templates = template.Must(template.ParseFiles(
	"../../internal/templates/register.html",
	"../../internal/templates/login.html",
	"../../internal/templates/profile.html",
	"../../internal/templates/edit-details.html",
))

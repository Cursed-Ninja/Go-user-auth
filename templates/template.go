package templates

import (
	"html/template"
)

var Templates = template.Must(template.ParseFiles(
	"../../templates/register.html",
	"../../templates/login.html",
	"../../templates/profile.html",
	"../../templates/edit-details.html",
))

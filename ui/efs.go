package ui

import "embed"

// Files contains all the contents of the ui/static/ directory
// because of the go:embed "comment directive".
// This also supports multiple paths: //go:embed "static/css" "static/img" "static/js".
// This also supports specific files: //go:embed "static/css/main.css" "static/img" "static/js"
// This also supports wildcard paths: //go:embed "static/css/*.css" "static/img" "static/js"
// This also supports files that start with a . Or _: //go:embed "all:static"
//
//go:embed "html" "static"
var Files embed.FS

package assets

import "embed"

//go:embed doc/swagger/**
var SwaggerFS embed.FS

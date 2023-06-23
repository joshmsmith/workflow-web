package utils

import (
	"os"
	"strings"
)

var log_level = strings.ToLower(os.Getenv("LOG_LEVEL"))

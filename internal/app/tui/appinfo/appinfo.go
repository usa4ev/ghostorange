package appinfo

import (
	"fmt"
)

var (
	Version = "0.0.1"
	BuildDate string
)

func AppInfo() string{
	return fmt.Sprintf("version: %s\nbuild date: %s\n", 
		Version, BuildDate)
}
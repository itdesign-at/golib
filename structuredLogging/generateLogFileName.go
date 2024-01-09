package structuredLogging

import (
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

// GenerateLogfileName generates a filename with a golang
// time layout as input. "user.Current" is replaced with the username
// of the user.Current() method.
//
// Example:
//
//	layout = /var/log/messenger-user.Current-2006-01.log
//
// returns
//
//	/var/log/messenger-root-2023-12.log
//
// With New() and Init():
//
//	logFilename := filepath.Join("/var/log",
//	  structuredLogging.GenerateLogfileName("messenger-user.Current-2006-01.log")
//	)
//	structuredLogging.New(logFilename).Init()
func GenerateLogfileName(layout string) string {
	var str string
	var username string

	if u, err := user.Current(); err == nil {
		username = u.Username
		if username == "" {
			username = u.Name
		}
		if username == "" {
			username = u.Uid
		}
	} else {
		uid := os.Getuid()
		username = strconv.Itoa(uid)
	}

	layout = strings.Replace(layout, "user.Current", username, -1)
	str = time.Now().Format(layout)

	return str
}

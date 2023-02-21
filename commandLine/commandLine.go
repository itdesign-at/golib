package commandLine

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/itdesign-at/golib/keyvalue"
	"gopkg.in/yaml.v3"
)

// CheckRootUserAndExit terminates the program when you are not "root".
// Can be switched off with
//
//	# export CLI_SKIPUSERCHECK=YES
func CheckRootUserAndExit() {
	myUser, _ := user.Current()

	// "root" check can be overruled with an external variable:
	// # export CLI_SKIPUSERCHECK=YES
	if myUser.Uid != "0" && myUser.Gid != "0" {
		if os.Getenv("CLI_SKIPUSERCHECK") != "YES" {
			_, _ = fmt.Fprintf(os.Stderr, "You must be 'root' to run this program\n")
			os.Exit(1)
		}
	}
}

// Parse is ported from PHP
//
//	/opt/watchit/lib/common/CommandLine.php
//
// and reads os.Args into the map returned
//
//	Examples:
//	./check_value -h server-1.demo.at -s "temperature" -value 17.3 -sensorOk -c 'eval(val("value")>25)'
//	all strings except "value" as float64 and "sensorOk" as bool
func Parse(args []string) keyvalue.Record {
	var opt = make(keyvalue.Record)
	n := len(args)
	if n < 2 {
		return opt
	}
	var key, value string
	for i := 1; i < n; i++ {
		key = ""
		if strings.HasPrefix(args[i], "-") {
			key = strings.TrimSpace(strings.TrimLeft(args[i], "-"))
		}
		if key == "" {
			if i+1 == n { // end reached?
				opt[args[i]] = true
				break
			}
			continue
		}
		if i+1 == n { // end reached?
			opt[key] = true
			break
		}
		value = args[i+1]
		// first character is a mask character
		// e.g. -negative "\-3"
		if strings.HasPrefix(value, `\`) {
			value = value[1:]
			if value == "" {
				opt[key] = `\`
			} else {
				if f, err := strconv.ParseFloat(value, 64); err == nil {
					opt[key] = f
				} else {
					opt[key] = value
				}
			}
			i++
			continue
		}
		// if "-key1" follows "-key2"
		if strings.HasPrefix(value, `-`) {
			opt[key] = true
			continue
		}
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			opt[key] = f
		} else {
			switch value {
			case "true":
				opt[key] = true
			case "false":
				opt[key] = false
			default:
				opt[key] = value
			}
		}
		i++
	}
	return opt
}

func PrintVersion(args keyvalue.Record, detailed bool) {
	if detailed {
		b, _ := yaml.Marshal(args)
		fmt.Print(string(b))
	} else {
		fmt.Println(args.String("version"))
	}
}

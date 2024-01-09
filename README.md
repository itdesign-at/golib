# golib
you have to prefix the tag with the folder name, e.g.: commandLine/v1.0.0

    git tag -a commandLine/v1.0.0 -m "Release 1.0.0"
    git tag -a converter/v1.0.0 -m "Release 1.0.0"
    git tag -a encoding/v1.0.0 -m "Release 1.0.0"
    git tag -a keyvalue/v1.0.0 -m "Release 1.0.0"
    git tag -a structuredLogging/v1.0.0 -m "Release 1.0.0"
    git push --tags

See [go module documentation](https://go.dev/doc/modules/managing-source) for more information.

# package commandLine

| Function                                               | Comment                                                             |
|--------------------------------------------------------|---------------------------------------------------------------------|
| func CheckRootUserAndExit()                            | CheckRootUserAndExit terminates the program when you are not "root" |
| func Parse(args []string) keyvalue.Record              | Reads os.Args into the map returned                                 |
| func PrintVersion(args keyvalue.Record, detailed bool) | Prints version info                                                 |

# package converter

| Function                                | Comment                                                   |
|-----------------------------------------|-----------------------------------------------------------|
| func Normalize(in string) string        | Normalize a string by removing all special characters     |
| func Denormalize(in string) string      | Denormalize does the opposite of Normalize                |
| func GetRandomString(length int) string | GetRandomString returns a random string with length given |

# package encoding

| Function                                       | Comment                            |
|------------------------------------------------|------------------------------------|
| func MarshalFile(file string, in any) error    | encode yaml, json or gob into file |
| func UnmarshalFile(file string, out any) error | decode yaml, json or gob from file |


# package keyvalue

| Function                                                     | Comment                         |
|--------------------------------------------------------------|---------------------------------|
| func NewRecord() Record                                      | does make for a new record      |
| func (r Record) Exists(key string) bool                      | checks the existence of the key |
| func (r Record) Value(key string) interface{}                | returns the value               |
| func (r Record) Bool(key string, convert ...bool) bool       | returns a bool value            |
| func (r Record) Float64(key string, convert ...bool) float64 | returns a float value           |
| func (r Record) Int(key string, convert ...bool) int         | returns an int value            |
| func (r Record) Int64(key string, convert ...bool) int64     | returns an int64 value          |
| func (r Record) String(key string, convert ...bool) string   | returns a string value          |
| func (r Record) Copy() Record                                | copy a record                   |
| func (r Record) GetSortedKeys() []string                     | returns keys sorted             |

# package structuredLogging

| Function                                                            | Comment                        |
|---------------------------------------------------------------------|--------------------------------|
| func New(dsn ...string) *SlogLogger                                 | New creates an slog writer     |
| func (sl *SlogLogger) Parameter(params keyvalue.Record) *SlogLogger | add params                     |
| func (sl *SlogLogger) InitJsonHandler() *slog.JSONHandler           | Get JSONHandler with writer(s) |
| func (sl *SlogLogger) Init() *SlogLogger                            | Init the logger with handler   |
| func (sl *SlogLogger) Write(p []byte) (int, error)                  | Write to all log writers       |
| func GenerateLogfileName(layout string) string                      | Get a logfilename              |                 

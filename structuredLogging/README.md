package structuredLogging

| Function                                                            | Comment                        |
|---------------------------------------------------------------------|--------------------------------|
| func New(dsn ...string) *SlogLogger                                 | New creates an slog writer     |
| func (sl *SlogLogger) Parameter(params keyvalue.Record) *SlogLogger | add params                     |
| func (sl *SlogLogger) InitJsonHandler() *slog.JSONHandler           | Get JSONHandler with writer(s) |
| func (sl *SlogLogger) Init() *SlogLogger                            | Init the logger with handler   |
| func (sl *SlogLogger) Write(p []byte) (int, error)                  | Write to all log writers       |
| func GenerateLogfileName(layout string) string                      | Get a logfilename              |                 
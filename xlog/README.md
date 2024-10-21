# package xlog

>Structured logging with slog. Supports logging to console, log files and nats.

## Log Record

>Every log record has these default slog attributes:
    
- `time` 
- `level`
- `msg`

>Every log record must have these additional attributes:

- `category`: e.g.: system, security, user, scheduler, ... 
- `action`: e.g.: startup, shutdown, login, failed login, job start, job end, ...
- `hostname`

>Every log record can have these additional attributes:

- `user`: authenticated user
- `extended`: additional information in key-value pairs

## Usage examples

### Initialization

```go   
    
    func initLogging(destinations []string, level string, addSource bool) *xlog.Client {

	if len(destinations) == 0 {
		destinations = []string{"stdout"}
	}

	client, err := xlog.NewClient(xlog.ClientOptions{
		Destinations: destinations,
		Level:        level,
		AddSource:    addSource,
	})
	if err != nil {
		log.Fatal("Logging initialization failed: ", err)
	}

	slog.Info("Logging initialized", xlog.Category("system"), xlog.Action("startup"), xlog.Extended("destinations", destinations, "level", level, "addSource", addSource))

	return client
}
```

### Logging
```go

        // Log a message 
	slog.Info("ITdesign WATCH IT initialized", xlog.Category("system"), xlog.Action("startup"), xlog.Extended("version", VERSION))

        // Log an error 
	slog.Error("Loading api.yaml failed", xlog.Category("system"), xlog.Action("startup"), xlog.Error(err))	

        // Audit log 
	slog.Log(ctx.Context(), xlog.LevelAudit, "User created", xlog.Category("user"), xlog.Action("create"), xlog.User(acl.GetUserName(ctx)), xlog.Extended(u))


```
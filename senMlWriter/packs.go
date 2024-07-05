package senMlWriter

import (
	"log/syslog"
	"math"
	"slices"
	"sync"
	"time"

	"github.com/mainflux/senml"

	"github.com/itdesign-at/golib/keyvalue"
)

// Config holds the configuration for the Sensorml handler.
type Config struct {
	// field(s) for senML
	BaseName string `json:"baseName" yaml:"baseName"`

	// field(s) for sending data
	Out string `json:"out" yaml:"out"`

	SyslogPriority syslog.Priority `json:"syslogPriority" yaml:"syslogPriority"`

	// TimePrecision is the number of digits after the decimal point.
	// The precision is used to round the time to reduce digits.
	TimePrecision int `json:"timePrecision" yaml:"timePrecision"`

	// FlushInterval is the interval in seconds to write the data to the Out.
	FlushInterval int `json:"flushInterval" yaml:"flushInterval"`
}

// Handler is a Sensorml handler.
// The handler is used to add data to the Sensorml handler.
type Handler struct {
	sync.Mutex

	//bn contains the senML records for each base name
	config        Config
	packs         map[string]senml.Pack
	writtenMinute int
	stop          chan struct{}
	done          chan struct{}
}

// New creates a new Sensorml handler with the given configuration.
// The handler is used to add data to the Sensorml handler.
// The handler writes the data to the configured Out every minute.
func New(c Config) *Handler {
	if c.SyslogPriority == 0 {
		c.SyslogPriority = syslog.LOG_INFO | syslog.LOG_LOCAL7
	}

	if c.FlushInterval == 0 {
		c.FlushInterval = 60
	}

	h := &Handler{
		config:        c,
		packs:         make(map[string]senml.Pack),
		writtenMinute: -1,
		stop:          make(chan struct{}),
		done:          make(chan struct{}),
	}

	go h.scheduler()
	return h
}

// Add adds the given data to the Sensorml handler.
// The base name is optional and is used as the base name for the senml records.
// If the base name is not set, the base name from the configuration is used.
func (h *Handler) Add(t time.Time, d map[string]any, baseName ...string) *Handler {
	h.Lock()
	defer h.Unlock()

	bn := h.config.BaseName
	if len(baseName) > 0 {
		bn = baseName[0]
	}

	if pack := h.add(t, d, bn); len(pack.Records) > 0 {
		h.packs[bn] = pack
	}

	return h
}

// Close writes the data to the configured Out.
// It returns an error if the writing fails.
func (h *Handler) Close() error {
	h.stop <- struct{}{}
	<-h.done
	return nil
}

// Flush writes the data to the configured Out.
// It returns an error if the writing fails.
func (h *Handler) Flush() error {
	var lastErr error

	h.Lock()
	defer h.Unlock()

	for bn, p := range h.packs {
		cfg := WriterConfig{
			BaseName:       bn,
			Out:            h.config.Out,
			SyslogPriority: h.config.SyslogPriority,
		}

		// if the writing fails, the data is kept in the handler
		// and will be written in the next minute
		if err := NewWriter(cfg).AddPack(p).Write(); err != nil {
			lastErr = err
		} else {
			delete(h.packs, bn)
		}
	}

	return lastErr
}

// add adds the given data to the Sensorml handler.
// The data are sorted by key and added to the pack.
func (h *Handler) add(t time.Time, data keyvalue.Record, baseName string) senml.Pack {

	if len(data) == 0 {
		return senml.Pack{}
	}

	// get the first key of the data map
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	var timeDelta float64
	firstKey := keys[0]
	firstValue := data.Float64(firstKey, true)

	pack, ok := h.packs[baseName]
	// The first record contains the first key and a time stamp.
	// If the device is not in the map, we have to set the base time.
	// If the device is in the map, we have to calculate the time delta.
	if !ok {
		timeDelta = 0
		pack.Records = append(pack.Records,
			senml.Record{
				BaseName: baseName,
				BaseTime: round(float64(t.UnixMilli())/1000, h.config.TimePrecision),
				Name:     firstKey,
				Value:    &firstValue,
			})
	} else {
		// round to x decimal place to reduce digits
		timeDelta = round(float64(t.UnixMilli())/1000-pack.Records[0].BaseTime, h.config.TimePrecision)
		pack.Records = append(pack.Records,
			senml.Record{
				Time:  timeDelta,
				Name:  firstKey,
				Value: &firstValue,
			})
	}

	// add all other records, ignore first key
	for _, key := range keys[1:] {

		v := data.Float64(key)
		r := senml.Record{
			Name:  key,
			Value: &v,
			Time:  timeDelta,
		}

		pack.Records = append(pack.Records, r)
	}

	return pack
}

// scheduler writes the data to the configured Out every minute.
// The scheduler is stopped by closing the stop channel.
// The data is written either every full second or minute (depending on the configuration).
func (h *Handler) scheduler() {

	interval := time.Second * time.Duration(h.config.FlushInterval)
	wait := time.Until(time.Now().Truncate(interval).Add(interval))

	timer := time.NewTimer(wait)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			timer.Stop() // Stop the initial timer after first tick
			ticker.Reset(interval)
			h.Flush()
		case <-ticker.C:
			h.Flush()
		case <-h.stop:
			h.Flush()
			close(h.done)
			return
		}
	}
}

// Round rounds the given float to the given precision.
// The precision is the number of digits after the decimal point.
func round(val float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(val*factor) / factor
}

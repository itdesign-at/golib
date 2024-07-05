package senMlWriter

import (
	"fmt"
	"github.com/mainflux/senml"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_add_WithEmptyDataReturnsEmptyPack(t *testing.T) {
	handler := New(Config{TimePrecision: 2, FlushInterval: 1})
	handler.Close()

	result := handler.add(time.Now(), map[string]any{}, "testBaseName")
	assert.Equal(t, senml.Pack{}, result)
}

func TestAddWithSingleRecord(t *testing.T) {
	handler := New(Config{TimePrecision: 2, FlushInterval: 1})
	handler.Close()

	data := map[string]any{"temperature": 22.5}
	timeStamp := time.Now()

	result := handler.add(timeStamp, data, "environment")

	assert.Equal(t, 1, len(result.Records))
	assert.Equal(t, "temperature", result.Records[0].Name)
	assert.Equal(t, 22.5, *result.Records[0].Value)
	assert.Equal(t, "environment", result.Records[0].BaseName)
	assert.NotEqual(t, 0.0, result.Records[0].BaseTime)
	assert.Equal(t, 0.0, result.Records[0].Time)
}

func TestAddWithMultipleRecords(t *testing.T) {
	handler := New(Config{TimePrecision: 2, FlushInterval: 1})
	handler.Close()

	data := map[string]any{"temperature": 22.5, "humidity": 60}
	timeStamp := time.Now()

	result := handler.add(timeStamp, data, "environment")

	assert.Equal(t, 2, len(result.Records))
	// Assuming slices.Sort sorts in alphabetical order
	assert.Equal(t, "humidity", result.Records[0].Name)
	assert.Equal(t, 60.0, *result.Records[0].Value)
	assert.Equal(t, "environment", result.Records[0].BaseName)
	assert.NotEqual(t, 0.0, result.Records[0].BaseTime)
	assert.Equal(t, 0.0, result.Records[0].Time)

	assert.Equal(t, "temperature", result.Records[1].Name)
	assert.Equal(t, 22.5, *result.Records[1].Value)
	assert.Equal(t, "", result.Records[1].BaseName)
	assert.Equal(t, 0.0, result.Records[1].BaseTime)
	assert.Equal(t, 0.0, result.Records[1].Time)
}

func TestAddWithValidDataAddsPack(t *testing.T) {
	handler := New(Config{TimePrecision: 2, FlushInterval: 1})
	handler.Close()

	data := map[string]any{"temperature": 22.5}
	timeStamp := time.Now()
	baseName := "environment"

	handler.Add(timeStamp, data, baseName)

	assert.Contains(t, handler.packs, baseName)
	assert.NotEmpty(t, handler.packs[baseName].Records)
}

func TestAddWithNoBaseNameUsesConfigBaseName(t *testing.T) {
	configBaseName := "defaultBaseName"

	handler := New(Config{TimePrecision: 2, FlushInterval: 1, BaseName: configBaseName})
	handler.Close()

	data := map[string]any{"humidity": 60}
	timeStamp := time.Now()

	handler.Add(timeStamp, data)

	assert.Contains(t, handler.packs, configBaseName)
	assert.NotEmpty(t, handler.packs[configBaseName].Records)
}

func TestAddWithEmptyDataDoesNotCreatePack(t *testing.T) {
	handler := New(Config{TimePrecision: 2, FlushInterval: 1})
	handler.Close()

	timeStamp := time.Now()
	baseName := "emptyData"

	handler.Add(timeStamp, map[string]any{}, baseName)

	assert.NotContains(t, handler.packs, baseName)
}

func TestAddWithMultipleBaseNamesCreatesSeparatePacks(t *testing.T) {
	handler := New(Config{TimePrecision: 2, FlushInterval: 1})
	handler.Close()

	data1 := map[string]any{"temperature": 22.5}
	data2 := map[string]any{"humidity": 60}
	timeStamp := time.Now()
	baseName1 := "environment1"
	baseName2 := "environment2"

	handler.Add(timeStamp, data1, baseName1)
	handler.Add(timeStamp, data2, baseName2)

	assert.Contains(t, handler.packs, baseName1)
	assert.Contains(t, handler.packs, baseName2)
	assert.NotEqual(t, handler.packs[baseName1], handler.packs[baseName2])
}

func Test_Add_WithMultipleRecords(t *testing.T) {
	handler := New(Config{TimePrecision: 2, FlushInterval: 1})
	handler.Close()

	data1 := map[string]any{"temperature": 22.5, "humidity": 60}
	timeStamp1 := time.Now()
	handler.Add(timeStamp1, data1, "environment")

	data2 := map[string]any{"temperature": 23.5, "humidity": 61}
	timeStamp2 := timeStamp1.Add(5 * time.Second)
	handler.Add(timeStamp2, data2, "environment")

	assert.Contains(t, handler.packs, "environment")

	assert.Equal(t, 4, len(handler.packs["environment"].Records))
	// Assuming slices.Sort sorts in alphabetical order
	assert.Equal(t, "humidity", handler.packs["environment"].Records[0].Name)
	assert.Equal(t, 60.0, *handler.packs["environment"].Records[0].Value)
	assert.Equal(t, "environment", handler.packs["environment"].Records[0].BaseName)
	assert.NotEqual(t, 0.0, handler.packs["environment"].Records[0].BaseTime)
	assert.Equal(t, 0.0, handler.packs["environment"].Records[0].Time)

	assert.Equal(t, "temperature", handler.packs["environment"].Records[1].Name)
	assert.Equal(t, 22.5, *handler.packs["environment"].Records[1].Value)
	assert.Equal(t, "", handler.packs["environment"].Records[1].BaseName)
	assert.Equal(t, 0.0, handler.packs["environment"].Records[1].BaseTime)
	assert.Equal(t, 0.0, handler.packs["environment"].Records[1].Time)

	assert.Equal(t, "humidity", handler.packs["environment"].Records[2].Name)
	assert.Equal(t, 61.0, *handler.packs["environment"].Records[2].Value)
	assert.Equal(t, "", handler.packs["environment"].Records[2].BaseName)
	assert.Equal(t, 0.0, handler.packs["environment"].Records[2].BaseTime)
	assert.Equal(t, 5.0, handler.packs["environment"].Records[2].Time)

	assert.Equal(t, "temperature", handler.packs["environment"].Records[3].Name)
	assert.Equal(t, 23.5, *handler.packs["environment"].Records[3].Value)
	assert.Equal(t, "", handler.packs["environment"].Records[3].BaseName)
	assert.Equal(t, 0.0, handler.packs["environment"].Records[3].BaseTime)
	assert.Equal(t, 5.0, handler.packs["environment"].Records[3].Time)

}

func Test_Example1(t *testing.T) {
	handler := New(Config{
		TimePrecision: 2,
		FlushInterval: 60,
		BaseName:      "hu/train1/wagon1/advantech3231/e0673922ebd5/gps/",
		Out:           "syslog://127.0.0.1:7814/senml02",
	})
	defer handler.Close()

	data1 := map[string]any{"latitude": 47.601987, "longitude": 17.249188, "altitude": 129.3, "speed": 64.2}
	data2 := map[string]any{"latitude": 47.601985, "longitude": 17.249185, "altitude": 128.3, "speed": 50.21}
	data3 := map[string]any{"latitude": 47.601984, "longitude": 17.249184, "altitude": 127.3, "speed": 20.9}
	data4 := map[string]any{"latitude": 47.601980, "longitude": 17.249183, "altitude": 120.9, "speed": 0}

	handler.Add(time.Now(), data1)
	handler.Add(time.Now().Add(time.Second), data2)
	handler.Add(time.Now().Add(2*time.Second), data3)
	handler.Add(time.Now().Add(3*time.Second), data4)

	fmt.Println("hu/train1/wagon1/advantech3231/e0673922ebd5/gps/", len(handler.packs["hu/train1/wagon1/advantech3231/e0673922ebd5/gps/"].Records), "values added")
}

func Test_Example2(t *testing.T) {
	handler := New(Config{
		TimePrecision: 2,
		FlushInterval: 60,
		Out:           "syslog://127.0.0.1:7814/senml02",
	})
	defer handler.Close()

	data1 := map[string]any{"latitude": 47.601987, "longitude": 17.249188, "altitude": 129.3, "speed": 64.2}
	data2 := map[string]any{"Ch1": 0.56985, "Ch2": 1.72185, "Ch3": 0, "Ch4": -14.5021}
	data3 := map[string]any{"latitude": 47.601984, "longitude": 17.249184, "altitude": 127.3, "speed": 20.9}
	data4 := map[string]any{"Ch1": 0.585, "Ch2": 1.185, "Ch3": 0.1, "Ch4": -14.21}

	handler.Add(time.Now(), data1, "hu/train1/wagon1/advantech3231/e0673922ebd5/gps/")
	handler.Add(time.Now().Add(time.Second), data2, "hu/train1/wagon1/advantech3231/e0673922ebd5/wago/")
	handler.Add(time.Now().Add(2*time.Second), data3, "hu/train1/wagon1/advantech3231/e0673922ebd5/gps/")
	handler.Add(time.Now().Add(3*time.Second), data4, "hu/train1/wagon1/advantech3231/e0673922ebd5/wago/")

	fmt.Println("hu/train1/wagon1/advantech3231/e0673922ebd5/gps/", len(handler.packs["hu/train1/wagon1/advantech3231/e0673922ebd5/gps/"].Records), "values added")
	fmt.Println("hu/train1/wagon1/advantech3231/e0673922ebd5/gps/", len(handler.packs["hu/train1/wagon1/advantech3231/e0673922ebd5/wago/"].Records), "values added")
}

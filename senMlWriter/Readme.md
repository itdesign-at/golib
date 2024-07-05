# Examples
## only one baseName
This is a simple example when only a BaseName is used. 
The BaseName is specified in the configuration, and the data is added using the Add method. 
The FlushInterval parameter defines how often data is sent to the server or saved as a file. T
he Out parameter specifies the server to which the data should be sent. 
The TimePrecision parameter specifies the number of decimal places for the time.

### Info: 
The Close() function closes the connection to the server, sends any outstanding data, and stops the automated data sending. 
The Flush() function sends the data to the server, but the handler remains open.

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


## multiple baseNames
The BaseName can also be specified with each Add call. 
This is useful when the data comes from different sensors and has different BaseNames. 
In this example, data from two sensors is sent. The first sensor sends GPS data, and the second sensor sends Wago data. 
The data is stored in two different BaseNames and sent to the server.

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

package main

import (
	"log"

	"io"

	"strings"

	"bufio"

	"github.com/tarm/serial"
)

/**
Reverse engineered commands:
M - get serial number
S - Next sample? After the initial is recieved by a button press
E - End sample?
R - Reboot/Reset device

EOL is 0a (\n I think)

I think the way it works is:
- Click button, device sends
Data:\r\n
1 value
- App sends S
- Device sends
1 value
- App sends S
- Device sends
1 value
...
- App sends E which makes the device listen for the button click again
*/

/*
Responses:
button click - Data\r\n
S - ##\r\n
*/

//0d = \r carriage return
//0a = \n linefeed

func receive(pipe chan string, port io.ReadWriteCloser) {

	portBuff := make([]byte, 100)

	for {

		//var buff bytes.Buffer

		var err error

		reader := bufio.NewReader(port)

		portBuff, err = reader.ReadBytes('\x0a') //io.ReaderFrom() //.ReadFull(port, portBuff) //port.Read(portBuff) //port.Read(buff)

		if string(portBuff[:len(portBuff)-1]) == "\n" {
			//done, ship it
			break
		}

		//remove newline and cr characters
		str := strings.Replace(string(portBuff), "\r", "", -1)
		str = strings.Replace(str, "\n", "", -1)

		output := strings.Split(str, "\n")

		for _, s := range output {
			//send the output into the channel
			pipe <- s
		}

		if err != nil {
			log.Fatal(err)
			break
		}

	}

}

func main() {

	output := make(chan string)
	port, err := setup()
	if err != nil {
		log.Fatal(err)
	}

	defer teardown(port)

	go receive(output, port)

	//port.Write([]byte("S"))

	//o := <-output
	//
	//if o != "" {
	//	log.Println(o)
	//}

	port.Write([]byte("E"))

	port.Write([]byte("S"))

	for {
		o := <-output

		if o != "" {
			log.Println(o)
		}
	}

}

func setup() (io.ReadWriteCloser, error) {

	c := &serial.Config{Name: "COM5", Baud: 9600}

	// Open the port.
	port, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}

	return port, nil
}

func teardown(port io.ReadWriteCloser) {
	// Make sure to close it later.
	port.Close()
}

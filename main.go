package main

import (
	"log"

	"io"

	"github.com/tarm/serial"
)

/**
Reverse engineered commands:
M - get serial number
S - Next sample? After the initial is recieved by a button press
E - End sample?
R - Reboot/Reset device


I think the way it works is:
- Click button, device sends
Data:\n\r
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
button click - Data\n\r
S - ##\n\r
*/

func receive(output chan string, port io.ReadWriteCloser) {

	buff := make([]byte, 128)

	for {

		var s string

		s = ""

		//TODO: figure out why this is only reading a single byte. Or loop until this has read the entire buffer, smash it together, split on \n\r, and then send each chunk to output
		n, err := port.Read(buff) //port.Read(buff)
		log.Printf("read")

		if n > 0 {
			s = string(buff[:n])
		}

		//send the average to the output
		output <- s

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

	//port.Write([]byte("S"))

	for {
		o := <-output

		if o != "" {
			log.Println(o)
		}
	}

	//X++
	//if X > MaxX {
	//	X = 1
	//}
	//
	//Y++
	//if Y > MaxY {
	//	Y = 1
	//}

	//log.Printf("X %x Y %x Value %x", X, Y, o)
	//}

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

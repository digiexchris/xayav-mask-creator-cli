package main

import (
	"log"

	"io"

	"strings"

	"bufio"

	"strconv"

	"time"

	"github.com/tarm/serial"
)

/**
Discovered commands:
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

type Output struct {
	//The nth sampling session, or one button press cycle
	Cell int
	//the averaged value for this cell
	Value int
}

func receive(pipe chan string, port io.ReadWriteCloser) {

	portBuff := make([]byte, 100)

	for {

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

/*
Managing responses back to the device
It appears that in order for more samples to happen after the button is clicked and one sample is sent,
we need to respond with an S for each additional sample or E to reset the button state.

If we receive Data:\r\n then read the following sample. Once the following sample is received, send a sample and receive it,
repeating sampling until we've reached the max samples. Once we've received the last sample and we're at max samples,
send E, average all of the samples, and send the resulting cell number and it's averaged value  to output
*/
func controller(pipe chan string, outputChannel chan Output, port io.ReadWriteCloser) {
	//number of samples to take before averaging and sending to output
	const maxSamples = 96
	//sample slower than once every 20ms
	const sampleMinInterval = time.Millisecond * 5

	var samples []int
	var cell int
	var nSamplesToRead int
	nSamplesToRead = 0
	cell = 0

	for {
		o := <-pipe

		if o == "Data: " {
			cell++

			// -1 because we get one sample immediately on button press, it's already in the pipe
			nSamplesToRead = maxSamples - 1
			samples = nil

			//there is a sample in the pipe following Data, so read it
			sample, err := strconv.Atoi(<-pipe)
			if err != nil {
				log.Fatal(err)
				break
			}

			samples = append(samples, sample)

			//trigger the next sample after waiting 100ms
			time.Sleep(sampleMinInterval)
			send(port, "S")
		} else if sample, err := strconv.Atoi(o); err == nil {
			//if it's a number it's a sample
			samples = append(samples, sample)
			nSamplesToRead--
			if nSamplesToRead < 1 {
				//If we're done this set of sampling
				var total int = 0
				for _, value := range samples {
					total += value
				}

				average := total / len(samples)

				output := Output{
					Cell:  cell,
					Value: average,
				}

				outputChannel <- output

				//we're done!
				send(port, "E")
			} else {
				time.Sleep(sampleMinInterval)
				send(port, "S")
			}
		} else {
			//it's something else other than Data: or a sample, so it's probably device info
			log.Println(o)
			//we're done!
			send(port, "E")
		}

	}
}

func send(port io.ReadWriteCloser, command string) {

	_, err := port.Write([]byte(command))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	reciever := make(chan string)
	output := make(chan Output)

	port, err := setup()
	if err != nil {
		log.Fatal(err)
	}

	defer teardown(port)

	go receive(reciever, port)
	go controller(reciever, output, port)

	//port.Write([]byte("S"))

	//o := <-output
	//
	//if o != "" {
	//	log.Println(o)
	//}

	for {
		o := <-output

		if o != (Output{}) {
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

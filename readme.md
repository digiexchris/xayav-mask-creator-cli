# Xayav Mask Creater CLI

## What
This is a cli tool for getting readings from the Xayav Mask Creator dongle, written in Go. This is useful if you are
taking readings using a computer that is not the computer doing primary HDMI output for a DLP/LCD 3d printer.
For example, if you use NanoDLP's mask creator gui and you just need to fill in values from a grid.

## Compatibility
It should be compatible with Windows, Linux, and MacOSX

## Why
I use NanoDLP running on a raspberry pi. I wanted to create a mask to even out the light distribution across
the LCD screen of my Wanhao Duplicator 7. There is at least one tutorial out there on creating a simple
Arduino device that will measure UV output. You then print out a locating grid, and use NanoDLP to generate
a mask using the data taken from the tool.

Xayav has created some pretty nice hardware to solve this priblem. The problem is that their software
doesn't match my use case. Their software runs on Windows only, and depends on having the printer HDMI
plugged in to the computer that is running the Xayav measurement and mask creating software. Since my
printer is plugged into a raspberry pi, this doesn't work for me. I don't want to go through the configuration
and calibration hassles required to get the D7 running on a windows computer just for mask creation purposes.
Not to mention, I don't have any Windows computers in the same building as the printer.

Considering Linux is prevalent in the maker spaces out there and being used for printer control everywhere,
I decided a linux tool was necessary. Xayav said that they could not give me an estimate on how soon a linux
version of the mask creator tool would be available, so I decided to write my own.

NOTE: I am not affiliated with Xayav in any way, and this tool is not supported by them. Don't send them
support requests about this tool. Since they have not published their protocol or software source, I can't
guarantee this will be free of problems. Even if they do publish details, If you use this software you are
taking responsibility for anything that happens and agree that I will not be held liable. It sucks that
we live in a world that I have to write that.

## How
Right now, you need to be able to compile this program in order to use it. I will be doing some refactoring to
allow for configurability without compiling, and I will release precompiled binaries for Windows (x64), Linux (X64 and Arm, for the pi)
and MacOSX.

You will need to edit main.go and set your COM port (Windows) or serial port path (Linux or MacOSX), and then
 compile and run it.

- Plug the Mask Creator dongle into your computer, and note the serial device assigned to it. On Windows for me,
this is COM5. On Linux for me, this is /dev/ttyUSB0. Your mileage will likely vary.
- Set your port in main.go in the setup() function near the bottom. For example `portName := "COM5"` for COM5.
- Compile and run!
- Display the NanoDLP alignment grid, and align your Mask Creator alignment grid to that
- Place the mask creator dongle into the alignment grid, and click and release the button. Hold the dongle in
position until a value is logged to terminal. Something like this will show
``` 2018/02/04 14:22:09 {1 196} ```. The (1 196) means 1st cell (or 1st button click), and the value is 196.
- Put the value into the correct cell in NanoDLP. Make sure you have generated the cell and grid in NanoDLP with
the correct orientation. The raspberry pi requires the D7 display to be rotated.
- Continue until done!

## Upcoming features
- Display the device serial number somehow
- Config file based configuration for serial port and grid size
- Startup routine that lists detected serial port and allows you to chose the right one
- Define grid size so the cell number correlates to a particular grid postion
- Perhaps generating and export the mask image directly
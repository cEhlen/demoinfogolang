package main

import (
	"fmt"
	"log"
)

type Dumper struct {
	demoFile    DemoFile
	frameNumber int
}

func (d *Dumper) Open(fileName string) {
	if err := d.demoFile.Open(fileName); err != nil {
		log.Println(err)
		log.Fatal("Unable to open demo file")
	}
}

func (d *Dumper) PrintHeader() {
	fmt.Println("FileInfo: ", string(d.demoFile.DemoHeader.fileInfo[:]))
	fmt.Println("Protocol Version: ", d.demoFile.DemoHeader.protocolVersion)
	fmt.Println("Net Protocol Version: ",
		d.demoFile.DemoHeader.networkProtocolVersion)
	fmt.Println("Server: ", string(d.demoFile.DemoHeader.server[:]))
	fmt.Println("Client: ", string(d.demoFile.DemoHeader.client[:]))
	fmt.Println("Map: ", string(d.demoFile.DemoHeader.mapName[:]))
	fmt.Println("Game Dir: ", string(d.demoFile.DemoHeader.gameDir[:]))
	fmt.Println("Playback time: ", d.demoFile.DemoHeader.playbackTime)
	fmt.Println("Playback ticks: ", d.demoFile.DemoHeader.playbackTicks)
	fmt.Println("Playback frames: ", d.demoFile.DemoHeader.playbackFrames)
	fmt.Println("Signon length: ", d.demoFile.DemoHeader.signonLength)
}

func (d *Dumper) Dump() {
	stop := false

	for !stop {
		fmt.Println("## READ LOOP ##")
		cmd, tick, player := d.demoFile.ReadCmdHeader()
		fmt.Printf("CMD: %v  |  TICK: %v  | PLAYER: %v\n", cmd, tick, player)
		switch cmd {
		case DEM_STOP:
			stop = true
		case DEM_CONSOLECMD:
			// read raw data
			fmt.Println("DEM_CONSOLECMD")
		case DEM_DATATABLES:
			// read some data
			fmt.Println("DEM_DATATABLES")
		case DEM_STRINGTABLES:
			// read raw data
			// dump string tables
			fmt.Println("DEM_STRINGTABLES")
		case DEM_USERCMD:
			// read user command
			fmt.Println("DEM_USERCMD")
		case DEM_SIGNON, DEM_PACKET, DEM_SYNCTICK:
			// handle packet
			fmt.Println("DEM_PACKET")
		}
	}
}

package main

import (
	"encoding/hex"
	"flag"
	"fmt"

	mpeg2ts "github.com/eds-cloud/go-mpeg2-ts"
)

var pmtPID int

func main() {

	_tsFilePath := flag.String("ts", "test.ts", "mpeg2ts file")
	flag.Parse()

	tsFilePath := *_tsFilePath

	stream, err := mpeg2ts.OpenFile(tsFilePath)
	if err != nil {
		panic(err)
	}

	var packet *mpeg2ts.Packet
	for {
		packet, err = stream.ReadPacket()
		if err != nil {
			break
		}
		fmt.Println("Packet Index", packet.Index, "PID", packet.PID)
		fmt.Printf("%s", hex.Dump(packet.Data))
		if packet.PID == 0 {
			pat, err := packet.ParsePAT()
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Print("PAT:\n")
			for _, program := range pat.Programs {
				fmt.Printf("Program_number: %d, ", program.ProgramNumber)
				if program.ProgramNumber != 0 {
					fmt.Printf("Program_map_PID: %d\n", program.ProgramMapPID)
					pmtPID = int(program.ProgramMapPID)
					break
				} else {
					fmt.Printf("network_PID: %d\n", program.NetworkPID)
				}
			}
		} else if pmtPID == int(packet.PID) {
			pmt, err := packet.ParsePMT(true)
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Print("PMT:\n")
			for _, s := range pmt.Streams {
				fmt.Printf("ElementaryPID: %d, Type: %d (0x%x)\n", s.ElementaryPID, s.Type, s.Type)
			}
		}
	}
	fmt.Println(err)
}

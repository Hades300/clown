package guard

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

var (
	device = "en2"
	handle *pcap.Handle
	expr   string
	bpfIns []pcap.BPFInstruction
	err    error
)

func init() {
	bpfIns = []pcap.BPFInstruction{}
	expr = ""
	handle, err = pcap.OpenLive(device, 1024, true, time.Second*30)
	if err != nil {
		log.Fatal("pcap捕捉流量失败", err)
	}
}

func listen() chan gopacket.Packet {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	bpfIns, err = handle.CompileBPFFilter(expr)
	if len(bpfIns) != 0 {
		err := handle.SetBPFInstructionFilter(bpfIns)
		if err != nil {
			panic(err)
		}
	}
	return packetSource.Packets()
}

func Send(data []byte) error {
	return handle.WritePacketData(data)
}

// example arp
// tcp and http

func Filter(e string) {
	expr = e
}

func Close() {
	handle.Close()
}

func ListenPacketInfo() []interface{} {
	packets := listen()
	i := 1
	for {
		select {
		case packet := <-packets:
			fmt.Println(i, "Last Layer:", getLastLayer(packet).LayerType())
			fmt.Println(i, "Last Layer:", fmt.Sprintf("%s", packet.Metadata().CaptureInfo))
			i++
		}
	}
}

// the last layer before the probably existed payload layer
func getLastLayer(packet gopacket.Packet) gopacket.Layer {
	mlayers := packet.Layers()
	var lastLayer gopacket.Layer
	lastLayer = mlayers[len(mlayers)-1]
	if lastLayer.LayerType() == gopacket.LayerTypePayload {
		lastLayer = mlayers[len(mlayers)-2]
	}
	return lastLayer
}

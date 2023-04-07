package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	pkt_id := 0
	pcap_path := "./dnp3_test_data_part2.pcap"
	if handle, err := pcap.OpenOffline(pcap_path); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		// for packet := range packetSource.Packets() {
		for packet := range packetSource.Packets() {
			if packet == nil {
				fmt.Println("oops")
			}
			// Get the TCP layer from this packet
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				fmt.Println("This is a TCP packet!")
				// Get actual TCP data from this layer
				tcp, _ := tcpLayer.(*layers.TCP)
				fmt.Printf("From src port %d to dst port %d\n", tcp.SrcPort, tcp.DstPort)
			}
			// Iterate over all layers, printing out each layer type
			for _, layer := range packet.Layers() {
				fmt.Println("PACKET LAYER:", layer.LayerType())
			}
			//
			// This is a TCP packet!
			//    From src port 1167 to dst port 20000
			// PACKET LAYER: Ethernet
			// PACKET LAYER: IPv4
			// PACKET LAYER: TCP
			// PACKET LAYER: Payload

			fmt.Println("looking at packet", pkt_id)
			pkt_id += 1
			// handlePacket(packet)  // Do something with a packet here.
		}
	}
}

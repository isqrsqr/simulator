package main

import (
	"fmt"
	"github.com/google/gopacket"
	_ "github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
    "github.com/i-sqr-s-sqr/i2s2/mrn"
)

func main() {
    var pcap_files [4]string
    pcap_files[0] = "./dnp3_test_data_part2.pcap"
    pcap_files[1] = "./dnp3_read.pcap"
    pcap_files[2] = "./dnp3_select_operate.pcap"
    pcap_files[3] = "./dnp3_write.pcap"

    for _, pcap_path := range pcap_files {
        if handle, err := pcap.OpenOffline(pcap_path); err != nil {
            panic(err)
        } else {
            packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
            // for packet := range packetSource.Packets() {
            for packet := range packetSource.Packets() {
                if packet == nil {
                    fmt.Println("oops")
                }
                mrp, _ := mrn.Eth_MultiRes_real(packet)
                if mrp != nil {
                    mrp.Print(false, true)
                    fmt.Println(" ")
                }
            }
        }
    }
}

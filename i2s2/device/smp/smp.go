// package for simulation-emulation packet
package smp

import (
    "net"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    )

type SMP struct {
    RawPacket []byte
    // MAC addresses
    SrcMAC net.Interface
    DstMAC net.Interface
    // IP address
    SrcIP net.IP
    DstIP net.IP
    // ports
    //
    SrcPort int
    DstPort int
    PayLoad []byte 
}


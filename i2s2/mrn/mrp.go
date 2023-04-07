package mrn

import ( 
    "fmt"
    _ "errors"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "net"
)

var pckt_id = 1

type MultiResState int
const ( 
    Empty MultiResState = iota
    Real
    Synth
    Symb
) 

type MultiResPacket interface {
    Resolution() MultiResState
    Id() int
    Data() []byte
    TransportPayload() []byte
    ApplicationPayload() []byte
    SrcMac() any
    DstMac() any
    SrcIP() any
    DstIP() any
    Protocol() any
    SrcPort() any
    DstPort() any
    Print(bool,bool)
}

type MultiResPacket_real struct {
    id       int
    data     []byte
    transport_payload  []byte
    application_payload []byte
    srcMAC, dstMAC net.HardwareAddr
    srcIP, dstIP net.IP
    srcPort, dstPort uint16
    protocol layers.IPProtocol
}
func (mrp *MultiResPacket_real) Print(prt_data bool, prt_payload bool) {
    fmt.Println("mrp id", mrp.id)
    if mrp.srcMAC != nil {
        fmt.Println("srcMac", mrp.srcMAC) 
    }
    if mrp.dstMAC != nil {
        fmt.Println("dstMac", mrp.dstMAC) 
    }
    if mrp.srcIP!= nil {
        fmt.Println("srcIP", mrp.srcIP) 
    }
    if mrp.dstIP != nil {
        fmt.Println("dstIP", mrp.dstIP) 
    }
    fmt.Println("srcPort", mrp.srcPort) 
    fmt.Println("dstPort", mrp.dstPort) 
    fmt.Println("protocol", mrp.protocol) 
    if prt_data {
        fmt.Println("org", mrp.data)
    }
    if prt_payload && len(mrp.transport_payload) > 0  {
        fmt.Println("transport payload", mrp.transport_payload)
    }
    if prt_payload && len(mrp.application_payload) > 0 {
        fmt.Println("application payload", mrp.application_payload)
    }
}

func (mrp *MultiResPacket_real) Resolution() MultiResState {
    return Real
}

func (mrp *MultiResPacket_real) Id() int {
    return mrp.id
}

func (mrp *MultiResPacket_real) Data() []byte {
    return mrp.data
}

func (mrp *MultiResPacket_real) ApplicationPayload() []byte {
    return mrp.application_payload
}

func (mrp *MultiResPacket_real) TransportPayload() []byte {
    return mrp.transport_payload
}

func (mrp *MultiResPacket_real) SrcMAC() any {
    return mrp.srcMAC
}

func (mrp *MultiResPacket_real) DstMAC() any {
    return mrp.DstMAC
}

func (mrp *MultiResPacket_real) SrcIP() any {
    return mrp.srcIP
}

func (mrp *MultiResPacket_real) DstIP() any {
    return mrp.DstIP
}

func (mrp *MultiResPacket_real) SrcPort() any {
    return mrp.srcPort
}

func (mrp *MultiResPacket_real) DstPort() any {
    return mrp.DstPort
}

func (mrp *MultiResPacket_real) Protocol() any {
    return mrp.protocol
}

// convert presented gopacket Pack into a Real Packet
// 
func Eth_MultiRes_real(packet gopacket.Packet) (*MultiResPacket_real, error) {

    mrp := new(MultiResPacket_real)

    // get the original packet bytes
    mrp.data = packet.Data()

    mrp.id = pckt_id 
    pckt_id += 1

    ethLayer := packet.Layer(layers.LayerTypeEthernet)
    if ethLayer != nil {
        eth, _ := ethLayer.(*layers.Ethernet)
        mrp.srcMAC = eth.SrcMAC
        mrp.dstMAC = eth.DstMAC
    }

    ipLayer := packet.Layer(layers.LayerTypeIPv4)
    if ipLayer != nil {
        ip, _ := ipLayer.(*layers.IPv4)
        mrp.srcIP = ip.SrcIP
        mrp.dstIP = ip.DstIP
        mrp.protocol = ip.Protocol
    } 

    tcpLayer := packet.Layer(layers.LayerTypeTCP)
    if tcpLayer != nil {
        tcp , _ := tcpLayer.(*layers.TCP)
        mrp.srcPort = uint16(tcp.SrcPort)
        mrp.dstPort = uint16(tcp.DstPort)
        mrp.transport_payload = tcp.Payload
    } else {
        udpLayer := packet.Layer(layers.LayerTypeUDP)
        if udpLayer != nil {
            udp , _ := udpLayer.(*layers.UDP)
            mrp.srcPort = uint16(udp.SrcPort)
            mrp.dstPort = uint16(udp.DstPort)
            mrp.transport_payload = udp.Payload
        }
    }
    appLayer := packet.ApplicationLayer()
    if appLayer != nil {
        mrp.application_payload = appLayer.Payload()
    }
    return mrp, nil 
}


package nic
import (
    _ "fmt"
    "net"
    "errors"
    "github.com/google/gopacket"
)

var media_id = 0
type MediaType int

var nic_id = 1


// a nic in i2s2 can be facing different kinds of connections.
const (
    PhysLink int = iota      // wired point-to-point link
    Broadcast                // wired broadcast link                
    DirectedWireless         // point-to-point wireless
    BroadcastWireless        // broadcast wireless
    Meta                     // aggregated connection
)

// a NIC object captures the important details of a NIC
type NIC struct {
    Id int                  // Id numbers are used to identify endpoints.  
                            // We maintain a map of endpoint ids to NIC struct pointers
    MAC net.Interface       // core golang net package representing a MAC address
    SymMAC int              // MACs that face Meta media have some symbolic representation give
                            // the index into a map that points to something describing it
    MAC_is_symbolic bool    // flag indicating whether MAC is physical (contained in MAC)
                            // or symbolic (contained in SymMAC
    IPAdrs net.IP           // using net packages IP representation, which includes IPv6
    Media MediaType

    // endpoints have integer identities that identify other NICs
    Endpoints map[int]bool
}


// remember that this NIC connects to the NIC with id ep
func (n *NIC) AddEndpoint(ep int) {
    n.Endpoints[ep] = true
}

// remember that this NIC connects to all the NIC in a set of them
func (n *NIC) AddEndpointSet(eps map[int]bool) {
    for ep, v := range eps {
        if v == false {
            n.Endpoints[ep] = true
        }
    }
}


// create a new instance of a NIC
func New(name string, macname string, ipaddr string, media_type MediaType) (*NIC, error) {

    // net.Interface.HardwareAddr holds the mac address if it exists
    intrfc := new(net.Interface)

    var symbolic_mac bool 
    mac, err := net.ParseMAC(macname)
    if err != nil {
        symbolic_mac = true
    } else {
        intrfc.HardwareAddr = mac
        symbolic_mac = false
    }
    // remember the string name 
    intrfc.Name = name
    intrfc.Index = media_id
    media_id += 1

    ip_addr := net.ParseIP(ipaddr)
    if ip_addr == nil {
        return nil, errors.New("ill-formed IP address for interface"+ipaddr)
    }
    media := media_type
    endpoints := make(map[int]bool)
    new_nic := NIC{Id:nic_id, MAC:intrfc, IPAdrs:ip_addr, 
        Media:media, Endpoints:endpoints, MAC_is_symbolic: symbolic_mac}

    nic_id += 1
    return &new_nic, nil
}



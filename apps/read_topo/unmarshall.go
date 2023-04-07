i8ba491u3na1



package main
import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "container/heap"
    "math"
)
type Interface struct {
    Bandwidth int
    Faces string 
}

type IntrfcList []Interface

type Router struct {
    Name string 
    Number int
    Intrfc IntrfcList
}

type Host struct {
    IP_Addr string
    Home string
    Name string
    Number int
}

type Network struct {
    Level string
    Name string
    Number int
}



type rtr_slice []Router
type host_slice []Host
type network_slice []Network

type topodict struct {
    Hosts host_slice `json:"Hosts"`
    Networks network_slice `json:"Networks"`
    Routers rtr_slice `json:"Routers"`
} 

type NetworkMap map[string]Network
type HostMap map[string]Host
type RouterMap map[string]Router

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func read_topo(topo_file_name string) *topodict {

    // create the template
    var IntrfcsExample IntrfcList
    IntrfcsExample = append(IntrfcsExample, Interface{Bandwidth:1000,Faces:"network1"})
    IntrfcsExample = append(IntrfcsExample, Interface{Bandwidth:5000,Faces:"network2"})
    rtr_example := Router{Name:"router", Number:4, Intrfc:IntrfcsExample}
    host_example := Host{IP_Addr:"0.0.0.0", Home:"LAN-1", Name:"server", Number: 1}
    network_example := Network{Level:"T1", Name:"backbone", Number:0}

    var rtrs rtr_slice
    rtrs = append(rtrs, rtr_example)

    var hosts host_slice
    hosts = append(hosts, host_example)

    var networks network_slice
    networks = append(networks, network_example)


    var example topodict
    example.Hosts = hosts
    example.Networks = networks
    example.Routers = rtrs

    dat, err := os.ReadFile(topo_file_name)
    check(err)

    json.Unmarshal(dat, &example)
    return &example
}


func main() {
    topo_file_name := "topo.json"
    topo_dict := read_topo(topo_file_name)
    // create Network, Host, and Router lookups

    NetworkByName := make(map[string]Network)
    HostByName := make(map[string]Host)
    RouterByName := make(map[string]Router)

    NetworkByNumber := make(map[int]Network)
    HostByNumber := make(map[int]Host)
    RouterByNumber := make(map[int]Router)

    for _, net := range((*topo_dict).Networks) {
        NetworkByName[ net.Name ] = net
        NetworkByNumber[ net.Number] = net
    }

    for _, host := range((*topo_dict).Hosts) {
        HostByName[host.Name] = host
        HostByNumber[host.Number] = host
    }

    for _, router := range((*topo_dict).Routers) {
        RouterByName[router.Name] = router
        RouterByNumber[router.Number] = router
    }

    // record connections between networks, hosts, and routers
    Edges := make(map[int][]int)

    // initialize storage for edges as viewed from networks
    for number, _ := range(NetworkByNumber) {
        Edges[number] = make([]int,0)
    }

    for host_number, host := range(HostByNumber) {

        // edges exist between a host and the network declared as
        // its home.  Represent that edge both on the host and on the network

        // one edge as viewed from the host 
        Edges[host_number] = make([]int,1)

        // Number of the network declared as host's home
        net_number := NetworkByName[host.Home].Number

        // save that edge as viewed from host
        Edges[host_number][0] = net_number

        // save that edge as viewed from home network 
        Edges[net_number] = append(Edges[net_number], host_number)
    }

    for rtr_number, router := range(RouterByNumber) {

        // the Intrfc slice gives interfaces which face networks.
        // A router's edges go to those networks.  We know how many
        // as viewed from the network
        Edges[rtr_number] = make([]int, len(router.Intrfc))

        // for every interface record as an edge the number of
        // the network faced by the interface

        // for a given router visit each interface
        for idx, intrfc := range(router.Intrfc) {

            // get the network name
            faced_network := NetworkByName[intrfc.Faces].Number 

            // save the edge as viewed from the router
            Edges[rtr_number][idx] = faced_network

            // save the edge as viewed from the network
            Edges[faced_network] = append(Edges[faced_network], rtr_number)
        }
    }

    fmt.Println("done")
}

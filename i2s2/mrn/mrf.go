package mrn

// multi-resolution flow structure
import ( 
    _ "fmt"
    _ "github.com/google/gopacket"
    _ "github.com/google/gopacket/layers"
    _ "net"
)

type FlowStage struct {
    StageId int
    Active bool
    Bandwidth float64
}

func (fs *FlowStage) SetBandwidth(bw float64) {
    fs.Bandwidth = bw
}

type MultiResFlow struct {
    Id int
    Edge int            // stage number where leading edge of multi-res flow is going
    Src, Dst int        // endpoint ids for routing
    Stages []*FlowStage
}

var mrf_counter int = 0

func NewMultiResFlow(src int, dst int, edges *[] ) (*MultiResFlow) {
    var mrf MultiResFlow

    mrf.Id = mrf_counter
    mrf_counter++
    mrf.Src = src
    mrf.Dst = dst
    mrf.Stages = make([]*FlowStage, len(path))

    for idx := 0; idx<len(path); idx++ {
        new_flow_stage := FlowStage{ StageId:path[idx], Active:false, Bandwidth:0.0}
        mrf.Stages[idx] = &new_flow_stage
        new_flow_stage.StageId = path[idx]
        new_flow_stage.Active = false
        new_flow_stage.Bandwidth = 0.0 
    }      
    return &mrf
}

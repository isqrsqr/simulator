
import (
    "errors"
	_ "strconv"
    _ "fmt"
    _ "github.com/i-sqr-s-sqr/i2s2/vrtime"
    _ "github.com/i-sqr-s-sqr/i2s2/evtq"
      "github.com/i-sqr-s-sqr/i2s2/pgn"

)

// particular ProtocolGraphNode instances should 
// instantiate the functions below
//
// type ProtocolGraphNodeInterface interface {
//    fromBelow([]*byte) 
//    toBelow([]*byte) 
//    fromAbove ([]*byte) 
//    toAbove ([]*byte) 
//    AddChild(int, func([]*byte, int))
//    AddParent(int, func([]*byte, int))
//}	

func New(string name, int id) error err {
    sn_ip, sn_err := sn.New(name, id)
    if err != nil {
        return nil, err
    }

    sn_ip.Children = create(map[int]*pgn.ProtocolGraphNode)
    sn_ip.Parents  = create(map[int]*pgn.ProtocolGraphNode)
}


// -------- Stack_node template ------
type ProtocolGraphNode struct {
	Name     string
	Id       int
    HandleMsgFromBelow func(*[]byte, int)
    HandleMsgFromAbove func(*[]byte, int)
    children map[int]bool
    parents map[int]bool
}
func New(name string, sn_id int,
    HandleMsgFromBelow func(*[]byte, int), 
    HandleMsgFromAbove func(*[]byte, int)) (*ProtocolGraphNode, error) {

        if ProtocolGraphNodeMap == nil {
                ProtocolGraphNodeMap = make(map[int]*ProtocolGraphNode)

        }
        if ProtocolGraphNodeMap[sn_id] != nil {
            return nil, errors.New("non-unique protocol graph node id")
        }

        snode := ProtocolGraphNode{Name: name, Id: sn_id, HandleMsgFromBelow: HandleMsgFromBelow,
            HandleMsgFromAbove: HandleMsgFromAbove}
        ProtocolGraphNodeMap[sn_id] = &snode
        return &snode, nil
}






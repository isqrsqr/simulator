package pgn

import (
	"strconv"
    "fmt"
    "github.com/isqrsqr/simulator/i2s2/vrtime"
    "github.com/isqrsqr/simulator/i2s2/evtm"
)

// ---------- interfaces ---------
type Child interface {
	Push_down([]byte, int)
}

type Parent interface {
	Push_up([]byte, int)
}

type Msg_xfer interface {
    Xfer([]byte) []byte
}

//--------------------------------

// -------- Protocol_node realization ------
type Protocol_node struct {
	Name     string
	Id       int
    Xfer_up   Msg_xfer
    Xfer_down Msg_xfer
	Parents  map[int]Parent
	Children map[int]Child
}

type Protocol_node_context struct {
    // what scheduling queue to use
    Evtm *evtm.EventManager

    // direction is pushing up from below or down from above
    down_direction bool

    // entry index is integer identity of source of message being 
    // passed along
    Index int
}

//   event handler for message delivered to Protocol_node
//   'context' and 'message' are from event, handed up from event scheduler
func (pn *Protocol_node) Handle(context any, message any) {

    // calls between Protocol_nodes always have Protocol_node_context structs as context
    context_struct := context.(Protocol_node_context)

    // calls between Protocol_nodes always carry byte slices
    msg :=  message.([]byte)

    if context_struct.down_direction {
        pn.Push_down(context_struct, msg, context_struct.Index) 
    } else {
        pn.Push_up(context_struct, msg, context_struct.Index) 
    }
}

// process a message from below, push up to all parents
//
func (pn *Protocol_node) Push_up(context Protocol_node_context, frame []byte, from int) {
	fmt.Println("push up at", pn.Id, "from", from)

    offset := vrtime.Time{Key1:0, Key2:0}

	if pn.Parents != nil {
		for _, parent := range pn.Parents {
            context.Index = pn.Id
            new_frame := pn.Xfer_up.Xfer(frame)

            // where does the event queue come from?
            // a context is passed in, can it be used?
            // the data should be the new_frame
            // time offset is created here

            new_evt_id, _ := context.Evtm.Schedule(context, new_frame, Handle, offset) 


			parent.Push_up(new_frame, pn.Id)
		}
	} else {
		for _, child := range pn.Children {
			child.Push_down(new_frame, pn.Id)
		}
	}
}

func (pn *Protocol_node) Push_down(frame []byte, from int) {

	passing := string(frame.Msg[frame.First:frame.Last])
	Add_on := ",d " + strconv.FormatInt(int64(pn.Id), 10)
	passing += Add_on
	new_frame := new([]byte)
	new_frame.Msg = []byte(passing)
	new_frame.First = frame.First
	new_frame.Last = frame.Last + len(Add_on)

	if pn.Children != nil {
		for _, child := range pn.Children {
			child.Push_down(new_frame, pn.Id)
		}
	} else {
		fmt.Println(new_frame.ToStr())
	}
}

func (pn *Protocol_node) Add_child(child *Protocol_node, child_id int) {
	if pn.Children == nil {
		pn.Children = make(map[int]Child, 1)
	}
	pn.Children[child_id] = child
}

func (pn *Protocol_node) Add_parent(parent *Protocol_node, parent_id int) {
	if pn.Parents == nil {
		pn.Parents = make(map[int]Parent, 1)
	}
	pn.Parents[parent_id] = parent
}

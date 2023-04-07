package main

// 'spark' code demonstrates scheduling and time advance
// in i2s2
//
// A doubly linked list of 'spark_context' structs are created.
// The code traverses the linked list, forward and then backward,
// advancing time by 0 in each call in the forward direction, advancing by 1
// in each call of the backward direction.  Each call prints the name of the 
// struct, the virtual time, and the identity of the event that scheduled the call.

// The i2s2 framework allows one to schedule a function call,
// associated with a simulation time.   The function call 
// 
//    func Spark(sys evtm.SystemContext, context any, data any) bool {
//
// is passed three
// objects.  The 'sys' structure is defined by the evtm module,
//
// type SystemContext struct {
//    EventMan     *EventManager
//    Time         vrtime.Time
//    EventId      int
//    EventHandler func(SystemContext, any, any) bool
//}
//
// and gives the context of the system for the call.  In particular it points
// to the event list manager from which the call was made, the virtual time,
// the unique integer-valued identity of the event being evaluated by the call,
// and a pointer to the calling function itself.
//
// the 'context' argument can be used to carry pointer to a structure with
// whatever the local context (e.g., here, a pointer to the spark_context struct
// associated with the call).  Finally the 'data' argument can be thought of as a pointer
// to a struct that carries a 'message' to the event handler.
// 
//
// A method 'Spark' when evaluated will print the struct's identity string
// and the time associated with the call.   The 'direction' (initially 0) determines
// what happens next.  If the direction is 0 and there is a 'back' link, then
// an evaluation of the spark_context structure pointed to by the back link is scheduled
// to occur 0 units of time later, e.g., immediately.   If the direction is 0 but there
// is no 'back' pointer the chain reverses, changing direction to 1, following 'front' pointers,
// and scheduling evaluations 1 unit of time in the future. 
//
// building and running manager.go should give the output
//
// % go run manager.go 
// hit Spark 1 at 1 with event id 1
// hit Spark 2 at 1 with event id 2
// hit Spark 3 at 1 with event id 3
// hit Spark 4 at 1 with event id 4
// hit Spark 3 at 1 with event id 5
// hit Spark 2 at 2 with event id 6
// hit Spark 1 at 3 with event id 7
// Done

import (
	"fmt"
	"github.com/isqrsqr/simulator/i2s2/evtm"
	"github.com/isqrsqr/simulator/i2s2/vrtime"
)

type spark_context struct {
	spark_string string
	front, back  *spark_context
}

// func Spark is the event handler.  Given system context, local context, and data to be acted on
//  The bool return value semantics is defined here by the user.  The code below returns true
// if the execution finished the list traverse and false otherwise, but in this example the return
// value has no impact on anything
//
func Spark(sys evtm.SystemContext, context any, data any) bool {

    // 'any' structs need to be cast.  Here we know that what was passed in 'context' was
    // a pointer to a spart_context struct
	scx := context.(*spark_context)

    // The data in these calls is the direction, and is an int
	direction := data.(int)

    // we choose here to schedule the next event using the same event manager as was used
    // to fire up this evaluation
	event_manager := sys.EventMan

    // not yet evident but time is a tuple of integers (assumed to be non-negative), a primary key
    // and a tie-breaker
	current_time := sys.Time

    // all the heavy lifting right here!
	fmt.Println("hit", scx.spark_string, "at", current_time.Key1, "with event id", sys.EventId)

    // we're going to add an 'offset' to the current time to create a time when the scheduled
    // evaluation takes place.  The offset needs to be of the i2s2 virtual time type
	offset := vrtime.Time{Key1: vrtime.Key1_type(direction), Key2: vrtime.Key2_type(0)}

    // the front to back and return calling pattern described earlier
	if direction == 0 && scx.back != nil {
		event_manager.Schedule(scx.back, 0, Spark, offset)
	} else if direction == 0 && scx.front != nil {
		event_manager.Schedule(scx.front, 1, Spark, offset)
	} else if direction == 1 && scx.front != nil {
		event_manager.Schedule(scx.front, 1, Spark, offset)
	} else {
        // when we run out of links to follow we are done
		fmt.Println("Done")
		return true
	}
	return false
}

func main() {

    // always need an event list
	evt_manager := evtm.New()

    // generate the list node structs, we'll link them up by separate code to come
	s1 := spark_context{spark_string: "Spark 1", front: nil, back: nil}
	s2 := spark_context{spark_string: "Spark 2", front: nil, back: nil}
	s3 := spark_context{spark_string: "Spark 3", front: nil, back: nil}
	s4 := spark_context{spark_string: "Spark 4", front: nil, back: nil}

    // separate code to code s1 <-> s2 <-> s3 <-> s4
	s1.back = &s2
	s2.front = &s1
	s2.back = &s3
	s3.front = &s2
	s3.back = &s4
	s4.front = &s3

	direction := 0

    // first event is scheduled for time 0+1
	offset := vrtime.Time{Key1: 1, Key2: 0}

    // Schedule arguments are local-context, data, calling function, offset ahead into virtual time
	evt_manager.Schedule(&s1, direction, Spark, offset)

    // off we go
	evt_manager.Run()
}

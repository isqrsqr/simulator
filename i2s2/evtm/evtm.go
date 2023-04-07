// package for managing the scheduling and execution of events
//   depends on the package for implementing the event list (evtq), 
//   and the package for implementing virtual time (vrtime)

package evtm

import (
	_ "fmt"
	"github.com/isqrsqr/simulator/i2s2/evtq"
	"github.com/isqrsqr/simulator/i2s2/vrtime"
	_ "strconv"
)


// Every event handling function receives as its first input parameter
// a 'SystemContext' structure.   The idea is to have this structure
// carry context and data independent information that the handler
// may need, e.g. the identity number of the event that causes the
// handler executing, the virtual time of the event, a pointer to
// the event manager structure that holds the event list from which
// the initiating event was drawn.  

type SystemContext struct {
	EventMan     *EventManager
	Time         vrtime.Time
	EventId      int
	EventHandler func(SystemContext, any, any) bool
}


// an Event structure packages up the context and data sensitive 
// information to be included when scheduling an event, and
// to be used in dispatching the event handler

type Event struct {
	Context any
	Data    any
	EventHandler func(SystemContext, any, any) bool
	EventId int
}

// An EventManager structure holds information needed
// to schedule and execute events.  It has a pointer to an
// EventQueue (defined in package evtq), and a copy of the virtual
// time of the last event to have been removed from the EventQueue (which
// is interpreted as the virtual time of that event).  It has a Boolean
// flag also which, if changed to false when processing an event, will
// inhibit the dispatch of further events until the event manager is told to run
// again

type EventManager struct {
	EventList *evtq.EventQueue
	Time      vrtime.Time
	RunFlag   bool
}

// an EventManager constructor.  It creates an empty event queue
// sets the virtual time to zero, initializes the 'running' flag to false.

func New() *EventManager {
	newEq := evtq.New()
	newEm := &EventManager{
		EventList: newEq,
		Time:      vrtime.ZeroTime(),
		RunFlag:   false,
	}
	return newEm
}

// function Run starts the event dispatch loop for an EventManager
// that has been inactive.  It will stay in this processing loop
// until either the EventQueue is exhausted of events, or
// the EventManager's RunFlag is set to false
//

func (evtm *EventManager) Run() {

    // as long as RunFlag is true the EventManager will stay in a loop
    // the next event is pulled from the EventQueue and dispatched

	evtm.RunFlag = true

    // keep working if the RunFlag is true and there are events to 
    // dispatch

	for evtm.RunFlag == true && evtm.EventList.Len() > 0 {

        // nxtEvt pulls off the package associated with the event with least
        // time-stamp and unpacks it into
        //   a) context is information the event handler may need about where and what 
        //      it is executing.  The code that schedules the event and the code that
        //      handles the event have to be using the formatting, as the representation
        //      internal to go is interface{}.
        //   b) data is information the event handler uses to execute the event, e.g.,
        //      a message or frame.   Again, the code that schedules the event and the
        //      code that executes it have to agree on the format, because Go type checking
        //      won't do that
        //   c) handler is the function to call to handle the event.  These all have the
        //      signature  func(SystemContext, context any, data any) bool
        //      The boolean return flags whether the event was dispatched without error
        //   d) Events are given unique integer id numbers when scheduled, and evt_id
        //      returns that of the event being dispatched
 
		context, data, handler, time, evt_id := nxtEvt(evtm.EventList)

        // update the EventManager's record of current virtual time, given that the minimum time
        // event was pulled off the EventQueue

		evtm.Time = time

        // make up a SystemContext based on the information unpacked from the next event
		system_context := SystemContext{EventMan: evtm, Time: time, EventHandler: handler, EventId: evt_id }

        // dispatch the event using the information carried along by the event
		handler(system_context, context, data)
	}

    // falling out of the displatch loop we know the EventManager isn't running anymore
	evtm.RunFlag = false
}

// Stop the event dispatch loop of the EventManager 
func (evtm *EventManager) Stop() {
	evtm.RunFlag = false
}

// Schedule puts information for dispatching an event into the EventManager's
// EventQue.  It bundles together the context and data information to be
// passed to the event, the function to call when dispatching the event, and
// the offset in time (from the time of scheduling) when the event handler
// will be called
//   Schedule returns the unique id of the event just scheduled (which can be used
// to remove the event from the list before dispatch), and the scheduled time of the event dispatch

func (evtm *EventManager) Schedule(context any, data any,
	handler func(SystemContext, any, any) bool, 
    offset vrtime.Time) (int, vrtime.Time) {

    // time of the last event to be pulled from the EventQueue
	current_time := evtm.Time

    // bundle together the information needed for event dispatch
	new_evt := Event{Context: context, EventHandler: handler, Data: data}

    // generalization of event time means explicit function call to add
    // the current time to the scheduling offset

	new_time := current_time.Plus(offset)

    // put the event bundle into the EventQueue with priority equal to the
    // scheduled time, and get in return the unique event id

	evt_num := evtm.EventList.Insert(&new_evt, new_time)

    // new_evt just got placed into the EventQueue but we can still get
    // at it and put in the identify of the event that carries it
	new_evt.EventId = evt_num

	return evt_num, new_time
}

// nxtEvt pulls off the minimum time event from an EventQueue and
// debundles the information it contains, returning the
// unbundled fields

func nxtEvt(queue *evtq.EventQueue) (any, any,
	func(SystemContext, any, any) bool, vrtime.Time, int) {

	item, evt_time, evt_id, _ := queue.Pop()
	evt_struct := item.(*Event)
	return evt_struct.Context, evt_struct.Data, evt_struct.EventHandler, evt_time, evt_id
}

// Given the identity of an event, remove it from the event list
// and return a flag indicating whether the event was found and removed

func (evtm *EventManager) RemoveEvent(evt_id int) bool {
	return evtm.EventList.Remove(evt_id)
}

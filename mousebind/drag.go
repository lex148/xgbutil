package mousebind

import (
	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

// Drag is the public interface that will make the appropriate connections
// to register a drag event for three functions: the begin function, the 
// step function and the end function.
func Drag(xu *xgbutil.XUtil, win xproto.Window, buttonStr string, grab bool,
	begin xgbutil.MouseDragBeginFun, step xgbutil.MouseDragFun,
	end xgbutil.MouseDragFun) {

	ButtonPressFun(
		func(xu *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
			dragBegin(xu, ev, win, begin, step, end)
		}).Connect(xu, win, buttonStr, false, grab)
}

// dragGrab is a shortcut for grabbing the pointer for a drag.
func dragGrab(xu *xgbutil.XUtil, win xproto.Window, cursor xproto.Cursor) bool {
	status, err := GrabPointer(xu, xu.Dummy(), xu.RootWin(), cursor)
	if err != nil {
		xgbutil.Logger.Printf("Mouse dragging was unsuccessful because: %v",
			err)
		return false
	}
	if !status {
		xgbutil.Logger.Println("Mouse dragging was unsuccessful because " +
			"we could not establish a pointer grab.")
		return false
	}

	xu.MouseDragSet(true)
	return true
}

// dragUngrab is a shortcut for ungrabbing the pointer for a drag.
func dragUngrab(xu *xgbutil.XUtil) {
	UngrabPointer(xu)
	xu.MouseDragSet(false)
}

// dragStart executes the "begin" function registered for the current drag.
// It also initiates the grab.
func dragBegin(xu *xgbutil.XUtil, ev xevent.ButtonPressEvent, win xproto.Window,
	begin xgbutil.MouseDragBeginFun, step xgbutil.MouseDragFun,
	end xgbutil.MouseDragFun) {

	// don't start a drag if one is already in progress
	if xu.MouseDrag() {
		return
	}

	// Run begin first. It may tell us to cancel the grab.
	// It can also tell us which cursor to use when grabbing.
	grab, cursor := begin(xu, int(ev.RootX), int(ev.RootY),
		int(ev.EventX), int(ev.EventY))

	// if we couldn't establish a grab, quit
	// Or quit if 'begin' tells us to.
	if !grab || !dragGrab(xu, win, cursor) {
		return
	}

	// we're committed. set the drag state and start the 'begin' function
	xu.MouseDragStepSet(step)
	xu.MouseDragEndSet(end)
}

// dragStep executes the "step" function registered for the current drag.
func dragStep(xu *xgbutil.XUtil, ev xevent.MotionNotifyEvent) {
	// If for whatever reason we don't have any *piece* of a grab,
	// we've gotta back out.
	if !xu.MouseDrag() || xu.MouseDragStep() == nil ||
		xu.MouseDragEnd() == nil {
		dragUngrab(xu)
		xu.MouseDragStepSet(nil)
		xu.MouseDragEndSet(nil)
		return
	}

	// now actually run the step
	xu.MouseDragStep()(xu, int(ev.RootX), int(ev.RootY),
		int(ev.EventX), int(ev.EventY))
}

// dragEnd executes the "end" function registered for the current drag.
func dragEnd(xu *xgbutil.XUtil, ev xevent.ButtonReleaseEvent) {
	if xu.MouseDragEnd() != nil {
		xu.MouseDragEnd()(xu, int(ev.RootX), int(ev.RootY),
			int(ev.EventX), int(ev.EventY))
	}

	dragUngrab(xu)
	xu.MouseDragStepSet(nil)
	xu.MouseDragEndSet(nil)
}

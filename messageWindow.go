// Copyright (c) 2014  by f4lk0r_ 
// Licensed under the GPLv3, see the License.txt file for details

package main

import(
//go ncurses library. Documentation here: https://godoc.org/code.google.com/p/goncurses
    gc "code.google.com/p/goncurses"
//    "log"
    //janimo's textsecure library. Documentation here: https://godoc.org/github.com/janimo/textsecure
//  ts "github.com/janimo/textsecure"
)

// Prints messages to the message window 
func printToMsgWindow(msg string, msgWin *gc.Window, amSending bool) {
    lines := int(len(msg) / (msgWinSize_x-1)) + 1
    y,x := globalInputWin.CursorYX()
    debugLog.Println("HIT")
    if amSending == true {
        msgWin.MovePrint((msgWinSize_y-lines),0,msg)
    } else{
        if lines > 1 {
            msgWin.MovePrint((msgWinSize_y-lines),int(msgWinSize_x * 3 / 4),msg)
            globalInputWin.Move(y,x)
        } else {
            space_buf := (msgWinSize_x) - len(msg)
            msgWin.MovePrint((msgWinSize_y-lines),space_buf,msg)
            globalInputWin.Move(y,x)
        }
    }
    msgWin.Refresh()
    globalInputWin.Refresh()
}

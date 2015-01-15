// Copyright (c) 2014  by f4lk0r_ 
// Licensed under the GPLv3, see the License.txt file for details

package main

import(
//go ncurses library. Documentation here: https://godoc.org/code.google.com/p/goncurses
    gc "code.google.com/p/goncurses"
//    "log"
    //janimo's textsecure library. Documentation here: https://godoc.org/github.com/janimo/textsecure
//  ts "github.com/janimo/textsecure"
    "strings"
)

// Prints messages to the message window 
func printToMsgWindow(msg string, msgWin *gc.Window, amSending bool) {
    lines := int(len(msg) / (msgWinSize_x-1)) + 1
    if amSending == true {
        msgWin.Scroll(lines)
        msgWin.ColorOn(2)
        msgWin.MovePrint((msgWinSize_y-lines),0,msg)
    } else{
        if strings.ContainsAny(msg,"\n") {
            printByLineBreakdown := strings.Split(msg,"\n")
            for i,val := range printByLineBreakdown {
                if i!= 0 {
                    msgWin.Scroll(1)
                }
                lines2 := int(len(val) / (msgWinSize_x-1)) + 1
                if lines2 > 1 {
                    msgWin.Scroll(lines2)
                    msgWin.ColorOn(1)
                    msgWin.MovePrint((msgWinSize_y-lines2),int(msgWinSize_x * 3 / 4),val)
                } else {
                    msgWin.Scroll(lines2)
                    msgWin.ColorOn(1)
                    space_buf := (msgWinSize_x) - len(val)
                    msgWin.MovePrint((msgWinSize_y-lines),space_buf,val)
                    msgWin.Scroll(-1)
                }
            }
        } else {
            if lines > 1 {
                msgWin.Scroll(lines)
                msgWin.ColorOn(1)
                msgWin.MovePrint((msgWinSize_y-lines),int(msgWinSize_x * 3 / 4),msg)
            } else {
                msgWin.Scroll(lines)
                msgWin.ColorOn(1)
                space_buf := (msgWinSize_x) - len(msg)
                msgWin.MovePrint((msgWinSize_y-lines),space_buf,msg)
                msgWin.Scroll(-1)
            }
        }
    }
    msgWin.Refresh()
    globalInputWin.Refresh()
}

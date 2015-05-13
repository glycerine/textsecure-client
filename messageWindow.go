// Copyright (c) 2014  by f4lk0r_ 
// Licensed under the GPLv3, see the License.txt file for details

package main

import(
//go ncurses library. Documentation here: https://godoc.org/code.google.com/p/goncurses
    gc "github.com/rthornton128/goncurses"
//    "log"
    //janimo's textsecure library. Documentation here: https://godoc.org/github.com/janimo/textsecure
//  ts "github.com/janimo/textsecure"
    "strings"
    "time"
    "bytes"
    "strconv"
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


//Handles keeping track of all the data structures for changing a contact.
func changeContact(contactsMenuWin * gc.Window, contactMenu * gc.Menu) {
    globalMsgWin.Erase()
    currentContact = getTel(contactMenu.Current(nil).Name())
    rows := getConversation(currentContact)
    defer rows.Close()
    var msg string
    var src string
    var t time.Time
    var b bytes.Buffer
    for rows.Next() {
        rows.Scan(&msg,&src, &t)
        if src != "You" {
            if time.Now().AddDate(0,0,-1).Before(t) {
                if t.Hour() < 10 {
                    b.WriteString("0")
                }
                b.WriteString(strconv.Itoa(t.Hour()))
                b.WriteString(":")
                if t.Minute() < 10 {
                    b.WriteString("0")
                }
                b.WriteString(strconv.Itoa(t.Minute()))
                b.WriteString(": ")
                b.WriteString(msg)
                printToMsgWindow(b.String(),globalMsgWin,false)
                b.Reset()
            } else if (time.Now()).AddDate(0,0,-6).Before(t) {
                b.WriteString(string([]byte(t.Weekday().String())[0:3]))
                b.WriteString(" at ")
                if t.Hour() < 10 {
                    b.WriteString("0")
                }
                b.WriteString(strconv.Itoa(t.Hour()))
                b.WriteString(":")
                if t.Minute() < 10 {
                    b.WriteString("0")
                }
                b.WriteString(strconv.Itoa(t.Minute()))
                b.WriteString(": ")
                b.WriteString(msg)
                printToMsgWindow(b.String(),globalMsgWin,false)
                b.Reset()
            } else {
                b.WriteString(t.Local().String())
                b.WriteString(": ")
                b.WriteString(msg)
                printToMsgWindow(b.String(),globalMsgWin,false)
                b.Reset()
            }
        } else {
            if time.Now().AddDate(0,0,-1).Before(t) {
                if t.Hour() < 10 {
                    b.WriteString("0")
                }
                b.WriteString(strconv.Itoa(t.Hour()))
                b.WriteString(":")
                if t.Minute() < 10 {
                    b.WriteString("0")
                }
                b.WriteString(strconv.Itoa(t.Minute()))
                b.WriteString(": ")
                b.WriteString(msg)
                printToMsgWindow(b.String(),globalMsgWin,true)
                b.Reset()
            } else if (time.Now()).AddDate(0,0,-6).Before(t) {
                b.WriteString(string([]byte(t.Weekday().String())[0:3]))
                b.WriteString(" at ")
                if t.Hour() < 10 {
                    b.WriteString("0")
                }
                b.WriteString(strconv.Itoa(t.Hour()))
                b.WriteString(":")
                if t.Minute() < 10 {
                    b.WriteString("0")
                }
                b.WriteString(strconv.Itoa(t.Minute()))
                b.WriteString(": ")
                b.WriteString(msg)
                printToMsgWindow(b.String(),globalMsgWin,true)
                b.Reset()
            } else {
                b.WriteString(t.Local().String())
                b.WriteString(": ")
                b.WriteString(msg)
                printToMsgWindow(b.String(),globalMsgWin,true)
                b.Reset()
            }
        }
    }
    globalMsgWin.Refresh()
}



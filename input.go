// Copyright (c) 2014  by f4lk0r_ 
// Licensed under the GPLv3, see the License.txt file for details

package main

import(
//go ncurses library. Documentation here: https://godoc.org/code.google.com/p/goncurses
    gc "code.google.com/p/goncurses"
)

// handles keyboard input
func inputHandler( inputWin *gc.Window, stdscr *gc.Window, contactsMenuWin *gc.Window, contactMenu *gc.Menu) {
    var c gc.Char
    var rawInput gc.Key
    max_y, max_x := inputWin.MaxYX()
    for {
        rawInput = inputWin.GetChar()
        c = gc.Char(rawInput)
        //debugLog.Println(rawInput)
        //debugLog.Println(c)

        //Escape to Quit
        if c == gc.Char(27) {
            break
        } else if rawInput == gc.KEY_BACKSPACE {
            //Delete Key
            y,x := inputWin.CursorYX()
            if x != 0 {
                inputWin.MoveDelChar(y,x-1)
                inputBuffer = inputBuffer[0:len(inputBuffer)-1]
                //debugLog.Println(inputBuffer)
            } else {
                if y!=0 {
                    inputWin.MoveDelChar(y-1,max_x)
                    inputBuffer = inputBuffer[0:len(inputBuffer)-1]
                    //debugLog.Println(inputBuffer)
                } else {
                    continue
                }
            }
        } else if c == gc.KEY_LEFT {
            y,x := inputWin.CursorYX()
            if x != 0 {
                inputWin.Move(y,x-1)
            }
        } else if c == gc.KEY_RIGHT {
            y,x := inputWin.CursorYX()
            if x != max_x {
                inputWin.Move(y,x+1)
            }
        } else if c == gc.KEY_UP {
            y,x := inputWin.CursorYX()
            if y != 0 {
                inputWin.Move(y-1,x)
            }
        } else if c == gc.KEY_DOWN {
            y,x := inputWin.CursorYX()
            if y != max_y {
                inputWin.Move(y+1,x)
            }
        } else if rawInput == gc.KEY_TAB {
            y,x := inputWin.CursorYX()
            gc.Cursor(0)
            escapeHandler := contactsWindowNavigation(contactsMenuWin, contactMenu)
            if escapeHandler == 1 {
                return
            }
            gc.Cursor(1)
            inputWin.Move(y,x)
        } else if rawInput == gc.KEY_RETURN {
            clearScrSendMsg(inputWin)
        } else if rawInput == gc.KEY_SRIGHT {
            inputWin.Print("\n")
            inputBuffer = append(inputBuffer,byte(10))
        } else if rawInput == gc.KEY_SLEFT {
        } else {
            inputWin.Print(string(c))
            inputBuffer = append(inputBuffer,byte(c))
            //debugLog.Println(inputBuffer)
        }
    }
}

// contact Menu Window Navigation input handler
func contactsWindowNavigation(contactsMenuWin * gc.Window, contactMenu * gc.Menu) int {
    var c gc.Char
    var rawInput gc.Key
    for {
        gc.Update()
        rawInput = contactsMenuWin.GetChar()
        c = gc.Char(rawInput)
        if rawInput == gc.KEY_TAB {
            return 0
        } else if  c == gc.Char(27) {
            return 1
        } else if rawInput == gc.KEY_RETURN {
            currentContact = getTel(contactMenu.Current(nil).Name())
            return 0
        } else if c == gc.Char('g') {
            contactMenu.Driver(gc.REQ_FIRST)
        } else if c == gc.Char('G') {
            contactMenu.Driver(gc.REQ_LAST)
        } else {
            contactMenu.Driver(gc.DriverActions[rawInput])
        }
    }
}

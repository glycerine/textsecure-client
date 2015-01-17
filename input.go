// Copyright (c) 2014  by f4lk0r_ 
// Licensed under the GPLv3, see the License.txt file for details

package main

import(
//go ncurses library. Documentation here: https://godoc.org/code.google.com/p/goncurses
    gc "code.google.com/p/goncurses"
)

type newLine struct {
    _cursorX int
    _placer int
}

// handles keyboard input
func inputHandler( inputWin *gc.Window, stdscr *gc.Window, contactsMenuWin *gc.Window, contactMenu *gc.Menu) {
    var placer = -1
    var NLlocate = map[int]newLine {
    }
    var c gc.Char
    var rawInput gc.Key
    max_y, max_x := inputWin.MaxYX()
    for {
        rawInput = inputWin.GetChar()
        c = gc.Char(rawInput)
        // debugLog.Println(rawInput)
        // debugLog.Println(c)

        //Escape to Quit
        if c == gc.Char(27) {
            break
        } else if rawInput == gc.KEY_BACKSPACE || c == gc.Char(127) {
            //Delete Key
            y,x := inputWin.CursorYX()
            var del = byte('F')
            if x != 0 {
                inputWin.MoveDelChar(y,x-1)
                del = inputBuffer[placer]
                copy(inputBuffer[placer: len(inputBuffer) - 1], inputBuffer[placer + 1:])
                inputBuffer = inputBuffer[0:len(inputBuffer)-1]
                placer--;
                if del != byte('\n') && NLlocate[y]._cursorX > x {
                    temp := newLine{NLlocate[y]._cursorX - 1, NLlocate[y]._placer - 1}
                    NLlocate[y] = temp
                }
                //debugLog.Println(inputBuffer)
            } else if y!=0 {
                    inputWin.Move(y - 1, max_x - 1)
                    inputWin.MoveDelChar(y - 1,max_x - 1)
                    del = inputBuffer[placer]
                    copy(inputBuffer[placer : len(inputBuffer) - 1], inputBuffer[placer + 1:])
                    inputBuffer = inputBuffer[0:len(inputBuffer)-1]
                    //debugLog.Println(inputBuffer)
                    placer--;
                }             
            if del == byte('\n') {
                inputWin.Move(y - 1, NLlocate[y - 1]._cursorX)
                delete (NLlocate, y - 1);
            }
        } else if c == gc.KEY_LEFT {
            y,x := inputWin.CursorYX()
            if x != 0 {
                inputWin.Move(y,x-1)
                placer--
            } else if y != 0 {
                inputWin.Move(y - 1, max_x - 1)
                placer--
            }
            if len(inputBuffer) > 0 && inputBuffer[placer + 1] == byte('\n') {
                inputWin.Move(y - 1, NLlocate[y - 1]._cursorX)
            }
        } else if c == gc.KEY_RIGHT {
            y,x := inputWin.CursorYX()
            placer++
            if inputBuffer == nil || placer == len(inputBuffer) {
                inputBuffer = append(inputBuffer,byte(' '))
            }
            if inputBuffer[placer] == byte('\n') {
                inputWin.Move(y + 1, 0)
            } else if x != max_x {
                inputWin.Move(y,x+1)
            } else {
                inputWin.Move(y + 1, x)
            }
        } else if c == gc.KEY_UP {
            y,x := inputWin.CursorYX()
            if y != 0{
                inputWin.Move(y-1,x)
                placer -= max_x
            }
            if NLlocate[y - 1]._placer != 0 {
                if NLlocate[y-1]._cursorX < x {
                    placer = NLlocate[y-1]._placer
                    inputWin.Move(y - 1, NLlocate[y - 1]._cursorX)
                } else {
                    placer = NLlocate[y-1]._placer - (NLlocate[y-1]._cursorX - x)
                }
            }
        } else if c == gc.KEY_DOWN {
            y,x := inputWin.CursorYX()
            if y != max_y {
                inputWin.Move(y+1,x)
                if NLlocate[y]._placer == 0 {
                    placer += max_x
                } else {
                    placer = NLlocate[y]._placer + x + 1
                }
            }
            if placer >= len(inputBuffer) {
                for i:= len(inputBuffer); i < placer + 1; i++ {
                    inputBuffer = append(inputBuffer, byte(' '))
                }
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
            placer = 0;
            for i := range NLlocate {
                delete (NLlocate, i);
            }
            clearScrSendMsg(inputWin)
        } else if rawInput == gc.KEY_SRIGHT {
            y,x := inputWin.CursorYX()
            inputWin.Print("\n")
            temp := newLine{x, placer}
            NLlocate[y] = temp
            placer++
            inputBuffer = append(inputBuffer,byte('\n'))
        } else if rawInput == gc.KEY_SLEFT {
        } else {
            y,x := inputWin.CursorYX()
            if inputBuffer == nil || placer == len(inputBuffer) - 1 {
                inputWin.Print(string(c))
                inputBuffer = append(inputBuffer,byte(c))
            } else {
                inputWin.Erase()
                inputBuffer = append(inputBuffer,byte(c))
                copy(inputBuffer[placer + 1:], inputBuffer[placer:])
                inputBuffer[placer + 1] = byte(c)
                inputWin.Print(string(inputBuffer))
            }
            if NLlocate[y]._cursorX > x {
                temp := newLine{NLlocate[y]._cursorX + 1, NLlocate[y]._placer + 1}
                NLlocate[y] = temp
            }
            placer++
            inputWin.Move(y,x + 1)

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
            currentContact = getTel(contactMenu.Current(nil).Name())
            return 0
        } else if  c == gc.Char(27) {
            return 1
        } else if rawInput == gc.KEY_RETURN {
            currentContact = getTel(contactMenu.Current(nil).Name())
            return 0
        } else if c == gc.Char('j') {
            contactMenu.Driver(gc.REQ_DOWN)
        } else if c == gc.Char('k') {
            contactMenu.Driver(gc.REQ_UP)
        } else if c == gc.Char('g') {
            contactMenu.Driver(gc.REQ_FIRST)
        } else if c == gc.Char('G') {
            contactMenu.Driver(gc.REQ_LAST)
        } else {
            contactMenu.Driver(gc.DriverActions[rawInput])
        }
    }
}

// Copyright (c) 2014  by f4lk0r_ 
// Licensed under the GPLv3, see the License.txt file for details

package main

import(
//go ncurses library. Documentation here: https://godoc.org/code.google.com/p/goncurses
    gc "github.com/rthornton128/goncurses"
)

var scroll = 0

type newLine struct {
    _cursorX int
    _placer int
}

// handles keyboard input
func inputHandler( inputWin *gc.Window, stdscr *gc.Window, contactsMenuWin *gc.Window, contactMenu *gc.Menu, msgWin *gc.Window) {
    var NLlocate = map[int]newLine {}
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
                if del != byte('\n') && NLlocate[y+scroll]._cursorX > x {
                    temp := newLine{NLlocate[y+scroll]._cursorX - 1, NLlocate[y+scroll]._placer - 1}
                    NLlocate[y+scroll] = temp
                }
                //debugLog.Println(inputBuffer)
            } else if y!=0 {//when x==0 and y!=0
                    inputWin.Move(y - 1, max_x - 1)
                    inputWin.MoveDelChar(y - 1,max_x - 1)
                    del = inputBuffer[placer]
                    copy(inputBuffer[placer : len(inputBuffer) - 1], inputBuffer[placer + 1:])
                    inputBuffer = inputBuffer[0:len(inputBuffer)-1]
                    //debugLog.Println(inputBuffer)
                    placer--;
                }
            if del == byte('\n') {
                inputWin.Erase()
                inputWin.Print(string(inputBuffer))
                inputWin.Move(y - 1, NLlocate[y - 1 + scroll]._cursorX)
                temp, check := NLlocate[y + scroll];
                var temp_cursor = temp._cursorX
                var temp_placer = temp._placer
                if check && NLlocate[y - 1 + scroll]._cursorX + temp_cursor >= max_x {
                    _newLine := newLine{NLlocate[y - 1 + scroll]._cursorX + temp_cursor - max_x , NLlocate[y + scroll]._placer - 1}
                    NLlocate[y + scroll] = _newLine
                    delete (NLlocate, y - 1)
                } else if  check {                                  // check if there are any '\n' this line
                    var largest = -1                                // if yes, select all '\n' and move
                    for i := range NLlocate {                       // placer by 1 and adjust cursor
                        if i >= y+scroll {                          // accordingly
                            if next_nl,ok := NLlocate[i + 1]; ok {
                                new_nl := newLine{next_nl._cursorX, next_nl._placer - 1}
                                NLlocate[i] = new_nl
                            }
                        }
                        if i > largest {
                            largest = i
                        }
                    }
                    delete (NLlocate, largest)                      // delete last map entry
                    _newLine := newLine{NLlocate[y - 1+scroll]._cursorX + temp_cursor , NLlocate[y - 1+scroll]._placer + temp_placer - 1}
                    NLlocate[y - 1+scroll] = _newLine
                } else {
                    delete (NLlocate, y - 1+scroll)
                }
            }
        } else if c == gc.KEY_PAGEDOWN {
            //debugLog.Println("HIT DOWN")
            msgWin.Scroll(-10)
            msgWin.Refresh()
            inputWin.Refresh()
        } else if c == gc.KEY_PAGEUP {
            //debugLog.Println("HIT UP")
            msgWin.Scroll(10)
            msgWin.Refresh()
            inputWin.Refresh()
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
                inputWin.Move(y - 1, NLlocate[y - 1+scroll]._cursorX)
            }
        } else if c == gc.KEY_RIGHT {
            y,x := inputWin.CursorYX()
            placer++
            if inputBuffer == nil || placer == len(inputBuffer) {
                inputBuffer = append(inputBuffer,byte(' '))
            }
            if inputBuffer[placer] == byte('\n') || x >= max_x - 1 {
                inputWin.Move(y + 1, 0)
            } else {
                inputWin.Move(y,x+1)
            }
        } else if c == gc.KEY_UP {
            y,x := inputWin.CursorYX()
            if y == 0 && placer == 0 {
                continue
            } else if y==0 && scroll > 0 {
                inputWin.Move(0,x)
                inputWin.Scroll(-1)
                scroll -= 1
                if(NLlocate[y-2+scroll]._placer != 0) {
                    inputWin.Erase()
                    inputWin.Print(string(inputBuffer[(NLlocate[y-2+scroll]._placer):]))
                } else if (placer-max_x-x >0) {
                    inputWin.Erase()
                    inputWin.Print(string(inputBuffer[(placer-x-max_x):]))
                } else {
                    inputWin.Erase()
                    inputWin.Print(string(inputBuffer))
                }
            }
            if y != 0{
                inputWin.Move(y-1,x)
                placer -= max_x
                if placer < 0 {
                    placer =0
                }
            }
            if NLlocate[y - 1 +scroll]._placer != 0 {
                if NLlocate[y-1+scroll]._cursorX < x {
                    placer = NLlocate[y-1+scroll]._placer
                    inputWin.Move(y - 1, NLlocate[y - 1+scroll]._cursorX)
                } else {
                    placer = NLlocate[y-1+scroll]._placer - (NLlocate[y - 1 + scroll]._cursorX - x)
                }
            }
        } else if c == gc.KEY_DOWN {
            y,x := inputWin.CursorYX()
            if y != max_y {
                inputWin.Move(y+1,x)
                if NLlocate[y+scroll]._placer == 0 {
                    placer += max_x
                } else {
                    placer = NLlocate[y+scroll]._placer + x + 1
                }
            } else if y == max_y {
                inputWin.Scroll(1)
                scroll += 1
                inputWin.Move(max_y-1,x)
                if NLlocate[y+scroll]._placer == 0 {
                    placer += max_x
                } else {
                    placer = NLlocate[y+scroll]._placer + x + 1
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
            placer = -1;
            for i := range NLlocate {
                delete (NLlocate, i);
            }
            sendMsg(inputWin, globalMsgWin)
        } else if rawInput == gc.KEY_SRIGHT {
            y,x := inputWin.CursorYX()
            if inputBuffer == nil || placer == len(inputBuffer) - 1 {
                if y == max_y-1 {
                    scroll++
                }
                inputWin.Print("\n")
                temp := newLine{x, placer}
                NLlocate[y+scroll] = temp
                placer++
                inputBuffer = append(inputBuffer,byte('\n'))
            } else {
                inputWin.Erase()
                inputBuffer = append(inputBuffer,byte('\n'))
                copy(inputBuffer[placer + 1:], inputBuffer[placer:])
                inputBuffer[placer + 1] = byte('\n')
                inputWin.Print(string(inputBuffer))
                temp := newLine {x, placer}
                placer ++
                nextholder, check := NLlocate[y + 1 + scroll]
                if check {
                    for i:= range NLlocate {
                        if i == y + scroll {
                            _newLine := newLine{NLlocate[i]._cursorX + 1 - x, NLlocate[i]._placer + 1}
                            nextholder := NLlocate[i + 1]
                            _ = nextholder
                            NLlocate[i + 1] = _newLine
                        } else if i > y {
                            temp := NLlocate[i + 1]
                            NLlocate[i + 1] = nextholder
                            nextholder = temp
                        }
                    }
                }
                NLlocate[y+scroll] = temp
                inputWin.Move(y + 1, 0)
            }
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
                for i := range NLlocate {
                    if i > y+scroll {
                        tempLine := newLine{NLlocate[i]._cursorX, NLlocate[i]._placer + 1}
                        NLlocate[i] = tempLine
                    }
                }
                if NLlocate[y+scroll]._cursorX >= x {
                    if NLlocate[y+scroll]._cursorX == max_x {
                        copy(inputBuffer[NLlocate[y+scroll]._placer: len(inputBuffer) - 1], inputBuffer[NLlocate[y+scroll]._placer + 1:])
                        inputBuffer = inputBuffer[0:len(inputBuffer)-1]
                        delete (NLlocate, y)
                    } else {
                        temp := newLine{NLlocate[y+scroll]._cursorX + 1, NLlocate[y+scroll]._placer + 1}
                        NLlocate[y+scroll] = temp
                    }
                }
            }
            placer++
            inputWin.Move(y,x + 1)
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
        if rawInput == gc.KEY_TAB || rawInput == gc.KEY_RETURN {
            return 0
        } else if  c == gc.Char(27) {
            return 1
        } else if c == gc.Char('j') || rawInput == gc.KEY_DOWN {
            contactMenu.Driver(gc.REQ_DOWN)
            changeContact(contactsMenuWin,contactMenu)
        } else if c == gc.Char('k') || rawInput == gc.KEY_UP {
            contactMenu.Driver(gc.REQ_UP)
            changeContact(contactsMenuWin,contactMenu)
        } else if c == gc.Char('g') {
            contactMenu.Driver(gc.REQ_FIRST)
            changeContact(contactsMenuWin,contactMenu)
        } else if c == gc.Char('G') {
            contactMenu.Driver(gc.REQ_LAST)
            changeContact(contactsMenuWin,contactMenu)
        } else {
            continue
        }
    }
}

// Gets password from input
func getPass() string {
    gc.Echo(false)
    stdscr := gc.StdScr()
    var returnString []byte

    x := 0
    var c gc.Char
    var rawInput gc.Key
    for {
        rawInput = stdscr.GetChar()
        c = gc.Char(rawInput)
        if rawInput == gc.KEY_BACKSPACE || c == gc.Char(127) {
            if x != 0 {
                returnString = returnString[0:len(returnString)-1]
                x--
            } else {
                continue
            }
        } else if rawInput == gc.KEY_RETURN {
            if x !=0 {
                return string(returnString)
            }
        } else if c > 31 && c < 127 {
            returnString = append(returnString,byte(c))
            x++
        } else {
            continue
        }
    }
}


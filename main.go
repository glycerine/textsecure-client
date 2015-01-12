package main

import (
	"fmt"
    //janimo's textsecure library. Documentation here: https://godoc.org/github.com/janimo/textsecure
	ts "github.com/janimo/textsecure"
    "os"
    "log"
    //go ncurses library. Documentation here: https://godoc.org/code.google.com/p/goncurses
    gc "code.google.com/p/goncurses"
)

//var telToName map[string]string
var debugLog = startDebugLogger()
var inputBuffer []byte

// This is here for conveniences sake when testing
// The default, built-in logging package in Go does not have the ability to log out to a file.
// I hate this crap language with the fiery passion of a thousand suns
func startDebugLogger() *log.Logger {
    file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        fmt.Println("Failed to open log file", ":", err)
        panic(err)
        return nil
    }

    MyFile := log.New(file,
        "PREFIX: ",
        log.Ldate|log.Ltime|log.Lshortfile)
    return MyFile
}

//sets up the initial configuration of curses. Keeps code in main cleaner.
func configCurses(stdscr *gc.Window) {
    if !gc.HasColors() {
        log.Fatal("Example requires a colour capable terminal")
    }
    if err := gc.StartColor(); err != nil {
        log.Fatal("Starting Colors failed:",err)
    }
    gc.Echo(false)

    if err := gc.InitPair(1, gc.C_RED, gc.C_BLACK); err != nil {
        log.Fatal("InitPair failed: ", err)
    }
    gc.InitPair(2, gc.C_BLUE, gc.C_BLACK)
    gc.InitPair(3, gc.C_GREEN, gc.C_BLACK)
    gc.InitPair(4, gc.C_YELLOW, gc.C_BLACK)
    gc.InitPair(5, gc.C_CYAN, gc.C_BLACK)
    gc.InitPair(6, gc.C_MAGENTA, gc.C_WHITE)
    gc.InitPair(7, gc.C_MAGENTA, gc.C_BLACK)

    //set background color to black
    stdscr.SetBackground(gc.Char(' ') | gc.ColorPair(0))

    stdscr.Keypad(true)
    gc.Cursor(1)
    gc.CBreak(true)
    stdscr.Clear()
    stdscr.ColorOn(1)
}

//creates the three Main big windows that make up the GUI
func createMainWindows(stdscr *gc.Window ) (*gc.Window, *gc.Window, *gc.Window, *gc.Window) {
    rows, cols := stdscr.MaxYX()
    height, width := rows,int(float64(cols) * .2)

    var contactsWin,messageWin,inputBorderWin, inputWin *gc.Window
    var err error
    contactsWin, err = gc.NewWindow(height, width, 0, 0)
    if err != nil {
        log.Fatal("Failed to create Contact Window:",err)
    }
    err = contactsWin.Border('|','|','-','-','+','+','+','+')
    if err != nil {
        log.Fatal("Failed to create Border of Contact Window:",err)
    }

    begin_x := width+1
    height = int(float64(rows) * .8)
    width =  int(float64(cols) * .8)
    messageWin, err = gc.NewWindow(height, width, 0, begin_x)
    if err != nil {
        log.Fatal("Failed to create Message Window:",err)
    }
    err = messageWin.Border('|','|','-','-','+','+','+','+')
    if err != nil {
        log.Fatal("Failed to create Border of Message Window:",err)
    }

    begin_y := int(float64(rows) * .8 ) + 1
    height = int(float64(rows) * .2)
    inputBorderWin, err = gc.NewWindow(height, width, begin_y, begin_x)
    if err != nil {
        log.Fatal("Failed to create InputBorder Window:",err)
    }
    err = inputBorderWin.Border('|','|','-','-','+','+','+','+')
    if err != nil {
        log.Fatal("Failed to create Border of the InputBorder Window:",err)
    }
    // inputBorderWin is just the border. InputWin is the actual input window
    // doing this simplifies handling text a fair amount.
    // also the derived function never errors. Which seems dangerous 
    inputWin, err = gc.NewWindow((height -2),(width-2),(begin_y+1),(begin_x+1))
    if err != nil {
        log.Fatal("Failed to create the inputWin Window:",err)
    }

    inputWin.Keypad(true)

    return contactsWin,messageWin,inputBorderWin,inputWin
}

func messageHandler(msg *ts.Message) {
    err := ts.SendMessage(msg.Source(), msg.Message())
    if err != nil {
        log.Println(err)
    }
    return
}

//creates a curses based TUI for the textsecure library
func main() {
    stdscr, err := gc.Init()
    client := &ts.Client{
        RootDir:        ".",
        ReadLine:       ts.ConsoleReadLine,
        MessageHandler: messageHandler,
    }
    ts.Setup(client)
    if err != nil {
        log.Fatal("Error initializing curses:", err)
    }
    defer gc.End()
    configCurses(stdscr)

    //Hello dialog. So the nooblets know the controls
    stdscr.Refresh()
    stdscr.MovePrintln(5,5, "Controls:")
    stdscr.Println("Escape: Exits the program.")
    stdscr.Println("Tab: Switches between the input window and the Message window.")
    stdscr.Println("Return: Sends a message")
    stdscr.Println(`Shift + Right/Left: puts in a new line '\n' character (like shift+return in facebook chat)`)
    stdscr.Println("Press any key to continue...")
    stdscr.GetChar()
    stdscr.Clear()
    stdscr.Refresh()

    contactsWin, messageWin, inputBorderWin, inputWin := createMainWindows(stdscr)
    messageWin.Refresh()
    inputBorderWin.Refresh()
    contactsWin.Refresh()
    inputWin.Move(0,0)
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
                inputWin.MoveDelChar(y-1,max_x)
                inputBuffer = inputBuffer[0:len(inputBuffer)-1]
                //debugLog.Println(inputBuffer)
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
        } else if rawInput == gc.KEY_RETURN {
            //SEND THE MESSAGE HERE WE GO BOYS AND GIRLS!
            if len(inputBuffer) != 0 {
                msg := string(inputBuffer)
                to := string("+12345678910")
                debugLog.Println(msg)
                debugLog.Println(to)
                err = ts.SendMessage(to,msg)
                if err != nil {
                    gc.End()
                    log.Fatal("SendMessage failed yo: ",err)
                }
            }
        } else if rawInput == gc.KEY_SRIGHT {
            inputWin.Print("\n")
            inputBuffer = append(inputBuffer,byte(10))
        } else {
            inputWin.Print(string(c))
            inputBuffer = append(inputBuffer,byte(c))
            //debugLog.Println(inputBuffer)
        }
    }
}

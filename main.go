// Copyright (c) 2014  by f4lk0r_ 
// Licensed under the GPLv3, see the License.txt file for details

package main

import (
    //janimo's textsecure library. Documentation here: https://godoc.org/github.com/janimo/textsecure
	ts "github.com/janimo/textsecure"
    //go ncurses library. Documentation here: https://godoc.org/code.google.com/p/goncurses
    gc "code.google.com/p/goncurses"
    "log"
//    "reflect"
)


//Global variables that make things way more convenient, rather than passing copies all the time
var(
//    contactsWin *gc.Window
//    messageWinBorder *gc.Window
      globalMsgWin *gc.Window
//    inputBorderWin *gc.Window
      globalInputWin *gc.Window
//    menu_items []*gc.MenuItem
//    contactMenu *gc.Menu
//    contactsMenuWin *gc.Window
    debugLog = startDebugLogger()
    inputBuffer []byte
    currentContact string
    msgWinSize_x int
    msgWinSize_y int
)

// getName returns the local contact name corresponding to a phone number,
// or failing to find a contact the phone number itself
func getName(tel string) string {
    if n, ok := telToName[tel]; ok {
        return n
    }
    return tel
}

// getTel returns the local contact telephone number corresponding to a name,
// or failing to find a contact the name itself
func getTel(name string) string {
    if t,ok := nameToTel[name]; ok {
        return t
    }
    return name
}

var telToName map[string]string
var nameToTel map[string]string

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
}

//creates the three Main big windows that make up the GUI
func createMainWindows(stdscr *gc.Window) (*gc.Window, *gc.Window, *gc.Window, *gc.Window, *gc.Window) {
    rows, cols := stdscr.MaxYX()
    height, width := rows,int(float64(cols) * .2)

    var contactsWin,messageWinBorder,inputBorderWin, inputWin *gc.Window
    var err error
    contactsWin, err = gc.NewWindow(height, width, 0, 0)
    if err != nil {
        log.Fatal("Failed to create Contact Window:",err)
    }
    err = contactsWin.Border('|','|','-','-','+','+','+','+')
    if err != nil {
        log.Fatal("Failed to create Border of Contact Window:",err)
    }

    // messageWinBorder is just the border. msgWin is the actual input window
    // doing this simplifies handling text a fair amount.
    // also the derived function never errors. Which seems dangerous 
    begin_x := width+1
    height = int(float64(rows) * .8)
    width =  int(float64(cols) * .8)
    messageWinBorder, err = gc.NewWindow(height, width, 0, begin_x)
    if err != nil {
        log.Fatal("Failed to create Message Border Window:",err)
    }
    err = messageWinBorder.Border('|','|','-','-','+','+','+','+')
    if err != nil {
        log.Fatal("Failed to create Border of Message Window:",err)
    }
    msgWin, err := gc.NewWindow((height -2),(width-2),1,(begin_x+1))
    if err != nil {
        log.Fatal("Failed to create the message Window:",err)
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
    msgWin.ScrollOk(true)

    return contactsWin,messageWinBorder,msgWin,inputBorderWin,inputWin
}


// In addition to sending a message using janimo's library, also clears screen and resets buffer
func sendMsg (inputWin *gc.Window, msgWin *gc.Window) {
    if len(inputBuffer) != 0 {
        msg := string(inputBuffer)
        to := currentContact
        err := ts.SendMessage(to,msg)
        printToMsgWindow(msg,msgWin, true)
        if err != nil {
            gc.End()
            log.Fatal("SendMessage failed yo: ",err)
        }
        inputBuffer = []byte{}
        inputWin.Erase()
    }
}


// Hello dialog. So the nooblets know the controls
func doHello( stdscr *gc.Window) {
    stdscr.Refresh()
    stdscr.MovePrintln(5,5, "Controls:")
    stdscr.Println("Escape: Exits the program.")
    stdscr.Println("Tab: Switches between the input window and the Message window.")
    stdscr.Println("Return: Sends a message")
    stdscr.Println(`Shift + Right Arrow: puts in a new line '\n' character (like shift+return in facebook chat)`)
    stdscr.Println("Press any key to continue...")
    stdscr.GetChar()
    stdscr.Clear()
    stdscr.Refresh()
}

// free memory for various things at the end of the program
func cleanup (menu_items []*gc.MenuItem, contactMenu *gc.Menu) {
    for i:=0; i <len(menu_items); i++  {
            menu_items[i].Free()
    }
    contactMenu.UnPost()
    contactMenu.Free()
}

// makes the menu inside the contacts window
func makeContactsMenu(contacts []ts.Contact, contactsWin *gc.Window) ([]*gc.MenuItem, *gc.Menu, *gc.Window) {
    menu_items := make([]*gc.MenuItem, len(contacts))
    var err error
    for i, val := range contacts {
        menu_items[i], err = gc.NewItem((val.Name), "")
        if err != nil {
            log.Fatal("Error making item for contact menu... ", err)
        }
    }
    contactMenu, err := gc.NewMenu(menu_items)
    if err != nil {
        log.Fatal("Error creating contact menu from menu_items... ", err)
    }
    contactsWinSizeY, contactsWinSizeX := contactsWin.MaxYX()
    contactsWin.Keypad(true)
    contactsMenuWin := contactsWin.Derived((contactsWinSizeY-5),(contactsWinSizeX-2),3,1)
    contactMenu.SetWindow(contactsMenuWin)
    contactMenu.Format(len(contacts),1)
    contactMenu.Mark(" * ")

    title := "Contacts"
    contactsWin.MovePrint(1, (contactsWinSizeX/2)-(len(title)/2), title)
    contactsWin.HLine(2, 1, '-', contactsWinSizeX-2)
    contactsMenuWin.Keypad(true)
    contactsMenuWin.ScrollOk(true)
    return menu_items, contactMenu, contactsMenuWin
}

// creates a curses based TUI for the textsecure library
func main() {
    client := &ts.Client{
        RootDir:        ".",
        ReadLine:       ts.ConsoleReadLine,
        MessageHandler: recieveMessage,
    }
    ts.Setup(client)
    stdscr, err := gc.Init()
    if err != nil {
        log.Fatal("Error initializing curses:", err)
    }
    defer gc.End()
    configCurses(stdscr)
    doHello(stdscr)

    contacts, err := ts.GetRegisteredContacts()
    if err != nil {
        log.Fatal("Could not get contacts: %s\n", err)
    }

    telToName = make(map[string]string)
    for _, c := range contacts {
        telToName[c.Tel] = c.Name
    }
    nameToTel = make(map[string]string)
    for _, c := range contacts {
        nameToTel[c.Name] = c.Tel
    }

    contactsWin, messageWinBorder, msgWin, inputBorderWin, inputWin := createMainWindows(stdscr)
    menu_items, contactMenu, contactsMenuWin := makeContactsMenu(contacts, contactsWin)
    globalInputWin = inputWin
    globalMsgWin = msgWin

    contactsMenuWin.Touch()
    contactMenu.Post()
    contactsWin.Refresh()
    messageWinBorder.Refresh()
    inputBorderWin.Refresh()
    msgWin.Refresh()
    msgWinSize_y,msgWinSize_x = msgWin.MaxYX()
    //msgWin.MovePrint(0,0,"test")
    //msgWin.Refresh()

    inputWin.Move(0,0)
    currentContact = getTel(contactMenu.Current(nil).Name())
    go ts.ListenForMessages()
    inputHandler(inputWin, stdscr, contactsMenuWin, contactMenu, msgWin)
    cleanup(menu_items, contactMenu)
}

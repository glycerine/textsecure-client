// Copyright (c) 2014  by f4lk0r_ 
// Licensed under the GPLv3, see the License.txt file for details

package main

import (
    //janimo's textsecure library. Documentation here: https://godoc.org/github.com/janimo/textsecure
    ts "github.com/janimo/textsecure"
//    //go ncurses library. Documentation here: https://godoc.org/code.google.com/p/goncurses
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
    pass string
    placer = -1
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

// In addition to sending a message using janimo's library, also clears screen and resets buffer
func sendMsg (inputWin *gc.Window, msgWin *gc.Window) {
    if len(inputBuffer) != 0 {
        msg := string(inputBuffer)
        to := currentContact
        err := ts.SendMessage(to,msg)
        if err != nil {
            gc.End()
            log.Fatal("SendMessage failed yo: ",err)
        }
        printToMsgWindow(msg,msgWin, true)
        insertMessage("You",currentContact,[]byte(msg),nil)
        inputBuffer = []byte{}
        inputWin.Erase()
    }
}

// free memory for various things at the end of the program
func cleanup (menu_items []*gc.MenuItem, contactMenu *gc.Menu) {
    for i:=0; i <len(menu_items); i++  {
            menu_items[i].Free()
    }
    contactMenu.UnPost()
    contactMenu.Free()
}

// creates a curses based TUI for the textsecure library
func main() {
    client := &ts.Client{
        RootDir:        ".",
        ReadLine:       ts.ConsoleReadLine,
        MessageHandler: recieveMessage,
    }
    stdscr, err := gc.Init()
    if err != nil {
        log.Fatal("Error initializing curses:", err)
    }
    defer gc.End()
    configCurses(stdscr)

    for i:=0;i<3;i++ {
        doHello(stdscr)
        var locked bool = passphraseUnlock(client)
        if locked == false {
            stdscr.ColorOn(1)
            stdscr.MovePrint(0,0,"Password incorrect, hit any key...")
            stdscr.ColorOff(1)
            stdscr.GetChar()
            stdscr.Clear()
            stdscr.Refresh()
        } else {
            stdscr.Clear()
            stdscr.Refresh()
            break
        }
        if i == 2 {
            log.Fatal("Password wrong too many times. Exiting...")
        }
    }
    ts.Setup(client)
    db = setupDatabase()


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

    currentContact = getTel(contactMenu.Current(nil).Name())
    changeContact(contactsMenuWin,contactMenu)
    inputWin.Move(0,0)
    go ts.ListenForMessages()
    inputHandler(inputWin, stdscr, contactsMenuWin, contactMenu, msgWin)
    cleanup(menu_items, contactMenu)
}

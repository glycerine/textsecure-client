// Copyright (c) 2014  by f4lk0r_ 
// Licensed under the GPLv3, see the License.txt file for details

package main

import (
    "bufio"
//    "encoding/hex"
    "log"
    "os"
    "path/filepath"
    //janimo's textsecure library. Documentation here: https://godoc.org/github.com/janimo/textsecure
	ts "github.com/janimo/textsecure"
    //go ncurses library. Documentation here: https://godoc.org/code.google.com/p/goncurses
    gc "code.google.com/p/goncurses"
    "golang.org/x/crypto/bcrypt"
)

// Unlocks the passphrase to the chat database
func passphraseUnlock(c *ts.Client) bool {
    passFile := filepath.Join(c.RootDir, ".storage/pass")
    file, err := os.Open(passFile)
    defer file.Close()
    if err != nil {
        if os.IsNotExist(err) {
            createPassphrase(passFile)
            return true
        } else {
            log.Fatal("Problem accessing passphrase file:",err)
        }
    } else {
        readFile := bufio.NewScanner(file)
        readFile.Scan()
        var str []byte = readFile.Bytes()

        stdscr := gc.StdScr()
        y,x := stdscr.MaxYX()
        stdscr.Println()
        stdscr.MovePrintln((2*int(y/3))-2,int(x/2)-16,"Please enter your passphrase:")
        stdscr.Move((2*int(y/3))-1,int(x/2)-16)

        pass = getPass()
        var passByte []byte = []byte(pass)
        err = bcrypt.CompareHashAndPassword(str,passByte) // this is constant time, I checked, it uses subtle on the backend
        if err == nil{
            return true
        }
    }
    return false
}

// creates the passphrase for the chat database using bcrypt library
func createPassphrase(passFile string) {

    stdscr := gc.StdScr()
    y,x := stdscr.MaxYX()
    stdscr.MovePrintln((2*int(y/3)) - 2,int(x/2)-37,"It appears you do not have a Passphrase. Please enter one now (ASCII only):")
    stdscr.Move((2*int(y/3))-1,int(x/2)-37)

    file, err := os.Create(passFile)
    if err != nil {
        log.Fatal("Problem creating passphrase file", err)
    }
    defer file.Close()
    w := bufio.NewWriter(file)
    pass = getPass()
    var passByte []byte = []byte(pass)

    var bHash []byte
    bHash, err = bcrypt.GenerateFromPassword(passByte,10)
    if err != nil {
        log.Fatal("Problem when bcrypt hashing passphrase",err)
    }

    _,err = w.Write(bHash)
    if err != nil {
        log.Fatal("Error writing passhash to File", err)
    }
    file.Sync()
    w.Flush()
}

// Copyright (c) 2014  by f4lk0r_ 
// Licensed under the GPLv3, see the License.txt file for details

package main

import(
    "log"
    "os"
)


// This is here for conveniences sake when testing, for logging out to a file to not interfere with curses
func startDebugLogger() *log.Logger {
    file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Println("Failed to open log file", ":", err)
        panic(err)
        return nil
    }

    MyFile := log.New(file,
        "PREFIX: ",
        log.Ldate|log.Ltime|log.Lshortfile)
    return MyFile
}

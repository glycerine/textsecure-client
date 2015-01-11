package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"github.com/janimo/textsecure"
    "os"
    "log"
//    "reflect"
)

// Simple command line test app for TextSecure.
// It can act as an echo service, send one-off messages and attachments,
// or carry on a conversation with another client

var (
	echo       bool
	to         string
	message    string
	attachment string
)

func init() {
	flag.BoolVar(&echo, "echo", false, "Act as an echo service")
	flag.StringVar(&to, "to", "", "Contact name to send the message to")
	flag.StringVar(&message, "message", "", "Single message to send, then exit")
	flag.StringVar(&attachment, "attachment", "", "File to attach")
}

var (
	red   = "\x1b[31m"
	green = "\x1b[32m"
	blue  = "\x1b[34m"
)

// conversationLoop sends messages read from the console
func conversationLoop() {
	for {
		message := textsecure.ConsoleReadLine(fmt.Sprintf("%s>", blue))
		if message == "" {
			continue
		}
		err := textsecure.SendMessage(to, message)
		if err != nil {
			debugLog.Println(err)
		}
	}
}

func messageHandler(msg *textsecure.Message) {
	if echo {
		err := textsecure.SendMessage(msg.Source(), msg.Message())
		if err != nil {
			debugLog.Println(err)
		}
		return
	}

	if msg.Message() != "" {
		fmt.Printf("\r                                               %s%s : %s%s%s\n>", red, getName(msg.Source()), green, msg.Message(), blue)
	}

	for _, a := range msg.Attachments() {
		handleAttachment(msg.Source(), a)
	}

	// if no peer was specified on the command line, start a conversation with the first one contacting us
	if to == "" {
		to = msg.Source()
		go conversationLoop()
	}
}

func handleAttachment(src string, b []byte) {
	f, err := ioutil.TempFile(".", "TextSecure_Attachment")
	if err != nil {
		debugLog.Println(err)
		return
	}
	debugLog.Printf("Saving attachment of length %d from %s to %s", len(b), src, f.Name())
	f.Write(b)

}

// getName returns the local contact name corresponding to a phone number,
// or failing to find a contact the phone number itself
func getName(tel string) string {
	if n, ok := telToName[tel]; ok {
		return n
	}
	return tel
}

func startDebugLogger() *log.Logger {
    // The default, built-in logging package in Go does not have the ability to log out to a file.
    // I hate this fucking language with the fiery passion of a thousand suns
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

func oldmain() {
	flag.Parse()
	debugLog.SetFlags(0)
	client := &textsecure.Client{
		RootDir:        ".",
		ReadLine:       textsecure.ConsoleReadLine,
		MessageHandler: messageHandler,
	}
	textsecure.Setup(client)

	if !echo {
		contacts, err := textsecure.GetRegisteredContacts()
		if err != nil {
			debugLog.Println("Could not get contacts: %s\n", err)
		}

		telToName = make(map[string]string)
		for _, c := range contacts {
			telToName[c.Tel] = c.Name
		}

		// If "to" matches a contact name then get its phone number, otherwise assume "to" is a phone number
		for _, c := range contacts {
			if strings.EqualFold(c.Name, to) {
				to = c.Tel
				break
			}
		}

		if to != "" {
			// Send attachment with optional message then exit
			if attachment != "" {
				err := textsecure.SendFileAttachment(to, message, attachment)
				if err != nil {
					debugLog.Fatal(err)
				}
				return
			}

			// Send a message then exit
			if message != "" {
				err := textsecure.SendMessage(to, message)
				if err != nil {
					debugLog.Fatal(err)
				}
				return
			}

			// Enter conversation mode
			go conversationLoop()
		}
	}

	err := textsecure.ListenForMessages()
	if err != nil {
		debugLog.Println(err)
	}
}

var telToName map[string]string
var debugLog = startDebugLogger()

func main() {
}

#!/usr/bin/env python3

from curses import wrapper
import curses
import curses.textpad
import yaml
from subprocess import call
import logging

#which person are we talking to
currentContact = ""
contactList = []

logging.basicConfig(filename='debug.log', level=logging.DEBUG)
logging.debug('Opening Debug Log')


def handleInput(c):
    #does not work, apparently enter key is really tweaky in curses
    #not sure if this is worth fixing seeing as you can send with control+g
    #if c == curses.KEY_ENTER:
    #    call(["./textsecure","-to=" + currentContact,"-message=" + data])
    #    inputWin.clear()
    #    inputWin.refresh()
    #    return 0
    #if c == ord("c"):
    #    logging.debug("hit")
    #    logging.debug(str(c))
    #    return "d"
    #if c == 27:
    #    #logging.debug("AYE MATEY WE HIT THE GREAT BIG WHALE")
    #    return "s"
    #logging.debug(c)
    return c

def main(stdscr):
    # Clear screen
    stdscr.clear()
    #curses.nonl()
    curses.start_color()
    # Set up our colors
    curses.init_pair(1, curses.COLOR_RED, curses.COLOR_BLACK)
    curses.init_pair(2, curses.COLOR_BLUE, curses.COLOR_BLACK)
    curses.init_pair(3, curses.COLOR_GREEN, curses.COLOR_BLACK)
    curses.init_pair(4, curses.COLOR_YELLOW, curses.COLOR_BLACK)
    curses.init_pair(5, curses.COLOR_CYAN, curses.COLOR_BLACK)
    curses.init_pair(6, curses.COLOR_MAGENTA, curses.COLOR_WHITE)
    curses.init_pair(7, curses.COLOR_MAGENTA, curses.COLOR_BLACK)

    global contactList
    global currentContact
    
    currentContact = "Phone Me"
    
    # Draw the 3 windows that will show on the screen. A Contacts panel (possibly bigger than the screen)
    # A message body panel that shows the current conversation
    # And a textinput box for you to input messages to send. 

    begin_x = 0; begin_y = 0; width = int(curses.COLS * .2)
    # The height of the pad is dependent upon the number of contacts the user has, or its just the size of the screen. 
    f = open("./.config/contacts.yml","r")
    yamldata = yaml.load(f)
    contactList = yamldata['contacts']

    if ((int(curses.LINES - 1)) < (len((yamldata['contacts'])) * 2 + 2)):
        height = (len((yamldata['contacts'])) * 2 + 6)
    else:
        height = (int(curses.LINES - 1))

    contactsWin = curses.newpad(height, width)
    stdscr.clear()
    contactsWin.clear()
    contactsWin.border("|","|","-","-","+","+","+","+")
    contactsWin.refresh(0,0, begin_x,begin_y, width,(int(curses.LINES -1)))

    begin_x = int(curses.COLS * .2)+1
    height = int(curses.LINES * .8 ); width = int(curses.COLS * .8)
    messageWin = curses.newwin(height, width, begin_y, begin_x)
    messageWin.clear()
    messageWin.border("|","|","-","-","+","+","+","+")
    messageWin.refresh()

    begin_y = int(curses.LINES * .8 ) + 1
    height = int(curses.LINES * .2) - 1
    # Because this curses module is stupid, I am forced to create a 
    # subwindow inside this window in order to get textboxes to behave 
    # as expected... otherwise it intermixes the border with user input.
    # so inputWinBorder is a window whose only purpose is to provide the border
    # and then inputWin is the actual window where the textbox is.
    # God bless this POS.
    inputWinBorder = curses.newwin(height, width, begin_y, begin_x)
    inputWinBorder.clear()
    inputWinBorder.refresh()
    
    # populate the contacts pane with data from the contacts.yml file
    contactsWin.addstr(1,1,"Contacts:", curses.color_pair(2))
    for i in range(0,(len(contactList))):
        # if there's more contacts than can fit on the screen, don't draw the extras... user has to scroll down to see those.
        if (2 + (2*(i+2)) < (curses.LINES - 1)):
            contactsWin.addstr((2*(i+2)),1,(contactList[i])['name'], curses.color_pair(0))

    contactsWin.border("|","|","-","-","+","+","+","+")
    contactsWin.refresh(0,0, 0,0, (int(curses.LINES -1)),(int(curses.COLS * .2)))

    # draw the inputbox
    inputWinBorder.border("|","|","-","-","+","+","+","+")
    inputWinBorder.refresh()
    winsize = inputWinBorder.getmaxyx()
    inputWin = inputWinBorder.derwin((winsize[0]-2),(winsize[1]-2),1,1)
    inputBox = curses.textpad.Textbox(inputWin)
    inputWin.move(0,0)

    # Aight everything's set up, let's handle some input
    while True:
        data = ""
        data = inputBox.edit(handleInput)
        call(["./textsecure","-to=" + currentContact,"-message=" + data])
        inputWin.clear()
        inputWin.refresh()



wrapper(main)

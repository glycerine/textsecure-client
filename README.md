# TextSecure for Terminal

Status
-------
Currently, this needs to be ported to the latest version of janimo's API and does not work. It is on my backburner to get it up and running again... Pull requests appreciated!


Many thanks to the people behind golang and all its libraries, ncurses, Rob Thornton for his [goncurses] (https://code.google.com/p/goncurses/) library, @janimo for his [golang textsecure] (https://github.com/janimo/textsecure.git) library, and @mattn for his [go-sqlite3] (https://github.com/mattn/go-sqlite3) driver.

Screenshot
----------

![screenshot of textsecure for terminal](https://github.com/f41c0r/textsecure-client/wiki/screenshots/output.gif)

Hint: click on the gif for full screen better quality image

Installation
------------

    go get github.com/f41c0r/textsecure-client

For more details, including setting up Go, check out janimo's [wiki] (https://github.com/janimo/textsecure/wiki/Installation)

cd to the directory where the application is located. for example:
    
    cd $GOPATH/src/github.com/f41c0r/textsecure-client

then build the application with go build

    go build

and then modify your configuration settings, such as phone numbers, etc etc

    nano .config/config.yml

and then add your contacts to the application

    nano .config/contacts.yml

and run the application!
    
    ./textsecure-client 

Configuration
-------------

Before starting the Application for the first time, go into the .config directory and edit contacts.yml and config.yml to your liking. You should then have everything you need to get started.

Usage
-----

This works, it's buggy as all hell and in alpha but you can use it to send and recieve messages. No group support yet, but that's todo.

Discussions
-----------

User and developer discussions happen on the [mailing list] (https://groups.google.com/forum/#!forum/textsecure-go)

License
-------

This code in this repository is under the GNU GPL v3.0. The Go programming language, the bcrypt library, and the goncurses library are all under a 3-clause BSD License, @janimo's textsecure library is under the GPL v3.0, and mattn's go-sqlite3 driver is under an MIT License. That means that all the code required for this software to run is Free, Libre and Open Source Software (FLOSS).

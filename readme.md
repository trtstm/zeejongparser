ZeeJong Parser
==============

About
-----

*ZeeJong Parser* is a tool, written in [*GO*](http://golang.org), which collects and builds a database containing an archive of soccer information. This tool is written for the *ZeeJong* web app.


Installing & Running
--------------------

Copy the folder to the `$GOPATH/src` folder.  
Build the program:

    go build

Run the `zeejongparser` executable

    ./zeejongparser


*See below for input/output information*

Usage
-----

The parser's input consists of a series of football competitions, with their corresponding archive url. The parser will parse these competitions in a concurrent way. All pages are cached in the `cache` folder, in case you'd like to run it more than once. Images are saved in the `images` folders. The names of those images correspond to the id of the 'owner'. (e.g. the player with id `512`, will have an image named `Player-512.png`)

While parsing, you can see the stats on a webpage. Open your browser and go to `http://localhost:8080`. This page will you show you the progress of the parser.

The output of the parser is a JSON file, containing the whole database.

You can import this JSON file in the main *ZeeJong* installation, running the importer.

> You can also install preloaded data. We provide a JSON file contain an archive of competitions.  
> In order to install this, you must run yourdomain.com/core/importer.php. This may take some time.  
> Do this after the installation. Note that the script will empty all tables before loading the sample data.

*More information on installing and setting up ZeeJong, can be found in the readme of the ZeeJong project.*



Authors
-------

- Timo Truyts
- Mathias Beke


Acknowledgements
----------------

- Parsing/browsing html: [GoQuery](https://github.com/PuerkitoBio/goquery)
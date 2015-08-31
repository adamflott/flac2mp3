# About

flac2mp3, a command line tool to parallel conversion of flacs to V0 mp3s.

# Usage

    $ ./flac2mp3 -h
    NAME:
       flac2mp3 - parallel conversion of flacs to V0 mp3s
    
    USAGE:
       flac2mp3 [global options] command [command options] [arguments...]
       
    VERSION:
       1
       
    AUTHOR(S):
       Adam Flott <adam@adamflott.com> 
       
    COMMANDS:
       scan, s	scan
       convert, c	convert
       help, h	Shows a list of commands or help for one command
       
    GLOBAL OPTIONS:
       --source-directory "/home/adam/src/adam/flac2mp3"			<dir>
       --destination-directory "/home/adam/src/adam/flac2mp3/destination"	<dir>
       --help, -h								show help
       --version, -v							print the version
       
    COPYRIGHT:
       2015
    
    $ ./flac2mp3 --source-direcotry /music/flacs --destination-directory /music/mp3s convert

# Notes

* defaults to all available CPUs as reported by runtime.NumCPU()
* hard coded to only do V0 (from lame) MP3s
* works for me

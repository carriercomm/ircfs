// Copyright 2011 Alex Ray. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
    "flag"
    "log"
    "github.com/ajray/go-fuse/fuse"
    "os"
//    "time"
//    "os/signal"
)

var (
    n    = 3
    nickFile = NewNickFile("ajray")
    ctlFile  = NewCtlFile()
)

// Flags
var (
    fa = flag.String("a","addr","Use address as irc server.")
    ff = flag.String("f","fromhost","Make the connection from ip address or host name fromhost.")
    fn = flag.String("n","nick","Use nick as nickname when connecting.")
    fl = flag.String("l","logpath","Log contents of the data files of all the irc dirs to logpath/.")
    ft = flag.Bool("t",false,"Use long timestamps, in both logging and the data file.")
    fd = flag.Bool("d",false,"Enable debugging.")
)

type IrcFs struct {
    fuse.DefaultFileSystem
}

func (me *IrcFs) GetAttr(name string) (*os.FileInfo, fuse.Status) {
    log.Print("GetAttr " + name)
    switch name {
    case "file.txt":
        return &os.FileInfo{Mode: fuse.S_IFREG | 0444, Size: int64(len(name))}, fuse.OK
    case "ctl":
        return &os.FileInfo{Mode: fuse.S_IFREG | 0222, Size: int64(len(name))}, fuse.OK
    case "event":
        return &os.FileInfo{Mode: fuse.S_IFREG | 0444, Size: int64(len(name))}, fuse.OK
    case "nick":
        return &os.FileInfo{Mode: fuse.S_IFREG | 0444, Size: int64(len(nickFile.Nick))}, fuse.OK
    case "raw":
        return &os.FileInfo{Mode: fuse.S_IFREG | 0666, Size: int64(len(name))}, fuse.OK
    case "pong":
        return &os.FileInfo{Mode: fuse.S_IFREG | 0444, Size: int64(len(name))}, fuse.OK
    case "":
        return &os.FileInfo{Mode: fuse.S_IFDIR | 0755}, fuse.OK
    }
    return nil, fuse.ENOENT
}

func (me *IrcFs) OpenDir(name string) (stream chan fuse.DirEntry, code fuse.Status) {
    log.Print("OpenDir " + name)
    if name == "" {
        stream = make(chan fuse.DirEntry, 6) // , n + 5) // MUAHAHA NO BUFFER
        stream <- fuse.DirEntry{Name: "file.txt", Mode: fuse.S_IFREG}
        stream <- fuse.DirEntry{Name: "ctl", Mode: fuse.S_IFREG}
        stream <- fuse.DirEntry{Name: "event", Mode: fuse.S_IFREG}
        stream <- fuse.DirEntry{Name: "nick", Mode: fuse.S_IFREG}
        stream <- fuse.DirEntry{Name: "raw", Mode: fuse.S_IFREG}
        stream <- fuse.DirEntry{Name: "pong", Mode: fuse.S_IFREG}
        close(stream)
        return stream, fuse.OK
    }
    return nil, fuse.ENOENT
}

func (me *IrcFs) Open(name string, flags uint32) (file fuse.File, code fuse.Status) {
    log.Print("Open " + name)
    switch name {
    case "file.txt":
        return fuse.NewReadOnlyFile([]byte(name)), fuse.OK
    case "ctl":
        return ctlFile, fuse.OK
    case "event":
        return fuse.NewReadOnlyFile([]byte(name)), fuse.OK
    case "nick":
        return nickFile, fuse.OK
    case "raw":
        return fuse.NewReadOnlyFile([]byte(name)), fuse.OK
    case "pong":
        return fuse.NewReadOnlyFile([]byte(name)), fuse.OK
    }
    return nil, fuse.ENOENT
}

func main() {
    flag.Parse()
    if len(os.Args) < 2 {
        log.Fatal("Usage:  ircfs MOUNTPOINT")
    }
    state, _, err := fuse.MountFileSystem(os.Args[1], &IrcFs{}, nil)
    if err != nil {
        log.Fatal("Mount fail:", err)
    }
    //go func() { for { sig := <-signal.Incoming
    //        time.Sleep(1) // FIXME TODO XXX hack to make the goroutine scheduler switch
    //        log.Print("Reading signal: " + sig.String() ) } }()
    state.Loop(true)
}

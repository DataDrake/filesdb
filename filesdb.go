package main

import (
	"flag"
	"fmt"
	"github.com/DataDrake/filesdb/core"
	"os"
	"bufio"
)

func usage() {
	fmt.Println("Usage: filesdb [CMD] [QUERY]")
}

func main() {
	flag.Usage = func() { usage() }
	flag.Parse()

	args := flag.Args()
	switch len(args) {
	case 2:
		if args[0] != "search" {
			usage()
			os.Exit(1)
		}
		f, err := os.Open("/var/log/filesdb/db.cbor")
		if err != nil {
			fmt.Printf("Could not open database, reason: %s\n", err.Error())
			os.Exit(1)
		}
		db := bufio.NewReader(f)
		filesdb.Search(args[1], db)
		f.Close()
	case 1:
		if args[0] != "update" {
			usage()
			os.Exit(1)
		}
		dir, err := os.Open("/var/log/filesdb")
		if err != nil {
			err = os.Mkdir("/var/log/filesdb", 0744)
			if err != nil {
				fmt.Printf("Could not create filesdb directory, reason: %s\n", err.Error())
				os.Exit(1)
			}
		}
		dir.Close()
		f, err := os.Create("/var/log/filesdb/db.cbor")
		if err != nil {
			fmt.Printf("Could not create db.cbor, reason: %s\n", err.Error())
			os.Exit(1)
		}
		db := bufio.NewWriter(f)
		filesdb.Fill("/", db)
		db.Flush()
		f.Sync()
		f.Close()
	default:
		usage()
		os.Exit(1)
	}

}

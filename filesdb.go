package main

import (
	"flag"
	"fmt"
	"github.com/DataDrake/filesdb"
	"os"
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
		t := filesdb.NewTree()
		err = t.FromCBOR(f)
		f.Close()
		if err != nil {
			fmt.Printf("Could not import database, reason: %s\n", err.Error())
			os.Exit(1)
		}
		t.Search(args[1])
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
		db, err := os.Create("/var/log/filesdb/db.cbor")
		if err != nil {
			fmt.Printf("Could not create db.cbor, reason: %s\n", err.Error())
			os.Exit(1)
		}
		t := filesdb.NewTree()
		t.Fill("/")
		t.ToCBOR(db)
		db.Close()
	default:
		usage()
		os.Exit(1)
	}

}

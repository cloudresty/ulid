package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	ulid "github.com/cloudresty/ulid"
)

var (
	versionFlag = flag.Bool("version", false, "Print version information")
	timeFlag    = flag.Uint64("time", 0, "Generate ULID with specified timestamp (milliseconds)")
	version     string // This will be set during build
)

func main() {

	flag.Parse()

	if *versionFlag {
		fmt.Println("ulid", version, "- https://github.com/cloudresty/ulid")
		os.Exit(0)
	}

	var ulidStr string
	var err error

	if *timeFlag > 0 {
		ulidStr, err = ulid.NewTime(*timeFlag)
	} else {
		ulidStr, err = ulid.New()
	}

	if err != nil {
		log.Fatalf("Error generating ULID: %v", err)
	}

	fmt.Println(ulidStr)

}

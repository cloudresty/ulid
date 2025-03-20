package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	ulid "github.com/cloudresty/goulid"
)

var (
	versionFlag = flag.Bool("version", false, "Print version information")
	timeFlag    = flag.Uint64("time", 0, "Generate ULID with specified timestamp (milliseconds)")
)

func main() {

	flag.Parse()

	if *versionFlag {
		fmt.Println("ulid v1.0.5")
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

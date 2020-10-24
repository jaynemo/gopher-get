
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
    "os"
	"strings"


	"github.com/prologic/go-gopher"
)

var (
	json = flag.Bool("json", false, "display gopher directory as JSON")
)


func check(e error) {
    if e != nil {
        panic(e)
        	os.Exit(1)

    }
}

func fatal(format string, a ...interface{}) {
	
    format = "*** " + format + " ***\n"

	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)

}

func saveFile(contents []byte, path string) {
    
    //fmt.Print("write to "+ path)

    
    d1 := []byte(contents)
    err := ioutil.WriteFile(path, d1, 0644)
    
    
    check(err)
    
}

func main() {
	var uri string
    var out string

	flag.Parse()
    
    
    uri = flag.Arg(0)
    out = flag.Arg(1)
	
	if uri == "" {
		log.Fatal("You must pass a Gopher URI")
	}
	if out == "" {
		log.Fatal("You must pass a file to write to")
	}
	res, err := gopher.Get(uri)
    
    if res != nil {
		fmt.Println(res.Type)
	}
	
	if err != nil {
		if strings.Contains(err.Error(), "invalid scheme for uri") {
			log.Fatal("Only gopher:// URIs are supported\n" + err.Error())
		} else if strings.Contains(err.Error(), "connection failed because connected host has failed to respond.") {
			log.Fatal("Connection timeout- the server did not respond in a reasonable period of time.\n" + err.Error())
		}
		log.Fatal(err)
	}
	
	if res.Body != nil {
		contents, err := ioutil.ReadAll(res.Body)
		if err != nil {
            fmt.Println(err)
			log.Fatal(err)
		}
		
		//Special handling for SDF, other Gophernicus installs
        if string(contents) == "Error: File or directory not found!\r\n" {
			log.Fatal("The server was unable to locate the file at: " + uri + ". " + string(contents))
		}
        //save body to file
        saveFile(contents, out)
        
	} else {
		var (
			bytes []byte
			err   error
		)

		if *json {
			bytes, err = res.Dir.ToJSON()
		} else {
			bytes, err = res.Dir.ToText()
		}
		if err != nil {
			log.Fatal(err)
		}

        //save dir contents or other bytest to file
        saveFile(bytes, out)
	}
}

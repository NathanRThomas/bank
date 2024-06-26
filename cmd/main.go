/** ****************************************************************************************************************** **
	The Bank

** ****************************************************************************************************************** **/

package main

import (
	"github.com/jessevdk/go-flags"
	
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
	"syscall"
	"context"
)

  //-------------------------------------------------------------------------------------------------------------------//
 //----- CONSTS ------------------------------------------------------------------------------------------------------//
//-------------------------------------------------------------------------------------------------------------------//

// give us a name
const apiName = "The Bank"
const apiVersion = "0.2.0"

  //-------------------------------------------------------------------------------------------------------------------//
 //----- CONFIG ------------------------------------------------------------------------------------------------------//
//-------------------------------------------------------------------------------------------------------------------//

// final local options object for this executable
var opts struct {
	Help bool `short:"h" long:"help" description:"Shows help message"`
	Port string `short:"p" long:"port" description:"Specifies the target port to run on"`
	Ledgers string `short:"l" long:"ledgers" description:"Directory with the ledgers" default:"/var/www/ledgers"`
	Templates string `short:"t" long:"templates" description:"Directory with the templates" default:"/var/www/go/bank/templates"`
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information -v max of -vv"`
}

  //-------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS ---------------------------------------------------------------------------------------------------//
//-------------------------------------------------------------------------------------------------------------------//

func showHelp () {
	fmt.Printf("*****************************\n%s : Version %s\n\n", apiName, apiVersion)

	fmt.Printf("\n*****************************\n")
}

// handles parsing command arguments as well as setting up our opts object
func parseCommandLineArgs () ([]string) {
	// parse things
	args, err := flags.Parse(&opts)
	if err != nil { log.Fatal(err) }

	if opts.Help {
		showHelp()
		os.Exit(0)
	}

	// see what they're trying to do here
	if len(opts.Port) == 0 {
		opts.Port = "8080" // default port
	}

	// check any args
	for _, arg := range args {
		switch strings.ToLower(arg) {
		case "help":
			showHelp()
			os.Exit(0)

		case "version":
			fmt.Printf("%s v%s\n", apiName, apiVersion)
			os.Exit(0)
		}
	}

	return args // return any arguments we don't know what to do with... yet

}

  //-------------------------------------------------------------------------------------------------------------------//
 //----- APP ---------------------------------------------------------------------------------------------------------//
//-------------------------------------------------------------------------------------------------------------------//

type app struct {
	running bool 
}

  //-------------------------------------------------------------------------------------------------------------------//
 //----- FUNCTIONS ---------------------------------------------------------------------------------------------------//
//-------------------------------------------------------------------------------------------------------------------//

  //-------------------------------------------------------------------------------------------------------------------//
 //----- MAIN --------------------------------------------------------------------------------------------------------//
//-------------------------------------------------------------------------------------------------------------------//

func main() {
	log.SetFlags(log.Lshortfile)
	// first step, parse the command line params
	parseCommandLineArgs()

	log.Printf("Starting %s v%s\n", apiName, apiVersion)

	// main app for everything
	app := &app{
		running: true,
	}

	// create our server server
	srv := &http.Server {
		Addr: fmt.Sprintf(":%s", opts.Port),
		Handler: app.routes(), 
		ReadTimeout: time.Second * 30,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-c // this sits until something comes into the channel, eg the notify interupts from above
		app.running = false

		srv.Shutdown(context.Background()) // shutdown the server
	}()

	log.Printf("%s v%s started on port %s\n", apiName, apiVersion, opts.Port) // going to always record this starting message
	if err := srv.ListenAndServe(); err != http.ErrServerClosed { // Error starting or closing listener:
		log.Printf("ListenAndServe: %v", err) // we want to know if this failed for another reason
	}

	log.Println("exiting...")
	
	os.Exit(0) //final exit
}

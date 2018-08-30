package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"runtime/trace"
	_ "github.com/emretiryaki/merkut/pkg/services/sqlstore"
	"github.com/emretiryaki/merkut/pkg/log"
	 srv "github.com/emretiryaki/merkut/pkg/server"
)

var version = "1.0.0"
var commit = "NA"


var configFile = flag.String("config","","path config file")
var homePath = flag.String("homepath", "", "defaults to working directory")

func  main(){

	v := flag.Bool("v", false, "prints current version and exits")


	flag.Parse()

	if *v {
		fmt.Printf("Version %s (commit: %s)\n", version, commit)
		os.Exit(0)
	}

	server := srv.NewMerkutServer()

	go listenToSystemSignals(server)

	err :=server.Run(*configFile,*homePath,version,commit)

	code :=server.Exit(err)

	trace.Stop()

	log.Close()

	os.Exit(code)

}

func listenToSystemSignals(server *srv.MerkutServerImpl){
	signalChan := make(chan os.Signal,1)
	ignoreChan := make(chan os.Signal, 1)

	signal.Notify(ignoreChan, syscall.SIGHUP)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		server.Shutdown(fmt.Sprintf("System signal: %s", sig))
	}
}


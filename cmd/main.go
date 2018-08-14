package main

import (
	"flag"
	"fmt"
	"os"
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/emretiryaki/merkut/pkg/setting"
	"github.com/emretiryaki/merkut/pkg/log"
)

var version = "1.0.0"
var commit = "NA"


var configFile = flag.String("config","","path config file")
func  main(){

	v := flag.Bool("v", false, "prints current version and exits")


	flag.Parse()

	if *v {
		fmt.Printf("Version %s (commit: %s)\n", version, commit)
		os.Exit(0)
	}



}

type MerkutServerImpl struct {
	context            context.Context
	shutdownFn         context.CancelFunc
	childRoutines      *errgroup.Group
	log                log.Logger
	cfg                *setting.Cfg
	shutdownReason     string
	shutdownInProgress bool

	//RouteRegister routing.RouteRegister `inject:""`
	//HttpServer    *api.HTTPServer       `inject:""`
}
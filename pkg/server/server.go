package server

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/facebookgo/inject"
	"golang.org/x/sync/errgroup"

	"github.com/emretiryaki/merkut/pkg/api"
	"github.com/emretiryaki/merkut/pkg/bus"
	"github.com/emretiryaki/merkut/pkg/log"
	"github.com/emretiryaki/merkut/pkg/registry"
	"github.com/emretiryaki/merkut/pkg/routing"
	"github.com/emretiryaki/merkut/pkg/setting"
	_ "github.com/emretiryaki/merkut/pkg/services/alerting/notifiers"
	_ "github.com/emretiryaki/merkut/pkg/services/alerting"
)

type MerkutServerImpl struct {
	context            context.Context
	shutdownFn         context.CancelFunc
	childRoutines      *errgroup.Group
	log                log.Logger
	cfg                *setting.Cfg
	shutdownReason     string
	shutdownInProgress bool
	RouteRegister routing.RouteRegister `inject:""`
	HttpServer    *api.HTTPServer       `inject:""`
}

func NewMerkutServer() *MerkutServerImpl {
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	childRoutines, childCtx := errgroup.WithContext(rootCtx)

	return &MerkutServerImpl{
		context:       childCtx,
		shutdownFn:    shutdownFn,
		childRoutines: childRoutines,
		log:           log.New("server"),
		cfg:           setting.NewCfg(),


	}
}

func (m *MerkutServerImpl) Shutdown(reason string)  {
	m.log.Info("Shutdown started", "reason", reason)
	m.shutdownReason = reason
	m.shutdownInProgress = true

	m.shutdownFn()

	m.childRoutines.Wait()

}
func (m *MerkutServerImpl) Run(configFile string, homePath string ,version string , commit string)  error{

	m.loadConfiguration(configFile,homePath,version,commit)


	serviceGraph := inject.Graph{}
	serviceGraph.Provide(&inject.Object{Value:bus.GetBus()})
	serviceGraph.Provide(&inject.Object{Value:m.cfg})
	serviceGraph.Provide(&inject.Object{Value: routing.NewRouteRegister()})

	services := registry.GetServices()

	for _, service := range services {
		serviceGraph.Provide(&inject.Object{Value: service.Instance})
	}

	serviceGraph.Provide(&inject.Object{Value: m})


	if err := serviceGraph.Populate(); err != nil{
		return fmt.Errorf("Failed to register dependecy : %v",err)
	}

	for _,service := range services{
		if registry.IsDisabled(service.Instance){
			continue
		}

		m.log.Info("Initializing " + service.Name)

		if err := service.Instance.Init() ; err != nil{
			return fmt.Errorf("service init failed : %v", err)
		}
	}

	for _,srv := range services {

		descriptor := srv

		service , ok := srv.Instance.(registry.BackgroundService)
		if !ok{
			continue
		}

		if registry.IsDisabled(descriptor.Instance) {
			continue
		}

		m.childRoutines.Go(func() error {
			if m.shutdownInProgress{
				return nil
			}

			err :=service.Run(m.context)

			// If error is not canceled then the service crashed
			if err != context.Canceled && err != nil {
				m.log.Error("Stopped "+descriptor.Name, "reason", err)
			} else {
				m.log.Info("Stopped "+descriptor.Name, "reason", err)
			}

			// Mark that we are in shutdown mode
			// So more services are not started
			m.shutdownInProgress = true
			return err
		})

	}

	sendSystemdNotification("READY=1")

	return m.childRoutines.Wait()

}

func (g *MerkutServerImpl) Exit(reason error) int {
	// default exit code is 1
	code := 1

	if reason == context.Canceled && g.shutdownReason != "" {
		reason = fmt.Errorf(g.shutdownReason)
		code = 0
	}

	g.log.Error("Server shutdown", "reason", reason)
	return code
}

func (g *MerkutServerImpl) loadConfiguration(configFile string, homePath string ,version string , commit string){
	err := g.cfg.Load(&setting.CommandLineArgs{
		Config:   configFile,
		HomePath: homePath,
		Args:     flag.Args(),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start grafana. error: %s\n", err.Error())
		os.Exit(1)
	}

	g.log.Info("Starting "+setting.ApplicationName, "version", version, "commit", commit, "compiled", time.Unix(setting.BuildStamp, 0))

}

func  sendSystemdNotification(state string)  error{

	notifySocket := os.Getenv("NOTIFY_SOCKET")

	if notifySocket == ""{
		return fmt.Errorf("NOTIFY_SOCKET environment variable empty or unset.")
	}
	socketAddr := &net.UnixAddr{
		Name: notifySocket,
		Net:  "unixgram",
	}

	conn, err := net.DialUnix(socketAddr.Net, nil, socketAddr)

	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(state))

	conn.Close()

	return err

}
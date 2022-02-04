package main

import (
	"log"
	"net"
	"os"

	"github.com/BurntSushi/toml"
	gocron_server "github.com/strnophix/gocron-server/pkg"
	pb "github.com/strnophix/gocron-server/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) < 2 {
		log.Fatalln("Expected a gocron-server config as the first argument.")
	}

	cfgPath := os.Args[1]
	cfg := Config{}
	_, err := toml.DecodeFile(cfgPath, &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	l, err := net.Listen("tcp", cfg.Server.Host)
	if err != nil {
		log.Fatalf("Failed to start server on %s: %v", cfg.Server.Host, err)
	}

	log.Printf("gocron-server is running on %s\n", cfg.Server.Host)

	gs := grpc.NewServer()

	s := gocron_server.NewSchedulerService()
	defer s.Shutdown()

	for _, unitCfg := range cfg.Unit {
		unit := unitCfg.ToSchedulerUnit()
		err = s.AddUnit(unit)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Succesfully added unit %s\n", unit.Name)
	}

	pb.RegisterSchedulerServer(gs, s)
	reflection.Register(gs)

	if err := gs.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

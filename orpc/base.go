package orpc

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/Mtbcooler/outrun/config"
	"github.com/Mtbcooler/outrun/orpc/rpcobj"
)

func Start() {
	rpc.Register(new(rpcobj.Toolbox))
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":"+config.CFile.RPCPort)
	if err != nil {
		log.Fatal("error listening in ORPC:", err)
	}
	log.Println("Starting ORPC server on port " + config.CFile.RPCPort)
	go http.Serve(listener, nil)
}

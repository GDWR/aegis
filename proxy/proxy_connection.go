package proxy

import (
	"aegis/utils"
	"io"
	"net"
)

func ProxyConnection(source net.Conn, destination string) {
	destCon, err := net.Dial("tcp", destination)
	utils.HandleError(err)
	go io.Copy(source, destCon)
	io.Copy(destCon, source)
}

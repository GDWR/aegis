package main

import (
	"aegis/proxy"
	"aegis/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexflint/go-arg"
	"net"
	"net/http"
)

var Args struct {
	Host            string `arg:"-h,--host" default:"0.0.0.0"`
	ControlPort     uint16 `arg:"-c,--control-port" help:"Port to host control http sever" default:"8765"`
	Udp             bool   `arg:"-u,--udp" help:"Use UDP" default:"false"`
	HostPort        uint16 `arg:"-p,--port,required" help:"HostPort to pkg on"`
	Destination     string `arg:"-d,--destination,required" help:"Destination of the server proxying to"`
	DestinationPort uint16 `arg:"--destination-port" help:"HostPort of the host to pkg to, defaults to provided HostPort"`
}

func receiveAddress(request *http.Request) (string, error) {
	var data string
	err := json.NewDecoder(request.Body).Decode(&data)

	if err != nil {
		return "", err
	}

	if net.ParseIP(data) == nil {
		return "", errors.New("input is not a valid ip address")
	}

	return data, nil
}

func main() {
	arg.MustParse(&Args)
	if Args.DestinationPort == 0 {
		Args.DestinationPort = Args.HostPort
	}

	banned := NewBanList()
	go listenAndProxy(banned)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		switch request.Method {

		case http.MethodGet:
			writer.WriteHeader(http.StatusOK)
			err := json.NewEncoder(writer).Encode(banned.banMap)
			utils.HandleError(err)
			break

		case http.MethodPost:
			ip, err := receiveAddress(request)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				_, err := writer.Write([]byte(fmt.Sprintf("%s", err)))
				utils.HandleError(err)
				break
			}

			writer.WriteHeader(http.StatusOK)
			banned.AddBan(ip)
			fmt.Printf("Added client address %s to the banlist\n", ip)
			break

		case http.MethodDelete:
			ip, err := receiveAddress(request)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				_, err := writer.Write([]byte(fmt.Sprintf("%s", err)))
				utils.HandleError(err)
				break
			}

			writer.WriteHeader(http.StatusOK)
			banned.RemoveBan(ip)
			fmt.Printf("Removed client address %s from the banlist\n", ip)
			break

		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			_, err := writer.Write([]byte("Unsupported HTTP method"))
			utils.HandleError(err)
			break
		}
	})

	fmt.Printf("Serving http control api on http://0.0.0.0:%d\n", Args.ControlPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", Args.ControlPort), nil)
	utils.HandleError(err)
}

func listenAndProxy(banList *BanList) {
	connectionString := fmt.Sprintf("0.0.0.0:%d", Args.HostPort)
	fmt.Printf("Listening for connections on tcp://%s to proxy\n", connectionString)

	listener, err := net.Listen("tcp4", connectionString)
	utils.HandleError(err)

	// Close listener when this function returns
	defer func(listener net.Listener) {
		utils.HandleError(listener.Close())
	}(listener)

	for {
		connection, err := listener.Accept()
		utils.HandleError(err)

		if banList.IsBanned(connection.RemoteAddr().String()) {
			fmt.Printf("Rejecting request from banned client %s\n", connection.RemoteAddr())
			continue
		}

		go proxy.ProxyConnection(connection, fmt.Sprintf("%s:%d", Args.Destination, Args.DestinationPort))
	}
}

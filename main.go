package main

import (
	"flag"
	"github.com/conformal/btcrpcclient"
	"github.com/conformal/btcutil"
	"github.com/conformal/btcwire"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

var address = flag.String("address", ":8080", "Websocket listen address")
var rpcHost = flag.String("rpchost", "localhost:8334", "Hostname for RPC Server")
var rpcUser = flag.String("rpcuser", "", "Username for RPC access")
var rpcPass = flag.String("rpcpass", "", "Password for RPC access")
var rpcCert = flag.String("rpccert", "", "Certificate for RPC access")

func main() {
	flag.Parse()

	ntfnHandlers := btcrpcclient.NotificationHandlers{
		OnTxAccepted: func(hash *btcwire.ShaHash, amount btcutil.Amount) {
			log.Printf("Tx Accepted: %v (%s)", hash, amount.String())
		},
	}

    if *rpcCert != "" {
        certFile := *rpcCert
    } else {
        btcdHomeDir := btcutil.AppDataDir("btcd", false)
        filepath.Join(btcdHomeDir, "rpc.cert")
        rpcCert = filepath.Join(btcdHomeDir, "rpc.cert")
    }

	certs, _ := ioutil.ReadFile(rpcCert)

	connCfg := &btcrpcclient.ConnConfig{
		Host:         *rpcHost,
		Endpoint:     "ws",
		User:         *rpcUser,
		Pass:         *rpcPass,
		Certificates: certs,
	}

	// Connect to btcd server
	client, err := btcrpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
	}

	client.NotifyNewTransactions(false)

	// Start Websocket server
	go h.run()
	http.HandleFunc("/", serveWs)
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Error starting http listener: ", err)
	}
}

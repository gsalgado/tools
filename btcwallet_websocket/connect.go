package btcwallet_websocket

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/monetas/btcutil"
	"github.com/monetas/websocket"
)

func Connect(port int) (*websocket.Conn, error) {
	// get the root cert for connecting to secure websocket
	btcwalletHomeDir := btcutil.AppDataDir("btcwallet", false)
	certs, err := ioutil.ReadFile(filepath.Join(btcwalletHomeDir, "rpc.cert"))

	if err != nil {
		return nil, err
	}
	// Setup TLS
	var tlsConfig *tls.Config
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(certs)
	tlsConfig = &tls.Config{
		RootCAs:    pool,
		MinVersion: tls.VersionTLS12,
	}

	// The RPC server requires basic authorization, so create a custom
	// request header with the Authorization header set.
	login := "user:pass"
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(login))
	requestHeader := make(http.Header)
	requestHeader.Add("Authorization", auth)

	dialer := websocket.Dialer{TLSClientConfig: tlsConfig}
	url := fmt.Sprintf("wss://localhost:%v/ws", port)
	conn, _, err := dialer.Dial(url, requestHeader)
	if err != nil {
		return nil, err
	}
	return conn, nil
}


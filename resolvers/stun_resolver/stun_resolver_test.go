package stunresolver

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/kitsuyui/myip/resolvers/base"
	"github.com/pion/stun"
)

func TestSTUNSuccess(t *testing.T) {
	serverAddr := startTestUDPSTUNServer(t, net.ParseIP("203.0.113.20"))
	h := STUNDetector{Host: "stun:" + serverAddr, Protocol: "udp"}
	ip, err := h.RetrieveIP()
	if err != nil {
		t.Errorf("Should be succeed")
	}
	if ip == nil || !ip.IP.Equal(net.ParseIP("203.0.113.20")) || ip.Score != 0.1 {
		t.Errorf("IP must not nil")
	}
}

func TestSTUNSSuccess(t *testing.T) {
	serverName, serverAddr, pool := startTestTLSSTUNServer(t, net.ParseIP("2001:db8::7"))
	originalTLSConfigForHost := tlsConfigForHost
	tlsConfigForHost = func(host string) *tls.Config {
		cfg := originalTLSConfigForHost(host)
		cfg.RootCAs = pool
		return cfg
	}
	t.Cleanup(func() {
		tlsConfigForHost = originalTLSConfigForHost
	})

	h := STUNDetector{Host: fmt.Sprintf("stuns:%s:%s", serverName, serverAddr.Port()), Protocol: "tcp"}
	ip, err := h.RetrieveIP()
	if err != nil {
		t.Errorf("Should be succeed")
	}
	if ip == nil || !ip.IP.Equal(net.ParseIP("2001:db8::7")) || ip.Score != 1.0 {
		t.Errorf("IP must not nil")
	}
}

func TestSTUNFail(t *testing.T) {
	ctx := context.Background()
	timeout := 500 * time.Millisecond
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	h := STUNDetector{Host: "stun:127.0.0.1:1000", Protocol: "udp"}
	var err error
	var ip net.IP
	type Result struct {
		IP  net.IP
		Err error
	}
	c := make(chan Result)
	go func() {
		ip, err := h.RetrieveIP()
		c <- Result{ip.IP, err}
	}()
	select {
	case <-ctx.Done():
		err = &base.TimeoutError{}
	case r := <-c:
		ip = r.IP
		err = r.Err
	}
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}

func TestSTUNInvalidAddress(t *testing.T) {
	h := STUNDetector{Host: "<>", Protocol: "udp"}
	ip, err := h.RetrieveIP()
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}

func TestSTUNInvalidProtocol(t *testing.T) {
	h := STUNDetector{Host: "stuns:<>", Protocol: "xxx"}
	ip, err := h.RetrieveIP()
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
	h = STUNDetector{Host: "stun:<>", Protocol: "xxx"}
	ip, err = h.RetrieveIP()
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}

func TestGetString(t *testing.T) {
	result := STUNDetector{Host: "stun.l.google.com:19302", Protocol: "udp"}.String()
	tobe := "stun.l.google.com:19302"
	if result != tobe {
		t.Errorf("The result must be %s", tobe)
	}
}

func startTestUDPSTUNServer(t *testing.T, responseIP net.IP) string {
	t.Helper()

	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen udp: %v", err)
	}
	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Errorf("close udp listener: %v", err)
		}
	})

	go handleSTUNPacketConn(t, conn, responseIP)

	return conn.LocalAddr().String()
}

func startTestTLSSTUNServer(t *testing.T, responseIP net.IP) (string, testHostPort, *x509.CertPool) {
	t.Helper()

	serverName := "localhost"
	certificate, pool := mustCreateTestCertificate(t, serverName)
	listener, err := tls.Listen("tcp", net.JoinHostPort(serverName, "0"), &tls.Config{
		Certificates: []tls.Certificate{certificate},
	})
	if err != nil {
		t.Fatalf("listen tls: %v", err)
	}
	t.Cleanup(func() {
		if err := listener.Close(); err != nil {
			t.Errorf("close tls listener: %v", err)
		}
	})

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go handleSTUNTCPConn(t, conn, responseIP)
		}
	}()

	return serverName, testHostPort(listener.Addr().String()), pool
}

type testHostPort string

func (p testHostPort) Port() string {
	_, port, err := net.SplitHostPort(string(p))
	if err != nil {
		return ""
	}
	return port
}

func handleSTUNPacketConn(t *testing.T, conn net.PacketConn, responseIP net.IP) {
	t.Helper()

	buffer := make([]byte, 1500)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			return
		}
		response, err := buildSTUNSuccessResponse(buffer[:n], responseIP)
		if err != nil {
			t.Errorf("build STUN response: %v", err)
			return
		}
		if _, err := conn.WriteTo(response, addr); err != nil {
			t.Errorf("write STUN response: %v", err)
			return
		}
	}
}

func handleSTUNTCPConn(t *testing.T, conn net.Conn, responseIP net.IP) {
	t.Helper()
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Errorf("set deadline: %v", err)
		return
	}

	buffer := make([]byte, 1500)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Errorf("read STUN request: %v", err)
		return
	}

	response, err := buildSTUNSuccessResponse(buffer[:n], responseIP)
	if err != nil {
		t.Errorf("build STUN response: %v", err)
		return
	}

	if _, err := conn.Write(response); err != nil {
		t.Errorf("write STUN response: %v", err)
	}
}

func buildSTUNSuccessResponse(rawRequest []byte, responseIP net.IP) ([]byte, error) {
	var request stun.Message
	request.Raw = append([]byte(nil), rawRequest...)
	if err := request.Decode(); err != nil {
		return nil, err
	}

	response := stun.MustBuild(
		stun.NewTransactionIDSetter(request.TransactionID),
		stun.BindingSuccess,
		&stun.XORMappedAddress{IP: responseIP, Port: 3478},
	)
	return response.Raw, nil
}

func mustCreateTestCertificate(t *testing.T, serverName string) (tls.Certificate, *x509.CertPool) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate private key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: serverName,
		},
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(time.Hour),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		DNSNames: []string{serverName},
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}

	certificatePEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	certificate, err := tls.X509KeyPair(certificatePEM, keyPEM)
	if err != nil {
		t.Fatalf("parse certificate: %v", err)
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(certificatePEM)
	return certificate, pool
}

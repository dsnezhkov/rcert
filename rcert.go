package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/CrowdSurge/banner" // for banner
	"github.com/fatih/color"       // for colors
)

// Flag parsing setup
type stringFlag struct {
	set   bool
	value string
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

var ipFile stringFlag

// Channels
var connUnverifiedTLSStatus = make(chan bool)
var connUnverifiedTLSData = make(chan []*x509.Certificate)

// Colors
var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var blue = color.New(color.FgBlue).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()
var magenta = color.New(color.FgMagenta).SprintFunc()

func DialTLS(network, addr string) (net.Conn, error) {
	conn, err := tls.Dial(network, addr, &tls.Config{
		InsecureSkipVerify: true,
	})

	if err == nil {
		cs := conn.ConnectionState()
		// Send certificates info back on channel
		connUnverifiedTLSStatus <- true
		connUnverifiedTLSData <- cs.PeerCertificates
	} else {
		connUnverifiedTLSStatus <- false
	}
	return conn, err
}
func rcert(h string, client *http.Client) {
	client.Get("https://" + h)
	/* resp, err = client.Get("https://" + h)
	fmt.Println(resp, err) */

	if connstatus := <-connUnverifiedTLSStatus; connstatus == true {
		fmt.Printf("%s\n", green("connected"))
		pcerts := <-connUnverifiedTLSData
		for _, c := range pcerts {
			fmt.Printf("%s: %s: %s\n", h, yellow("Issuer"), c.Issuer)
			fmt.Printf("%s: %s: %s\n", h, yellow("Subject"), c.Subject)
			for _, dn := range c.DNSNames {
				fmt.Printf("%s: %s: %s\n", h, blue("Ext: SAN DNS name"), dn)
			}
			for _, pdn := range c.PermittedDNSDomains {
				fmt.Printf("%s: %s: %s\n", h, magenta("Ext: SAN Permitted domain"), pdn)
			}
			for _, en := range c.EmailAddresses {
				fmt.Printf("%s: %s: %s\n", h, cyan("Ext: SAN Email"), en)
			}
		}
	} else {
		fmt.Printf("%s\n", red("not connected"))
	}
}

func usage() {
	fmt.Println("rcert -ipfile=/path/to/file")
	os.Exit(2)
}

func init() {
	flag.Var(&ipFile, "ipfile", "filename containing IPs, one per line")
}
func main() {

	flag.Parse()
	if !ipFile.set {
		fmt.Println("-ipfile not set")
		usage()
	}

	//timeout := time.Duration(10 * time.Second)
	//tlstimeout := time.Duration(5 * time.Second)
	timeout := time.Duration(6 * time.Second)
	tlstimeout := time.Duration(3 * time.Second)
	hosts := []string{}

	// Set color banner
	color.Set(color.FgRed)
	banner.Print("r-cert")
	fmt.Println("___________________________________")
	fmt.Println()
	color.Unset()

	// Read list of IP addresses from file
	fmt.Println("[*] Populating list of IPs from file")
	f, err := os.Open(ipFile.value)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		hosts = append(hosts, line)
	}
	f.Close()

	// Configure TLS client (untrusted certs included)
	fmt.Println("[*] Configuring TLS transport client")
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{

			DialTLS:             DialTLS,
			TLSHandshakeTimeout: tlstimeout,
		},
	}

	fmt.Println("[*] R-certing the IPs")
	fmt.Println("___________________________________")
	fmt.Println()

	for _, h := range hosts {
		fmt.Printf("%s: ", h)
		rcert(h, client)
	}
}

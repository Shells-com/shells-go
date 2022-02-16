package main

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/KarpelesLab/rest"
	"github.com/Shells-com/shells-go/shellsui"
	"github.com/Shells-com/shells-go/spicefyne"
)

type shellHost struct {
	ID   string `json:"Shell_Host__"`
	Name string
	IP   string
	IPv6 string
}

type shellsOS struct {
	ID     string `json:"Shell_OS__"`
	Code   string
	URL    string
	Family string // "linux", etc
	Boot   string
	Name   string
}

type shellSpice struct {
	Host     string `json:"host"`     // "la01-01.shellsnet.com"
	ID       string `json:"id"`       // Shell ID (shell-***)
	Key      string `json:"key"`      // socket key
	Password string `json:"password"` // spice password
	Port     int    `json:"port"`     // 443
	Protocol string `json:"protocol"` // "wss"
	Token    string `json:"token"`    // "shell-asqn2r-ucp5-c2fh-heow-ronwydbe.dQQsb3X9UHt0ZFyCqZnJrHPQniDW3TMZ"
	URL      string `json:"url"`      // "wss://..."
}

type shell struct {
	ID              string `json:"Shell__"`
	Label           string
	Engine          string // "full"
	Size            int    // 16
	Status          string // "valid"
	State           string // "running"
	Ssh_Port        int    // 12345
	Username        string
	Hostname        string
	MAC             string
	IPv4            string
	IPv6            string
	Created         rest.Time
	Expires         rest.Time
	Last_Snapshot   rest.Time
	Timer_Allowance int
	Host            *shellHost
	OS              *shellsOS

	w     fyne.Window
	a     fyne.App
	token *rest.Token
	spice shellSpice
}

func shellsList(a fyne.App, w fyne.Window, token *rest.Token) {
	// try to use token
	var shells []shell
	p := map[string]interface{}{
		"Status": "valid",
	}
	err := rest.Apply(token.Use(context.Background()), "Shell", "GET", p, &shells)
	if err != nil {
		log.Printf("failed to get: %s", err)
		a.Preferences().RemoveValue("token")
		loginWindow(a, w, shellsList)
		return
	} else {
		//log.Printf("got: %+v", shells)
	}

	var fields []fyne.CanvasObject
	title := widget.NewLabel("Choose:")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}
	fields = append(fields, title)

	for _, shl := range shells {
		shlCopy := shl
		btn := widget.NewButton(shl.Label, func() {
			shlCopy.w = w
			shlCopy.a = a
			shlCopy.token = token
			shlCopy.run()
		}) // TODO
		fields = append(fields, btn)
	}

	logout := widget.NewButton("Logout", func() {
		a.Preferences().SetString("token", "")
		a.Preferences().RemoveValue("token")
		loginWindow(a, w, shellsList)
	})
	logout.Importance = widget.LowImportance
	fields = append(fields, logout)
	w.SetContent(shellsui.GetMainContainer(container.NewVBox(fields...)))
}

func (s *shell) ProxyConnect(port uint16, compress bool, timeout time.Duration) (net.Conn, error) {
	// attempt to connect to a given port
	cfg := &tls.Config{
		NextProtos: []string{"shl-proxy"},
	}
	c, err := s.tlsDial("tcp", fmt.Sprintf("%s:%d", s.spice.Host, s.spice.Port), cfg)
	if err != nil {
		return nil, err
	}

	// check negociated protocol
	if c.ConnectionState().NegotiatedProtocol != "shl-proxy" {
		c.Close()
		return nil, errors.New("failed to establish shl-proxy")
	}

	var flags uint8

	if compress {
		flags |= 1 // compress flag
	}
	flags |= 2 // get connection established confirmation

	// send token
	tok := append([]byte{flags, byte((port >> 8) & 0xff), byte(port & 0xff)}, s.spice.Token...)
	ln := make([]byte, 2)
	binary.BigEndian.PutUint16(ln, uint16(len(tok)))

	c.Write(ln)
	c.Write(tok)

	// read result
	c.SetReadDeadline(time.Now().Add(timeout))
	res := make([]byte, 1)
	_, err = c.Read(res)
	if err != nil {
		return nil, err
	}
	if res[0] != 0 {
		return nil, errors.New("Connection has failed")
	}
	c.SetReadDeadline(time.Time{})

	// c is either ready to use, or closed.
	if compress {
		return &snappyConn{c: c}, nil
	}
	return c, nil
}

func (s *shell) SpiceConnect(compress bool) (net.Conn, error) {
	// connect to server
	cfg := &tls.Config{
		NextProtos: []string{"shl-spice"},
	}
	c, err := s.tlsDial("tcp", fmt.Sprintf("%s:%d", s.spice.Host, s.spice.Port), cfg)
	if err != nil {
		return nil, err
	}

	// check negociated protocol
	if c.ConnectionState().NegotiatedProtocol != "shl-spice" {
		c.Close()
		return nil, errors.New("failed to establish shl-spice")
	}

	var flags uint8

	if compress {
		flags |= 1 // compress flag
	}

	// send token
	tok := append([]byte{flags}, s.spice.Token...)
	ln := make([]byte, 2)
	binary.BigEndian.PutUint16(ln, uint16(len(tok)))

	c.Write(ln)
	c.Write(tok)

	// c is either ready to use, or closed.
	if compress {
		return &snappyConn{c: c}, nil
	}
	return c, nil
}

func (s *shell) doProxyPort(c net.Conn, port int, compress bool) {
	defer c.Close()
	// attempt to connect to port
	rc, err := s.ProxyConnect(uint16(port), compress, 30*time.Second)
	if err != nil {
		// failed :(
		return
	}

	// interconnect
	interco(rc, c, compress)
}

func (s *shell) startProxyPort(port int, compress bool) (*net.TCPListener, error) {
	// start a proxy server on a random local port
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		return nil, err
	}

	go func() {
		defer l.Close() // probably not needed
		for {
			c, err := l.AcceptTCP()
			if err != nil {
				return
			}
			go s.doProxyPort(c, port, compress)
		}
	}()

	return l, err
}

func (s *shell) tryRdp() {
	// let's try to connect to rdp port (3389)
	// first, let's see if xfreerdp exists (TODO support more options for rdp clients)
	p, err := exec.LookPath("xfreerdp")
	if err != nil {
		log.Printf("rdp: could not locate xfreerdp, not trying")
		return
	}

	log.Printf("rdp: attempting connection on port 3389 ...")
	// attempt to connect, if it works let's move to the actual running of xfreerdp
	c, err := s.ProxyConnect(3389, false, 2*time.Second)
	if err != nil {
		log.Printf("rdp: failed to connect to port 3389, giving up")
		return
	}
	c.Close()

	l, err := s.startProxyPort(3389, false)
	defer l.Close()

	addr := l.Addr().(*net.TCPAddr)
	port := addr.Port

	log.Printf("rdp: proxying through local port %d, running rdp client...", port)

	// run xfreerdp
	// /v:192.168.10.30 /w:1600 /h:1200
	// /f = fullscreen
	// /sound:sys:alsa
	// /microphone:sys:alsa
	// /multimedia:sys:alsa
	// /usb:id,dev:054c:0268
	cmd := exec.Command(p, fmt.Sprintf("/v:127.0.0.1:%d", port), "/cert:ignore", "/u:Administrator", "/p:")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("about to run: %s %s", cmd.Path, cmd.Args)
	err = cmd.Run()
	if err != nil {
		log.Printf("failed to run xfreerdp: %s", err)
		return
	}
	os.Exit(0)
}

func (s *shell) run() {
	// loading...
	s.w.SetContent(shellsui.GetMainContainer(container.NewVBox(widget.NewProgressBarInfinite())))

	// need to get spice access
	err := rest.Apply(s.token.Use(context.Background()), "Shell/"+s.ID+":spice", "POST", map[string]interface{}{}, &s.spice)
	if err != nil {
		log.Printf("failed: %s", err)
		return
	}

	log.Printf("got spice = %+v", s.spice)

	s.tryRdp()

	_, err = spicefyne.New(s.w, s.a, s, s.spice.Password)
	if err != nil {
		log.Printf("spice init failed: %s", err)
		os.Exit(1)
	}
}

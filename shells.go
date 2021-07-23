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

func (s *shell) SpiceConnect(compress bool) (net.Conn, error) {
	// connect to server
	cfg := &tls.Config{
		NextProtos: []string{"shl-spice"},
	}
	c, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.spice.Host, s.spice.Port), cfg)
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

func (s *shell) run() {
	s.w.SetContent(shellsui.GetMainContainer(container.NewVBox(widget.NewProgressBarInfinite())))

	// need to get spice access
	err := rest.Apply(s.token.Use(context.Background()), "Shell/"+s.ID+":spice", "POST", map[string]interface{}{}, &s.spice)
	if err != nil {
		log.Printf("failed: %s", err)
		return
	}

	log.Printf("got spice = %+v", s.spice)

	_, err = spicefyne.New(s.w, s.a, s, s.spice.Password)
	if err != nil {
		log.Printf("spice init failed: %s", err)
		os.Exit(1)
	}
}

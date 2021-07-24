package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/KarpelesLab/csscolor"
	"github.com/KarpelesLab/rest"
	"github.com/Shells-com/shells-go/res"
	"github.com/Shells-com/shells-go/shellsui"
	"github.com/pkg/browser"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/vincent-petithory/dataurl"
)

const clientID = "oaap-yg3n7j-nxzj-d5zm-eulq-oot76jli"

type loginShells struct {
	a       fyne.App
	w       fyne.Window
	session string
	vals    map[string]interface{}
	fields  []fyne.CanvasObject
	next    func(fyne.App, fyne.Window, *rest.Token)
}

type loginOauth2 struct {
	ID        string `json:"OAuth2_Consumer__"`
	Name      string
	TokenName string `json:"Token_Name"`
}

type loginField struct {
	Cat    string       `json:"cat"`
	Type   string       `json:"type"`
	Label  string       `json:"label"`
	Name   string       `json:"name"`
	Button *loginButton `json:"button"` // "button":{"background-color":"#1da1f2","logo":"
	Info   *loginOauth2 `json:"info"`
}

type loginButton struct {
	BackgroundColor string `json:"background-color"` // #ffffff #000 #1da1f2 etc...
	Logo            string `json:"logo"`             // data:image/svg+xml;base64,... or https://...
}

type loginRes struct {
	Complete bool `json:"complete"`
	Fields   []*loginField
	Initial  bool                   `json:"initial"` // if true, no "reset" button
	Message  string                 `json:"message"`
	Req      []string               `json:"req"`
	Session  string                 `json:"session"`
	User     map[string]interface{} `json:"user"`
	Token    *rest.Token            // bearer token, 1h lifetime
	Url      string                 `json:"url"`
}

var email string = ""

func loginWindow(a fyne.App, w fyne.Window, next func(fyne.App, fyne.Window, *rest.Token)) {
	w.Resize(fyne.Size{Width: WWidth, Height: WHeight})

	l := &loginShells{a: a, w: w, next: next}
	l.loading()
	w.Show()

	// if we have a saved session, do not go further
	if l.checkSession() {
		return
	}

	switch os.Getenv("SHELLS_LOGIN") {
	case "thin":
		// special mode, generate a QR login and ask user to login
		go func() {
			err := l.doThin()
			if err != nil {
				l.doErr(err)
			}
		}()
	case "flow":
		fallthrough
	default:
		go func() {
			time.Sleep(50 * time.Millisecond)
			l.call(nil)
		}()
	}
}

func (l *loginShells) doSubmit() {
	// check entry, etc
	l.call(l.vals)
}

func (l *loginShells) checkSession() bool {
	// check preferences for a token
	if jsB := l.a.Preferences().String("token"); jsB != "" {
		log.Printf("login: found saved token")
		if js, err := base64.RawURLEncoding.DecodeString(jsB); err == nil {
			var token *rest.Token
			if err := json.Unmarshal(js, &token); err == nil {
				token.ClientID = clientID
				go l.next(l.a, l.w, token)
				return true
			} else {
				log.Printf("json parse failed: %s", err)
			}
		} else {
			log.Printf("base64 parse failed: %s", err)
		}
	}
	return false
}

func (l *loginShells) call(vars map[string]interface{}) {
	req := make(map[string]interface{})
	if vars != nil {
		// duplicate values
		for k, v := range vars {
			req[k] = v
		}
	}
	req["client_id"] = clientID
	if l.session != "" {
		req["session"] = l.session
	}

	l.loading()
	go l.doCall(req)
}

func (l *loginShells) doErr(err error) {
	if err == nil {
		return
	}

	var msg string
	switch e := err.(type) {
	case *rest.Error:
		msg = e.Response.Error
	default:
		msg = err.Error()
	}

	d := dialog.NewError(errors.New(msg), l.w)

	d.SetOnClosed(func() {
		// restore focus
		l.restoreFocus()
	})

	l.w.SetContent(
		container.New(
			layout.NewBorderLayout(
				shellsui.GetMarginRectangle(),
				shellsui.GetMarginRectangle(),
				shellsui.GetMarginRectangle(),
				shellsui.GetMarginRectangle()),
			container.NewVBox(l.fields...),
		),
	)
	l.restoreFocus()
}

func (l *loginShells) doThin() error {
	// special process
	var res map[string]interface{}
	err := rest.Apply(context.Background(), "OAuth2/App/"+clientID+":token_create", "POST", map[string]interface{}{}, &res)
	if err != nil {
		return err
	}
	tok, ok := res["polltoken"].(string)
	if !ok {
		return fmt.Errorf("failed to fetch polltoken")
	}

	// see: https://www.shells.com/.well-known/openid-configuration?pretty
	tokuri := url.QueryEscape("polltoken:" + tok)
	fulluri := fmt.Sprintf("https://www.shells.com/_rest/OAuth2:auth?response_type=code&client_id=%s&redirect_uri=%s&scope=profile", clientID, tokuri)

	if u, ok := res["xox"].(string); ok {
		fulluri = u
	}

	qr, err := qrcode.New(fulluri, qrcode.Medium)
	if err != nil {
		return err
	}
	img := qr.Image(256)
	cimg := canvas.NewImageFromImage(img)
	cimg.FillMode = canvas.ImageFillOriginal

	//log.Printf("Please open this URL in order to access shells:\n%s", fulluri)
	var fields []fyne.CanvasObject

	fields = append(fields, shellsui.GetShellsLabel("Please authenticate this device"))
	fields = append(fields, cimg)
	fields = append(fields, shellsui.GetShellsLabel("Peering code: "+strings.TrimPrefix(fulluri, "xox.jp/")))

	l.w.SetContent(shellsui.GetMainContainer(container.NewVBox(fields...)))
	l.w.CenterOnScreen()

	// wait for login to complete
	for {
		var res map[string]interface{}
		err := rest.Apply(context.Background(), "OAuth2/App/"+clientID+":token_poll", "POST", map[string]interface{}{"polltoken": tok}, &res)
		if err != nil {
			return err
		}

		v, ok := res["response"]
		if !ok {
			time.Sleep(time.Second) // just in case
			continue
		}

		l.loading()

		resp, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid response from api, response of invalid type")
		}

		code, ok := resp["code"].(string)
		if !ok {
			return fmt.Errorf("invalid response from api, response not containing code")
		}

		log.Printf("fetching auth token...")

		// https://www.shells.com/_special/rest/OAuth2:token
		httpresp, err := http.PostForm("https://www.shells.com/_special/rest/OAuth2:token", url.Values{"client_id": {clientID}, "grant_type": {"authorization_code"}, "code": {code}})
		if err != nil {
			return fmt.Errorf("while fetching token: %w", err)
		}
		defer httpresp.Body.Close()

		if httpresp.StatusCode != 200 {
			return fmt.Errorf("invalid status code from server: %s", httpresp.Status)
		}

		body, err := io.ReadAll(httpresp.Body)
		if err != nil {
			return fmt.Errorf("while reading token: %w", err)
		}

		// decode token
		var token *rest.Token
		err = json.Unmarshal(body, &token)
		if err != nil {
			return fmt.Errorf("while decoding token: %w", err)
		}
		token.ClientID = clientID

		if js, err := json.Marshal(token); err == nil {
			// should be the case
			jsB := base64.RawURLEncoding.EncodeToString(js)
			l.a.Preferences().SetString("token", jsB)
		}

		l.next(l.a, l.w, token)

		return nil
	}
}

func (l *loginShells) doCall(vars map[string]interface{}) {
	var result loginRes
	err := rest.Apply(context.Background(), "User:flow", "POST", vars, &result)
	if err != nil {
		l.doErr(err)
		return
	}

	if result.Complete {
		log.Printf("login operation has completed")

		result.Token.ClientID = clientID

		if js, err := json.Marshal(result.Token); err == nil {
			// should be the case
			jsB := base64.RawURLEncoding.EncodeToString(js)
			l.a.Preferences().SetString("token", jsB)
		}

		l.next(l.a, l.w, result.Token)
		return
	}

	//log.Printf("got response: %+v", res)
	l.session = result.Session

	var fields []fyne.CanvasObject
	var focus fyne.Focusable

	vals := make(map[string]interface{})
	var vLock sync.Mutex
	socialbtnsC := container.NewHBox()

	for _, f := range result.Fields {
		switch f.Type {
		case "text", "email", "password", "phone":
			e := widget.NewEntry()
			switch f.Type {
			case "email":
				l1 := shellsui.GetShellsLabel("Please enter your Shells Account")
				l2 := shellsui.GetShellsLabel("email address to log in")

				fields = append(fields, l1)
				fields = append(fields, l2)
			case "password":
				wb := shellsui.GetShellsLabel("Welcome Back")
				emailStr := shellsui.GetShellsLabel(email)
				pwd := shellsui.GetShellsLabel("Please enter your password")

				fields = append(fields, wb)
				fields = append(fields, emailStr)
				fields = append(fields, shellsui.GetMarginRectangleWithHeight(30))
				fields = append(fields, pwd)
				e.Password = true
			}

			name := f.Name
			fType := f.Type
			e.OnChanged = func(s string) {
				vLock.Lock()
				defer vLock.Unlock()

				vals[name] = s
				if fType == "email" {
					email = s
				}
			}
			e.OnSubmitted = func(value string) {
				l.call(l.vals)
			}

			ctn := container.New(
				layout.NewBorderLayout(
					shellsui.GetMarginRectangleWithHeight(15),
					shellsui.GetMarginRectangleWithHeight(30),
					nil,
					nil,
				),
				e,
			)
			// TODO e.validator
			fields = append(fields, ctn)
			if focus == nil {
				focus = e
			}
		case "checkbox":
			name := f.Name
			e := widget.NewCheck(f.Label, func(v bool) {
				vLock.Lock()
				defer vLock.Unlock()

				vals[name] = v
			})
			fields = append(fields, e)
		case "oauth2":
			var icon fyne.Resource

			//icon := getOAuth2Icon(f.Info.TokenName)
			if f.Button != nil {
				logo := f.Button.Logo
				if logo != "" {
					if strings.HasPrefix(logo, "data:") {
						// parse as data uri
						d, err := dataurl.DecodeString(logo)
						if err == nil {
							icon = fyne.NewStaticResource(fmt.Sprintf("logo-%s.svg", f.Info.ID), d.Data)
						}
					} else {
						// url, download it
						icon, _ = fyne.LoadResourceFromURLString(logo)
					}
				}
			}

			id := f.Info.ID
			bgColor, err := csscolor.Parse(f.Button.BackgroundColor)
			if err != nil {
				bgColor = theme.PrimaryColor()
			}

			btn := shellsui.NewSocialButton(icon, bgColor, func() {
				l.oauth2(id)
			})
			socialbtnsC.Add(btn)

		default:
			log.Printf("unknown field = %+v", f)
		}
	}

	if len(socialbtnsC.Objects) > 0 {
		fields = append(fields, container.NewCenter(socialbtnsC))
	}

	submit := widget.NewButtonWithIcon("", res.RightArrow, l.doSubmit)

	submit.Importance = widget.HighImportance
	fields = append(fields, container.NewMax(submit))

	l.vals = vals
	l.fields = fields

	l.w.SetContent(shellsui.GetMainContainer(container.NewVBox(fields...)))
	l.w.CenterOnScreen()

	if focus != nil {
		l.w.Canvas().Focus(focus)
	}
}

func (l *loginShells) loading() {
	l.w.SetContent(shellsui.GetMainContainer(container.NewVBox(widget.NewProgressBarInfinite())))
}

func (l *loginShells) restoreFocus() {
	for _, w := range l.fields {
		if f, ok := w.(fyne.Focusable); ok {
			l.w.Canvas().Focus(f)
			break
		}
	}
}

func (l *loginShells) oauth2(id string) {
	l.loading()

	var result map[string]interface{}
	err := rest.Apply(context.Background(), "OAuth2/App/"+clientID+":token_create", "POST", map[string]interface{}{}, &result)

	if err != nil {
		log.Printf("failed to fetch the token: %s", err)
		l.doErr(err)
		return
	}

	tok, ok := result["polltoken"].(string)
	if !ok {
		log.Printf("failed to fetch polltoken")
		l.doErr(errors.New("invalid response from API"))
		return
	}
	tokuri := "polltoken:" + tok

	req := make(map[string]interface{})
	if l.vals != nil {
		// duplicate values
		for k, v := range l.vals {
			req[k] = v
		}
	}
	req["client_id"] = clientID
	if l.session != "" {
		req["session"] = l.session
	}
	req["oauth2"] = id
	req["redirect_uri"] = tokuri

	var resultUserFlow loginRes
	err = rest.Apply(context.Background(), "User:flow", "POST", req, &resultUserFlow)

	if err != nil {
		log.Printf("failed to update the User flow with oauth2: %s", err)
		l.doErr(err)
		return
	}
	if resultUserFlow.Url == "" {
		log.Printf("failed to get the url from the payload from User:flow")
		l.doErr(errors.New("missing login URL"))
		return
	}
	err = browser.OpenURL(resultUserFlow.Url)
	if err != nil {
		l.doErr(err)
		return
	}

	for {
		var res map[string]interface{}
		err := rest.Apply(context.Background(), "OAuth2/App/"+clientID+":token_poll", "POST", map[string]interface{}{"polltoken": tok}, &res)
		if err != nil {
			log.Printf("failed to fetch the token poll: %s", err)
			l.doErr(err)
			return
		}

		log.Printf("api response = %+v", res)

		v, ok := res["response"]
		if !ok {
			time.Sleep(time.Second) // just in case
			continue
		}

		resp, ok := v.(map[string]interface{})
		if !ok {
			log.Printf("invalid response from api, response of invalid type")
			l.doErr(errors.New("invalid response from API"))
			return
		}

		session, ok := resp["session"].(string)
		if !ok {
			log.Printf("invalid response from api, response not containing session")
			l.doErr(errors.New("invalid response from API"))
			return
		}
		l.session = session
		l.call(nil)
		return
	}
}

package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/Shells-com/shells-go/res"

	"fyne.io/fyne/v2/app"
	"github.com/KarpelesLab/rest"
	"github.com/gordonklaus/portaudio"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	log.Printf("Initializing PortAudio version %s", portaudio.VersionText())
	portaudio.Initialize()
	defer portaudio.Terminate()

	if dev, err := portaudio.DefaultOutputDevice(); err == nil {
		log.Printf("default output device: %s type %s", dev.Name, dev.HostApi.Name)
	}

	rest.Host = "www.shells.com"
	rest.Debug = true

	os.Setenv("FYNE_SCALE", "1")

	a := app.NewWithID("com.shells.app")
	a.Settings().SetTheme(getShellsTheme())
	w := a.NewWindow("Shells")
	w.SetIcon(res.ShellsIcon)

	loginWindow(a, w, shellsList)

	a.Run()
}

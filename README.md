# Shells™ Go client

This is a client to remotely access Shells™ virtual desktop instances.

It uses Spice protocol in order to provide access to the display, controls, etc.

## PortAudio / opus dependencies

Add the following to `go.mod` for statically including portaudio & opus as static libs:

    replace github.com/gordonklaus/portaudio => github.com/KarpelesLab/static-portaudio v0.6.190600
    replace github.com/hraban/opus => github.com/KarpelesLab/static-opus v0.5.131

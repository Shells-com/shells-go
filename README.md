# Shells™ Go Client

A cross-platform desktop client for remotely accessing [Shells™](https://www.shells.com) virtual desktop instances.

## Requirements

**A Shells™ account is required to use this application.** You can create an account at [shells.com](https://www.shells.com).

## Features

- Remote desktop access to your Shells™ virtual machines
- Cross-platform support (Linux, macOS, Windows)
- SPICE protocol for high-performance display, keyboard, and mouse input
- Audio support via PortAudio and Opus codec
- Clipboard sharing between local and remote systems
- OAuth2 authentication with QR code support for thin clients

## Building

### Prerequisites

- Go 1.24 or later
- CGO enabled (required for audio and graphics dependencies)

### Linux

Install the required development packages:

```bash
# Debian/Ubuntu
sudo apt-get install libgl1-mesa-dev libxcursor-dev libxrandr-dev \
    libxinerama-dev libxi-dev libxxf86vm-dev libasound2-dev pkg-config

# Fedora
sudo dnf install mesa-libGL-devel libXcursor-devel libXrandr-devel \
    libXinerama-devel libXi-devel libXxf86vm-devel alsa-lib-devel pkgconfig
```

### macOS

Xcode command line tools are required:

```bash
xcode-select --install
```

### Windows

MSYS2 with MinGW64 toolchain is recommended for building on Windows.

### Build

```bash
go build -v .
```

## Static Linking (PortAudio / Opus)

For producing portable binaries with statically linked audio libraries, add the following to `go.mod`:

```
replace github.com/gordonklaus/portaudio => github.com/KarpelesLab/static-portaudio v0.6.190600
replace github.com/hraban/opus => github.com/KarpelesLab/static-opus v0.6.152
```

## Environment Variables

- `SHELLS_FULLSCREEN=WIDTHxHEIGHT` - Start in fullscreen mode with specified resolution (e.g., `1920x1080`)
- `SHELLS_LOGIN=thin` - Use QR code login flow for thin client deployments

## License

This software is proprietary to Shells™. See [shells.com](https://www.shells.com) for terms of service.

# GoPortInspector

## Overview
GoPortInspector is a lightweight utility that periodically inspects the network ports open on your machine and shows the associated process for each port. It provides a quick snapshot of active TCP and UDP connections, making it useful for troubleshooting or monitoring purposes.

## Installation
### Prerequisites
- Go 1.24 or later
- Windows operating system (tested on Windows 10/11)

### Build
```
git clone https://example.com/GoPortInspector.git
cd GoPortInspector
go build
```

## Usage
Run the tool directly with `go run` or execute the compiled binary:
```
# run without building
go run portwatch.go

# or run the built binary
./portwatch
```
The output refreshes every few seconds, displaying protocol, local/remote addresses, connection state, and the owning process.

## Project Goals
- Provide an easy way to observe which processes are listening on or using network ports
- Offer a starting point for more advanced port‑monitoring tools
- Encourage contributions for cross‑platform support and additional features

## License
GoPortInspector is released under the [MIT License](LICENSE).

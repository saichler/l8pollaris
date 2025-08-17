# L8 Pollaris

[![Go Version](https://img.shields.io/badge/go-1.23.8-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

L8 Pollaris is a **Polling/Parsing & Populating model service** for agnostic collection, parsing & populating abstract models. It provides a flexible and extensible framework for defining device polling configurations and managing collection operations across different protocols and device types.

## Overview

L8 Pollaris is part of the Layer8 ecosystem, providing a centralized service for:

- **Device Polling Configuration**: Define how and what to poll from network devices
- **Protocol Abstraction**: Support for multiple protocols (SNMP, SSH, RESTCONF, NETCONF, gRPC, K8s)
- **Flexible Operations**: Support for different polling operations (GET, MAP, TABLE)
- **Group Management**: Organize devices and polling configurations into logical groups
- **Distributed Caching**: Built-in distributed caching for high performance

## Architecture

The project is structured into several key components:

### Core Components

- **PollarisService**: Main service implementation handling CRUD operations
- **PollarisCenter**: Central management hub for polling configurations and caching
- **Protocol Buffer Definitions**: Type-safe message definitions for all entities

### Protocol Buffer Types

#### Pollaris Configuration (`pollaris.proto`)
```protobuf
message Pollaris {
  string name = 1;
  string vendor = 2;
  string series = 3;
  string family = 4;
  string software = 5;
  string hardware = 6;
  string version = 7;
  repeated string groups = 8;
  map<string, Poll> polling = 9;
}
```

#### Device Configuration (`devices.proto`)
```protobuf
message Device {
  string device_id = 1;
  DeviceServiceInfo collect_service = 2;
  DeviceServiceInfo parsing_service = 3;
  DeviceServiceInfo inventory_service = 4;
  map<string, Host> hosts = 5;
}
```

#### Collection Jobs (`collect.proto`)
```protobuf
message CJob {
  string error = 1;
  bytes result = 2;
  int64 started = 3;
  int64 ended = 4;
  // ... additional fields for job management
}
```

## Supported Protocols

L8 Pollaris supports multiple network management protocols:

- **SNMP v2/v3**: Traditional SNMP polling
- **SSH**: Command-line interface polling
- **RESTCONF**: REST-based configuration protocol
- **NETCONF**: Network Configuration Protocol
- **gRPC**: Modern RPC framework
- **Kubernetes**: Container orchestration platform polling

## Supported Operations

- **OGet**: Simple get operations for single values
- **OMap**: Key-value mapping operations
- **OTable**: Tabular data collection

## Getting Started

### Prerequisites

- Go 1.23.8 or later
- Docker (for Protocol Buffer generation)
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/saichler/l8pollaris.git
cd l8pollaris
```

2. Navigate to the Go module:
```bash
cd go
```

3. Install dependencies:
```bash
go mod tidy
go mod vendor
```

### Building

```bash
go build ./...
```

### Testing

Run the test suite with coverage:
```bash
./test.sh
```

This will:
- Clean and reinitialize the Go module
- Fetch dependencies
- Run unit tests with coverage
- Generate and open a coverage report

### Protocol Buffer Generation

To regenerate Protocol Buffer bindings:

```bash
cd proto
./make-bindings.sh
```

This uses Docker to run the Protocol Buffer compiler and generates Go bindings in the `go/types/` directory.

## Configuration

### Creating a Pollaris Configuration

```go
import (
    "github.com/saichler/l8pollaris/go/pollaris"
    "github.com/saichler/l8pollaris/go/types"
)

// Create a new polling configuration
pollarisConfig := &types.Pollaris{
    Name:     "cisco-ios-xe",
    Vendor:   "cisco",
    Series:   "catalyst",
    Family:   "9000",
    Software: "ios-xe",
    Version:  "17.3",
    Groups:   []string{"switches", "enterprise"},
    Polling: map[string]*types.Poll{
        "interface-stats": {
            Name:      "interface-stats",
            What:      "1.3.6.1.2.1.2.2.1",
            Operation: types.Operation_OTable,
            Protocol:  types.Protocol_PSNMPV2,
            Cadence:   60000, // 60 seconds
            Timeout:   5000,  // 5 seconds
        },
    },
}
```

### Service Integration

```go
// Register and activate the Pollaris service
resources.Registry().Register(&types.Pollaris{})
resources.Services().Activate(
    pollaris.ServiceType, 
    pollaris.ServiceName, 
    0, 
    resources, 
    listener,
)

// Access the Pollaris center
center := pollaris.Pollaris(resources)
err := center.Add(pollarisConfig, false)
```

## API Reference

### Core Methods

#### PollarisCenter Methods

- `Add(pollaris *types.Pollaris, isNotification bool) error`
- `Update(pollaris *types.Pollaris, isNotification bool) error`
- `PollarisByName(name string) *types.Pollaris`
- `PollarisByKey(args ...string) *types.Pollaris`
- `PollsByGroup(groupName, vendor, series, family, software, hardware, version string) []*types.Pollaris`

#### Utility Functions

- `Poll(pollarisName, pollName string, resources ifs.IResources) (*types.Poll, error)`
- `PollarisByKey(resources ifs.IResources, args ...string) (*types.Pollaris, error)`
- `PollarisByGroup(resources ifs.IResources, groupName, vendor, series, family, software, hardware, version string) ([]*types.Pollaris, error)`

## Dependencies

L8 Pollaris integrates with several Layer8 components:

- **l8services**: Distributed services framework
- **l8types**: Common interfaces and types
- **l8utils**: Utility libraries
- **l8collector**: Data collection interfaces
- **l8parser**: Data parsing capabilities
- **l8srlz**: Serialization framework
- **l8test**: Testing infrastructure

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Add tests for new functionality
5. Run the test suite: `./test.sh`
6. Commit your changes: `git commit -am 'Add feature'`
7. Push to the branch: `git push origin feature-name`
8. Submit a pull request

## Testing

The project includes comprehensive unit tests. Key test areas:

- Pollaris configuration management
- Service lifecycle operations
- Group-based polling retrieval
- Key-based lookup functionality

Run tests with:
```bash
go test -v ./tests/...
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Project Structure

```
l8pollaris/
├── LICENSE                 # Apache 2.0 license
├── README.md              # This file
├── go/                    # Go source code
│   ├── go.mod            # Go module definition
│   ├── go.sum            # Go module checksums
│   ├── test.sh           # Test runner script
│   ├── pollaris/         # Main service implementation
│   │   ├── PollarisService.go    # Service interface implementation
│   │   ├── PollarisCenter.go     # Central management logic
│   │   └── PollarisUtils.go      # Utility functions
│   ├── tests/            # Unit tests
│   │   ├── Pollaris_test.go      # Main test suite
│   │   └── TestInit.go           # Test initialization
│   ├── types/            # Generated Protocol Buffer types
│   │   ├── pollaris.pb.go        # Pollaris message types
│   │   ├── devices.pb.go         # Device message types
│   │   └── collect.pb.go         # Collection job types
│   └── vendor/           # Vendored dependencies
└── proto/                # Protocol Buffer definitions
    ├── pollaris.proto    # Core polling configuration
    ├── devices.proto     # Device and connection types
    ├── collect.proto     # Collection job types
    └── make-bindings.sh  # Protocol Buffer generation script
```

## Support

For questions, issues, or contributions, please use the GitHub issue tracker or contact the Layer8 development team.
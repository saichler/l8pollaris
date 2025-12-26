# Layer 8 Data Mining Models (L8 Pollaris)

[![Go Version](https://img.shields.io/badge/go-1.23.8-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![Coverage](https://img.shields.io/badge/coverage-95%25-brightgreen.svg)]()

**© 2025 Sharon Aicler (saichler@gmail.com)**

**Layer 8 Data Mining Models (L8 Pollaris)** is an advanced **Polling/Parsing & Populating model service** designed for agnostic collection, parsing & populating abstract data mining models. It provides a highly flexible and extensible framework for defining device polling configurations and managing sophisticated collection operations across multiple protocols and device types in enterprise environments.

## Overview

L8 Pollaris is a core component of the Layer8 ecosystem, providing an intelligent centralized service for:

- **Advanced Device Polling Configuration**: Define sophisticated polling strategies for network devices and data sources
- **Target Management System**: Comprehensive target lifecycle management with state tracking, round-robin load balancing, and multi-collector support
- **Protocol Abstraction Layer**: Unified interface supporting multiple protocols (SNMP v2/v3, SSH, RESTCONF, NETCONF, gRPC, Kubernetes, GraphQL)
- **Flexible Data Mining Operations**: Support for diverse polling operations (GET, MAP, TABLE) with intelligent data extraction
- **Dynamic Group Management**: Organize devices and polling configurations into logical, hierarchical groups
- **High-Performance Distributed Caching**: Built-in distributed caching with improved logging and initialization
- **Enterprise-Grade Scalability**: Designed for large-scale enterprise environments with concurrent polling capabilities

## Architecture

The project is structured into several key components:

### Core Components

- **PollarisService**: Advanced service implementation handling comprehensive CRUD operations with enterprise-grade reliability
- **PollarisCenter**: Enhanced central management hub with improved logging, initialization tracking, and distributed caching capabilities
- **Protocol Buffer Definitions**: Type-safe message definitions for all entities with optimized serialization
- **Advanced Logging System**: Comprehensive logging infrastructure for debugging and monitoring service initialization
- **Distributed Cache Management**: High-performance caching layer with synchronization controls and initialization validation

### Target Management System

- **TargetService**: ORM-based service for persistent target management with PostgreSQL backend
- **TargetCallback**: Lifecycle hooks for target state changes with validation and collector routing
- **TargetLinks**: Abstraction layer for routing targets to appropriate collectors, parsers, caches, and persistence layers
- **InitTargets**: Automatic target initialization and recovery with round-robin distribution to collectors
- **StartStopTargets**: Bulk operations for starting/stopping targets by type with state management

### Protocol Buffer Types

#### Pollaris Configuration (`pollaris.proto`)
```protobuf
message L8Pollaris {
  string name = 1;
  string vendor = 2;
  string series = 3;
  string family = 4;
  string software = 5;
  string hardware = 6;
  string version = 7;
  repeated string groups = 8;
  map<string, L8Poll> polling = 9;
}

message L8Poll {
  string name = 1;
  string what = 2;
  L8C_Operation operation = 3;
  L8PProtocol protocol = 4;
  L8PCadencePlan cadence = 5;
  int64 timeout = 6;
  repeated L8PAttribute attributes = 7;
}
```

#### Target Configuration (`targets.proto`)
```protobuf
message L8PTarget {
  string target_id = 1;
  string links_id = 2;
  map<string, L8PHost> hosts = 3;
  L8PTargetState state = 4;
  L8PTargetType inventory_type = 5;
}

enum L8PTargetState {
  InvalidState = 0;
  Down = 1;
  Up = 2;
  Maintenance = 3;
  Offline = 4;
}

enum L8PTargetType {
  InvalidType = 0;
  Network_Device = 1;
  GPUS = 2;
  Hosts = 3;
  Virtual_Machine = 4;
  K8s_Cluster = 5;
  Storage = 6;
  Power = 7;
}
```

#### Collection Jobs (`jobs.proto`)
```protobuf
message CJob {
  string error = 1;
  bytes result = 2;
  int64 started = 3;
  int64 ended = 4;
  L8PCadencePlan cadence = 5;
  int64 timeout = 6;
  string target_id = 7;
  string host_id = 8;
  string pollaris_name = 9;
  string job_name = 10;
}
```

## Supported Protocols

L8 Pollaris supports multiple network management protocols:

- **SNMP v2/v3**: Traditional SNMP polling for network devices
- **SSH**: Command-line interface polling for device management
- **RESTCONF**: REST-based configuration protocol for modern infrastructure
- **NETCONF**: Network Configuration Protocol for structured device configuration
- **gRPC**: Modern RPC framework for high-performance communication
- **Kubernetes (Kubectl)**: Container orchestration platform polling
- **GraphQL**: Query language for flexible API interactions

## Supported Operations

- **L8C_Get**: Simple get operations for single values
- **L8C_Map**: Key-value mapping operations for structured data
- **L8C_Table**: Tabular data collection for bulk retrieval

## Supported Target Types

L8 Pollaris can manage various infrastructure target types:

- **Network Devices**: Routers, switches, firewalls, and other network equipment
- **GPUS**: GPU computing resources
- **Hosts**: Physical and virtual server hosts
- **Virtual Machines**: VM instances across hypervisors
- **Kubernetes Clusters**: K8s cluster monitoring and management
- **Storage**: Storage arrays and systems
- **Power**: Power distribution and UPS systems

## Recent Updates & Improvements

### Latest Enhancements (December 2025)

The following major improvements have been implemented in the latest version:

- **Apache 2.0 Licensing**: Added comprehensive Apache License 2.0 headers to all source files
- **Target Management System**: Introduced full-featured target lifecycle management with:
  - PostgreSQL-backed persistent storage via ORM
  - Round-robin load balancing to collectors
  - Bulk start/stop operations by target type
  - Address validation to prevent duplicate entries
  - Automatic target recovery on service restart
- **Enhanced Target States**: Support for Up, Down, Maintenance, and Offline states
- **Multi-Collector Architecture**: Targets are distributed across collectors using round-robin for load balancing
- **TargetLinks Abstraction**: Pluggable routing layer for directing targets to collectors, parsers, caches, and persistence

### Previous Enhancements

- **Enhanced Initialization Logging**: Comprehensive logging during PollarisCenter initialization to track init elements and service startup
- **Improved Cache Management**: Implemented `NewDistributedCacheNoSync` for better performance and control over synchronization
- **Initialization Data Handling**: Added dedicated `addForInit()` method for proper handling of initialization elements without triggering distributed cache events
- **Service Reliability**: Enhanced error handling and validation throughout the initialization process

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
- **l8parser**: Data parsing and boot configuration
- **l8srlz**: Serialization framework
- **l8test**: Testing infrastructure
- **l8bus**: Message bus and overlay networking with health monitoring
- **l8orm**: Object-relational mapping for persistent storage
- **l8ql**: Query language interpreter for data retrieval
- **probler**: Infrastructure probing framework

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

© 2025 Sharon Aicler (saichler@gmail.com)

This project is licensed under the Apache License, Version 2.0. You may obtain a copy of the License at:

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

## Project Structure

```
l8pollaris/
├── LICENSE                 # Apache 2.0 license
├── README.md               # This file
├── web.html                # Project website
├── go/                     # Go source code
│   ├── go.mod              # Go module definition
│   ├── go.sum              # Go module checksums
│   ├── pollaris/           # Main service implementation
│   │   ├── PollarisService.go    # Service interface implementation
│   │   ├── PollarisCenter.go     # Central management logic
│   │   ├── PollarisUtils.go      # Utility functions
│   │   └── targets/              # Target management subsystem
│   │       ├── TargetService.go      # Target service with ORM
│   │       ├── TargetCallback.go     # Lifecycle callbacks
│   │       ├── TargetLinks.go        # Routing abstraction
│   │       ├── InitTargets.go        # Target initialization
│   │       └── StartStopTargets.go   # Bulk operations
│   ├── tests/              # Unit tests
│   │   ├── Pollaris_test.go          # Pollaris test suite
│   │   ├── TargetService_test.go     # Target service tests
│   │   └── TestInit.go               # Test initialization
│   └── types/l8tpollaris/  # Generated Protocol Buffer types
│       ├── pollaris.pb.go        # Pollaris message types
│       ├── targets.pb.go         # Target message types
│       └── jobs.pb.go            # Collection job types
└── proto/                  # Protocol Buffer definitions
    ├── pollaris.proto      # Core polling configuration
    ├── targets.proto       # Target and host types
    └── jobs.proto          # Collection job types
```

## Support

For questions, issues, or contributions:
- **GitHub Issues**: https://github.com/saichler/l8pollaris/issues
- **Email**: saichler@gmail.com
- **Author**: Sharon Aicler
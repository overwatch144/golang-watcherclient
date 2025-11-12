# Golang Watcher Client

Go client library for OpenStack Watcher API - Resource Optimization Service.

## Installation
```bash
go get github.com/yourusername/golang-watcherclient
```

## Quick Start
```go
import "github.com/yourusername/golang-watcherclient/watcherclient"

client, err := watcherclient.NewClient(watcherclient.ClientOptions{
    AuthURL:         "http://keystone:5000/v3",
    Username:        "admin",
    Password:        "secret",
    ProjectName:     "admin",
    ProjectDomainID: "default",
    UserDomainID:    "default",
})

// Create an audit
audit := &watcherclient.Audit{
    Name:       "My Audit",
    AuditType:  "ONESHOT",
    Goal:       "server_consolidation",
    AutoTrigger: true,
}

result, err := client.CreateAudit(audit)
```

## Features

- Full Watcher API v1 support
- Keystone v3 authentication
- CRUD operations for all resource types:
  - Audits
  - Audit Templates
  - Action Plans
  - Actions
  - Goals
  - Strategies
  - Data Model

## Documentation

For detailed documentation, see the [examples](./examples) directory.

## License

Apache License 2.0
# Digo: Dependency injection for Go

Digo is a lightweight and extensible Dependency Injection container for Go — featuring constructor-based resolution, interface mapping, lifecycle scopes, and recursive instantiation.

## Features

### Constructor Registers  
Register factory functions (constructors) that the container uses to instantiate types automatically. Supports various types such as:
- Strings, numbers, booleans  
- Arrays and slices  
- Nested structs and complex object graphs  

### Interface Registers  
Map interfaces to their concrete implementations.  
This allows the container to resolve abstractions without hardcoding dependencies, making your architecture more flexible and testable.

### Scopes (Singleton / Transient)  
Supports lifecycle management for registered types:
- Singleton: A single instance shared across the application.
- Transient: A new instance is created on every resolution.

Scoped resolution ensures optimal performance and memory control for large applications.

### Recursive Resolution  
Dependencies are resolved recursively, including any sub-dependencies required by constructors.  
Supports:
- Optional fields (when nil is acceptable)  
- Default values (when handled in constructors)  

## Installation

Install the package using `go get`:

```bash
go get github.com/ChidemJean/digo
```

## Example

```go

// Init container
container := digo.New()

// Register a concrete implementation
container.Register(NewEmailNotifier, container.Singleton)

// Register a service that depends on the interface
container.Register(NewAlertService, container.Transient)

// Map the interface to the implementation
container.RegisterInterface((*Notifier)(nil), &EmailNotifier{})

// Resolve and use
alert := container.Resolve[Notifier]()
alert.Send()
# CONFDI (config-diff)
Is a small tool that takes a default config file and an override file, 
which then produces a new config file with merged values.     

## Usage
```go
go run main.go -default ./config-example/default.toml -override ./config-example/override.toml -merged merged.toml
```
This will generate a new `.toml` file with merged values.

## Flags
```go
 -default string
        default config path
  -log-format string
        log format (default "standard")
  -merged string
        merged config path
  -override string
        override config path

```

## Purpose
The main purpose for `confdi` is merging configmap values in K8s. 

When dealing with a lot of services, that need to share the same config, 
it can be challenging to keep track of the default configuration and the custom config 
for each service.    

Confdi allows a user to have a centralized configmap with default configuration and a 
centralized configmap with custom configuration for each service. 
Confdi will merge those two configuration files, and present a merged config file to the service.   
Confdi is meant to be ran as an `initContainer` where default and override config maps will be mounted to it, 
along with the `emptyDir` shared with pods, where the merged config file will be placed to.
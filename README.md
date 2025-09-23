# cqlshutil | Scylla DB assignment  

cqlsh is a command-line shell used to interact with ScyllaDB through CQL (Cassandra Query Language)


## How to use it

The specific assignment was developed using go 1.25.1

To build the project use the following command 
```bash
git clone https://github.com/otheodorakopoulos-lgtm/cqlshutil.git
go build -o cqlshutil.exe ./cmd/main.go
```

You can list the versions of cqlsh that are available for download in ScyllaDB Cloud using the following command:

```bash
./cqlshutil.exe list 
```
You can also use the -lt and -gt flags to:
-lt: Include only the versions that are lesser than the specified version.
-gt: Include only the versions that are greater than the specified version.

```bash
./cqlshutil.exe list -lt 2025.0.0 -gt 2025.1.3
```
You can download a specific version using the download command:
```
./cqlshutil.exe download 2025.2.2 -o test.tar.gz
or
./cqlshutil.exe download 2025.2.2 > test.tar.gz
```
Currently only the 2025.1 and 2025.2 minor releases are supported.

## Testing

To run all unit tests you can use the following command:

```bash
go test -v ./...
```


# fdio - Flogo Dot IO
A command-line interface for the Flogo Dot IO website. This tool is designed to help create the `items.toml` from which the [showcase](https://tibcosoftware.github.io/flogo/showcases/) and the [flogo cli](https://github.com/TIBCOSoftware/flogo-cli) can get their search results.

## Installing
There are a few ways to install this project

### Get the sources
You can get the sources for this project by simply running
```bash
$ go get -u github.com/retgits/fdio/...
```

You can create a binary using the `install` command
```bash
$ go install ./...
```

### Build from source
To build the fdio command-line interface simply run `go build`. This does require your system to have `gcc` installed. To build for Windows, you'll need to have CGO enabled
```bash 
$ GOOS=windows CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" go build
```

_You might need additional packages if you're running this command on a Linux system (like `apt-get install gcc-mingw-w64-x86-64 mingw-w64-x86-64-dev`)_

## Usage
```
A command-line interface for the Flogo Dot IO website

Usage:
  fdio [command]

Available Commands:
  backup      Performs a backup (upload) or restore (download) of a database file from Amazon S3
  crawl       Crawls GitHub to find new activities and triggers
  export      Export the database to a toml file
  help        Help about any command
  import      Import a toml file into the database
  init        Initialize the database in a new location
  query       Run a query against the database
  stats       Get statistics from the database

Flags:
      --db string   The path to the database (required)
  -h, --help        help for fdio
      --version     version for fdio

Use "fdio [command] --help" for more information about a command.
```
### Backup
```
Performs a backup (upload) or restore (download) of a database file from Amazon S3

Usage:
  fdio backup [flags]

Flags:
      --bucket string     The name of the S3 bucket (required)
  -h, --help              help for backup
      --overwrite         When restoring a file, should any existing file be overwritten (defaults to false)
      --region string     The region of the S3 bucket (defaults to us-west-2) (default "us-west-2")
      --restore           Restore (download) the file from Amazon S3

Global Flags:
      --db string   The path to the database (required)
```

### Crawl
```
Crawls GitHub to find new activities and triggers

Usage:
  fdio crawl [flags]

Flags:
  -h, --help            help for crawl
      --timeout float   The number of hours between now and the last repo update
      --type string     The type to look for, either trigger or activity (required)

Global Flags:
      --db string   The path to the database (required)
```
_The crawl command will create a `.crawl` file which lists the last date/time this command started_

### Export
```
Export the database to a toml file

Usage:
  fdio export [flags]

Flags:
  -h, --help          help for export
      --overwrite     Overwrite file if it exists
      --toml string   The path to the TOML file (required)

Global Flags:
      --db string   The path to the database (required)
```

### Import
```
Import a toml file into the database

Usage:
  fdio import [flags]

Flags:
  -h, --help          help for import
      --toml string   The path to the TOML file (required)

Global Flags:
      --db string   The path to the database (required)
```

### Init
```
Initialize the database in a new location

Usage:
  fdio init [flags]

Flags:
      --create   Create a new database if file doesn't exist
  -h, --help     help for init

Global Flags:
      --db string   The path to the database (required)
```

### Query
```
Run a query against the database

Usage:
  fdio query [flags]

Flags:
  -h, --help           help for query
  -q, --query string   The database query you want to run

Global Flags:
      --db string   The path to the database (required)
```
_With this command you can run any arbitrary query against the database, so do this at your own risk_

### Stats
```
Get statistics from the database

Usage:
  fdio stats [flags]

Flags:
  -h, --help   help for stats

Global Flags:
      --db string   The path to the database (required)
```

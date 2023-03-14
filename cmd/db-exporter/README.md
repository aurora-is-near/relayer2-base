# DB export tool

`db-exporter` is a tool to export and import a chain's data from relayer's internal database.
Exported data includes blocks, transactions, and logs, importing inserts data in the original format to the database including all indices.

## Install

```bash
$ go install .
```

## Usage

```bash
$ db-exporter -h
Usage of db-exporter:
  -a, --archive string   directory location for the exported data
  -c, --chainid uint     export/import data for this chainID
      --db string        the badgerDB database's directory
  -e, --export           export data from the given database
      --height uint      start export at specific block height
  -h, --help             print help
  -i, --import           import data to the given database
      --inspect          inspect data in the given database
```

### Examples

Export data of chain `1313161554` from relayer database `/tmp/relayer/data`.

```bash
$ db-exporter --export --db /tmp/relayer/data --chainid 1313161554 --archive ~/export-data
```

Import data into `/tmp/other/data` from archive at `~/export-data` under the chain `1313161554`.
The target database `/tmp/other/data` does not have to exist beforehand.

```bash
$ db-exporter --import --db /tmp/other/data --chainid 1313161554 --archive ~/export-data
```

Inspect the data inside `/tmp/other/data`.

```bash
$ db-exporter --inspect --db /tmp/other/data
```

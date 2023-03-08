package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aurora-is-near/relayer2-base/db/badger/core"
	"github.com/aurora-is-near/relayer2-base/db/codec"
	"github.com/dgraph-io/badger/v3"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
)

func main() {
	var (
		doExport    bool
		doImport    bool
		doInpect    bool
		help        bool
		dbPath      string
		archivePath string
		startHeight uint64
		chainID     uint64
	)

	pflag.BoolVarP(&help, "help", "h", false, "print help")
	pflag.BoolVarP(&doExport, "export", "e", false, "export data from the given database")
	pflag.BoolVarP(&doImport, "import", "i", false, "import data to the given database")
	pflag.BoolVar(&doInpect, "inspect", false, "inspect data in the given database")
	pflag.StringVar(&dbPath, "db", "", "the badgerDB database's directory")
	pflag.Uint64Var(&startHeight, "height", 0, "start export and specific block height")
	pflag.Uint64VarP(&chainID, "chainid", "c", 0, "export/import data for this chainID")
	pflag.StringVarP(&archivePath, "archive", "a", "", "directory location for the exported data")
	pflag.Parse()

	if help {
		pflag.Usage()
		os.Exit(0)
	}

	if info, err := os.Stat(dbPath); err != nil {
		panic(err)
	} else if !info.IsDir() {
		fmt.Println("path for DB", dbPath, "is not a directory")
		os.Exit(1)
	}

	dbConf := core.Config{
		MaxScanIterators:   1,
		ScanRangeThreshold: 1,
		FilterTtlMinutes:   1,
		GcIntervalSeconds: 1e9, // don't do GC
		BadgerConfig: badger.DefaultOptions(dbPath).
			WithDetectConflicts(false).
			WithNumCompactors(2).
			WithMemTableSize(64 << 15).
			WithNumMemtables(3).
			WithValueThreshold(3e5).
			WithValueLogFileSize(1<<28 - 1).
			WithLoggingLevel(badger.INFO),
	}

	fs := afero.NewBasePathFs(afero.NewOsFs(), archivePath)

	codec := codec.NewTinypackCodec()

	db, err := core.NewDB(dbConf, codec)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if doInpect {
		log.Printf(`inspecting relayer DB at %q`, dbPath)
		err := PrintDBInfo(db, os.Stdout)
		if err != nil {
			panic(err)
		}

	} else if doExport {
		log.Printf(`exporting relayer DB at %q to directory %q`, dbPath, archivePath)
		if info, err := os.Stat(archivePath); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(archivePath, 0)
			if err != nil {
				panic(err)
			}
		} else if !info.IsDir() {
			fmt.Println("archive path", archivePath, "is not a directory")
			os.Exit(1)
		}

		a, err := NewArchiver(fs, codec)
		if err != nil {
			panic(err)
		}
		e := Exporter{
			DB:       db.BadgerDB(),
			Archiver: a,
			ChainID:  chainID,
			Decoder:  codec,
			Height:   startHeight,
		}
		if err := e.Export(); err != nil {
			panic(err)
		}

	} else if doImport {
		log.Printf(`importing to relayer DB at %q from directory %q`, dbPath, archivePath)
		u, err := NewUnarchiver(fs, codec)
		if err != nil {
			panic(err)
		}
		i := Importer{
			DB:           db,
			Unarchiver:   u,
			ChainID:      chainID,
			PendingLimit: 1e5,
		}
		if err := i.Import(); err != nil {
			panic(err)
		}
	}
}

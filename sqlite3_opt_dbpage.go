// Copyright (C) 2019 Yasuhiro Matsumoto <mattn.jp@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

//go:build sqlite_dbpage
// +build sqlite_dbpage

package sqlite3

/*
#cgo CFLAGS: -DSQLITE_ENABLE_DBPAGE_VTAB=1
#cgo LDFLAGS: -lm

#define SQLITE_RSYNC_NO_MAIN
#define SQLITE_RSYNC_USE_H

#include "sqlite_rsync.h"

#include <stdlib.h>
#include <string.h>
#include <stdio.h>
*/
import "C"

import (
    "os"
    "errors"
	"unsafe"
)

func RsyncOrigin(dbPath string, isUri bool, pIn, pOut *os.File) error {
    cOriginPath := C.CString(dbPath)
	defer C.free(unsafe.Pointer(cOriginPath))

    ctx := C.SQLiteRsync{
		zOrigin:    cOriginPath,
		iProtocol:  C.PROTOCOL_VERSION,
		isRemote:   1,
        pIn:        C.fdopen((C.int)(pIn.Fd()), C.CString("rb")),
        pOut:       C.fdopen((C.int)(pOut.Fd()), C.CString("wb")),
	}

    flags := (C.int)(0)
    if isUri {
        flags = C.SQLITE_OPEN_URI
    }
    // start replicating
    C.originSide(&ctx, flags)

    if ctx.nErr > 0 {
        return errors.New("RSyncOrigin faild")
    }

    return nil
}

func RsyncReplica(replicaPath string, pIn, pOut *os.File) error {
    cReplicaPath := C.CString(replicaPath)
	defer C.free(unsafe.Pointer(cReplicaPath))

    ctx := C.SQLiteRsync{
		zReplica:   cReplicaPath,
		iProtocol:  C.PROTOCOL_VERSION,
		isRemote:   1,
        pIn:        C.fdopen((C.int)(pIn.Fd()), C.CString("rb")),
        pOut:       C.fdopen((C.int)(pOut.Fd()), C.CString("wb")),
	}

    // do the replica
    C.replicaSide(&ctx)
    
    if ctx.nErr > 0 {
        return errors.New("RSyncReplica faild")
    }

    return nil
}

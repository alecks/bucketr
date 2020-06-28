// Copyright 2020 The Bucketr Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func respond(w http.ResponseWriter, stat int, resp interface{}) {
	w.WriteHeader(stat)
	marshalled, _ := json.Marshal(resp)
	_, err := w.Write(marshalled)
	chk("[warn] http write", err)
}

func chk(scope string, err error) {
	if err != nil {
		fmt.Println(scope+":", err)
		// Allow the log to be a soft warning.
		if scope[:len(warnTag)] != warnTag {
			os.Exit(1)
		}
	}
}

//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"

	"github.com/edgexfoundry/edgex-go/internal/core/keeper"

	"github.com/gorilla/mux"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	keeper.Main(ctx, cancel, mux.NewRouter())
}

//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"context"
	"sync"

	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"

	"github.com/edgexfoundry/edgex-go/internal/core/keeper/container"
)

func BootstrapHandler(ctx context.Context, wg *sync.WaitGroup, _ startup.Timer, dic *di.Container) bool {
	dbClient := container.DBClientFrom(dic.Get)
	lc := bootstrapContainer.LoggingClientFrom(dic.Get)

	existedRegistrations, err := dbClient.Registrations()
	if err != nil {
		lc.Errorf("Failed to get registrations from database: %s", err.Error())
		return false
	}

	c := NewRegistry(ctx, wg, dic)
	for _, r := range existedRegistrations {
		c.Register(r)
	}

	dic.Update(di.ServiceConstructorMap{
		container.RegistryInterfaceName: func(get di.Get) interface{} {
			return c
		},
	})

	return true
}

//
// Copyright (C) 2024 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package interfaces

import "github.com/edgexfoundry/go-mod-core-contracts/v3/models"

// Registry defines the functionalities of a registry service
type Registry interface {
	// Register registers a service with the registration information,
	// and health check its status periodically
	Register(r models.Registration)
	// DeregisterByServiceId deregisters a service by its id and stop
	// health checking its status
	DeregisterByServiceId(id string)
}

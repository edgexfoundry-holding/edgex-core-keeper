//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package models

import "github.com/edgexfoundry/go-mod-core-contracts/v2/models"

type Registration struct {
	models.DBTimestamp
	ServiceId     string
	Status        string
	Host          string
	Port          int
	HealthCheck   HealthCheck
	LastConnected int64
}

type HealthCheck struct {
	Interval string
	Path     string
	Type     string
}

//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"

	"github.com/edgexfoundry/edgex-go/internal/core/keeper/container"
	"github.com/edgexfoundry/edgex-go/internal/core/keeper/dtos"
	"github.com/edgexfoundry/edgex-go/internal/core/keeper/models"
)

func AddRegistration(r models.Registration, dic *di.Container) errors.EdgeX {
	dbClient := container.DBClientFrom(dic.Get)
	r, err := dbClient.AddRegistration(r)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}

	registry := container.RegistryFrom(dic.Get)
	registry.Register(r)

	return nil
}

func UpdateRegistration(r models.Registration, dic *di.Container) errors.EdgeX {
	dbClient := container.DBClientFrom(dic.Get)
	err := dbClient.UpdateRegistration(r)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}

	registry := container.RegistryFrom(dic.Get)
	// remove the old service health check runner first, and then create a new one based on the updated registry
	registry.DeregisterByServiceId(r.ServiceId)
	registry.Register(r)

	return nil
}

func DeleteRegistration(id string, dic *di.Container) errors.EdgeX {
	if id == "" {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "serviceId is empty", nil)
	}

	dbClient := container.DBClientFrom(dic.Get)
	err := dbClient.DeleteRegistrationByServiceId(id)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}

	registry := container.RegistryFrom(dic.Get)
	registry.DeregisterByServiceId(id)

	return nil
}

func Registrations(dic *di.Container) ([]dtos.Registration, errors.EdgeX) {
	dbClient := container.DBClientFrom(dic.Get)
	registrations, err := dbClient.Registrations()
	if err != nil {
		return nil, errors.NewCommonEdgeXWrapper(err)
	}

	res := make([]dtos.Registration, len(registrations))
	for idx, r := range registrations {
		dto := dtos.FromRegistrationModelToDTO(r)
		res[idx] = dto
	}

	return res, nil
}

func RegistrationByServiceId(id string, dic *di.Container) (dtos.Registration, errors.EdgeX) {
	if id == "" {
		return dtos.Registration{}, errors.NewCommonEdgeX(errors.KindContractInvalid, "serviceId is empty", nil)
	}

	dbClient := container.DBClientFrom(dic.Get)
	r, err := dbClient.RegistrationByServiceId(id)
	if err != nil {
		return dtos.Registration{}, errors.NewCommonEdgeX(errors.KindEntityDoesNotExist, "registration not found by serviceId", err)
	}

	return dtos.FromRegistrationModelToDTO(r), nil
}

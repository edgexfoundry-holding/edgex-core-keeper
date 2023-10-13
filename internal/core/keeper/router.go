//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package keeper

import (
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/core/keeper/constants"
	controller "github.com/edgexfoundry/edgex-go/internal/core/keeper/controller/http"
	commonController "github.com/edgexfoundry/edgex-go/internal/pkg/controller/http"
	"github.com/edgexfoundry/edgex-go/internal/pkg/correlation"

	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"

	"github.com/gorilla/mux"
)

func LoadRestRoutes(r *mux.Router, dic *di.Container, serviceName string) {
	// Common
	cc := commonController.NewCommonController(dic, serviceName)
	r.HandleFunc(common.ApiPingRoute, cc.Ping).Methods(http.MethodGet)
	r.HandleFunc(common.ApiVersionRoute, cc.Version).Methods(http.MethodGet)
	r.HandleFunc(common.ApiConfigRoute, cc.Config).Methods(http.MethodGet)
	r.HandleFunc(common.ApiMetricsRoute, cc.Metrics).Methods(http.MethodGet)

	kv := controller.NewKVController(dic)
	r.HandleFunc(constants.ApiKVRoute, kv.Keys).Methods(http.MethodGet)
	r.HandleFunc(constants.ApiKVRoute, kv.AddKeys).Methods(http.MethodPut)
	r.HandleFunc(constants.ApiKVRoute, kv.DeleteKeys).Methods(http.MethodDelete)

	rc := controller.NewRegistryController(dic)
	r.HandleFunc(constants.ApiRegisterRoute, rc.Register).Methods(http.MethodPost)
	r.HandleFunc(constants.ApiRegisterRoute, rc.UpdateRegister).Methods(http.MethodPut)
	r.HandleFunc(constants.ApiAllRegistrationsRoute, rc.Registrations).Methods(http.MethodGet)
	r.HandleFunc(constants.ApiRegistrationByServiceIdRoute, rc.RegistrationByServiceId).Methods(http.MethodGet)
	r.HandleFunc(constants.ApiRegistrationByServiceIdRoute, rc.Deregister).Methods(http.MethodDelete)
	r.Use(correlation.ManageHeader)
}

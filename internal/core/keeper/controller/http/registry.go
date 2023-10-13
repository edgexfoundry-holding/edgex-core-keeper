//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"net/http"

	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
	"github.com/gorilla/mux"

	"github.com/edgexfoundry/edgex-go/internal/core/keeper/application"
	"github.com/edgexfoundry/edgex-go/internal/core/keeper/constants"
	"github.com/edgexfoundry/edgex-go/internal/core/keeper/dtos"
	requestDTO "github.com/edgexfoundry/edgex-go/internal/core/keeper/dtos/requests"
	"github.com/edgexfoundry/edgex-go/internal/core/keeper/dtos/responses"
	"github.com/edgexfoundry/edgex-go/internal/io"
	"github.com/edgexfoundry/edgex-go/internal/pkg"
	"github.com/edgexfoundry/edgex-go/internal/pkg/utils"
)

type RegistryController struct {
	reader io.DtoReader
	dic    *di.Container
}

func NewRegistryController(dic *di.Container) *RegistryController {
	return &RegistryController{
		reader: io.NewJsonDtoReader(),
		dic:    dic,
	}
}

func (rc *RegistryController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	lc := bootstrapContainer.LoggingClientFrom(rc.dic.Get)
	ctx := r.Context()

	var reqDTO requestDTO.AddRegistrationRequest
	edgexErr := rc.reader.Read(r.Body, &reqDTO)
	if edgexErr != nil {
		utils.WriteErrorResponse(w, ctx, lc, edgexErr, "")
		return
	}

	err := reqDTO.Validate()
	if err != nil {
		edgexErr = errors.NewCommonEdgeX(errors.KindContractInvalid, "bad AddRegistrationRequest type", err)
		utils.WriteErrorResponse(w, ctx, lc, edgexErr, "")
		return
	}

	registry := dtos.ToRegistrationModel(reqDTO.Registration)
	edgexErr = application.AddRegistration(registry, rc.dic)
	if edgexErr != nil {
		utils.WriteErrorResponse(w, ctx, lc, edgexErr, "")
		return
	}

	utils.WriteHttpHeader(w, ctx, http.StatusCreated)
}

func (rc *RegistryController) UpdateRegister(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	lc := bootstrapContainer.LoggingClientFrom(rc.dic.Get)
	ctx := r.Context()

	var reqDTO requestDTO.AddRegistrationRequest
	edgexErr := rc.reader.Read(r.Body, &reqDTO)
	if edgexErr != nil {
		utils.WriteErrorResponse(w, ctx, lc, edgexErr, "")
		return
	}

	err := reqDTO.Validate()
	if err != nil {
		edgexErr = errors.NewCommonEdgeX(errors.KindContractInvalid, "invalid registration request", err)
		utils.WriteErrorResponse(w, ctx, lc, edgexErr, "")
		return
	}

	registry := dtos.ToRegistrationModel(reqDTO.Registration)
	edgexErr = application.UpdateRegistration(registry, rc.dic)
	if edgexErr != nil {
		utils.WriteErrorResponse(w, ctx, lc, edgexErr, "")
		return
	}

	utils.WriteHttpHeader(w, ctx, http.StatusNoContent)
}

func (rc *RegistryController) Deregister(w http.ResponseWriter, r *http.Request) {
	lc := bootstrapContainer.LoggingClientFrom(rc.dic.Get)
	ctx := r.Context()

	// URL parameters
	vars := mux.Vars(r)
	id := vars[constants.ServiceId]

	err := application.DeleteRegistration(id, rc.dic)
	if err != nil {
		utils.WriteErrorResponse(w, ctx, lc, err, "")
		return
	}

	utils.WriteHttpHeader(w, ctx, http.StatusNoContent)
}

func (rc *RegistryController) Registrations(w http.ResponseWriter, r *http.Request) {
	lc := bootstrapContainer.LoggingClientFrom(rc.dic.Get)
	ctx := r.Context()

	dtos, err := application.Registrations(rc.dic)
	if err != nil {
		utils.WriteErrorResponse(w, ctx, lc, err, "")
		return
	}

	response := responses.NewMultiRegistrationsResponse("", "", http.StatusOK, uint32(len(dtos)), dtos)
	utils.WriteHttpHeader(w, ctx, http.StatusOK)
	pkg.EncodeAndWriteResponse(response, w, lc)
}

func (rc *RegistryController) RegistrationByServiceId(w http.ResponseWriter, r *http.Request) {
	lc := bootstrapContainer.LoggingClientFrom(rc.dic.Get)
	ctx := r.Context()

	// URL parameters
	vars := mux.Vars(r)
	id := vars[constants.ServiceId]

	dto, err := application.RegistrationByServiceId(id, rc.dic)
	if err != nil {
		utils.WriteErrorResponse(w, ctx, lc, err, "")
		return
	}

	response := responses.NewRegistrationResponse("", "", http.StatusOK, dto)
	utils.WriteHttpHeader(w, ctx, http.StatusOK)
	pkg.EncodeAndWriteResponse(response, w, lc)
}

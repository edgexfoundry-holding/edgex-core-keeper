//
// Copyright (C) 2022 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"encoding/json"
	"net/http"

	"github.com/edgexfoundry/edgex-go/internal/core/keeper/application"
	"github.com/edgexfoundry/edgex-go/internal/core/keeper/constants"
	requestDTO "github.com/edgexfoundry/edgex-go/internal/core/keeper/dtos/requests"
	responseDTO "github.com/edgexfoundry/edgex-go/internal/core/keeper/dtos/responses"
	kpContrUtils "github.com/edgexfoundry/edgex-go/internal/core/keeper/utils"
	edgexIO "github.com/edgexfoundry/edgex-go/internal/io"
	"github.com/edgexfoundry/edgex-go/internal/pkg"
	"github.com/edgexfoundry/edgex-go/internal/pkg/utils"

	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"

	"github.com/gorilla/mux"
)

type KVController struct {
	reader edgexIO.DtoReader
	dic    *di.Container
}

// NewKVController creates and initializes a KVController
func NewKVController(dic *di.Container) *KVController {
	return &KVController{
		reader: edgexIO.NewJsonDtoReader(),
		dic:    dic,
	}
}

func (rc *KVController) Keys(w http.ResponseWriter, r *http.Request) {
	lc := bootstrapContainer.LoggingClientFrom(rc.dic.Get)
	ctx := r.Context()

	// URL parameters
	vars := mux.Vars(r)
	key := vars[constants.Key]

	// parse URL query string for keyOnly and plaintext
	keysOnly, isRaw, err := kpContrUtils.ParseGetKeyRequestQueryString(r)
	if err != nil {
		utils.WriteErrorResponse(w, ctx, lc, err, "")
		return
	}

	resp, err := application.Keys(key, keysOnly, isRaw, rc.dic)
	if err != nil {
		utils.WriteErrorResponse(w, ctx, lc, err, "")
		return
	}

	response := responseDTO.NewMultiKVResponse("", "", http.StatusOK, resp)
	utils.WriteHttpHeader(w, ctx, http.StatusOK)
	pkg.EncodeAndWriteResponse(response, w, lc)
}

func (rc *KVController) AddKeys(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	lc := bootstrapContainer.LoggingClientFrom(rc.dic.Get)
	ctx := r.Context()

	// URL parameters
	vars := mux.Vars(r)
	key := vars[constants.Key]

	// parse URL query string for flatten
	isFlatten, err := kpContrUtils.ParseAddKeyRequestQueryString(r)
	if err != nil {
		utils.WriteErrorResponse(w, ctx, lc, err, "")
		return
	}

	var reqDTO requestDTO.AddKeysRequest
	err = rc.reader.Read(r.Body, &reqDTO)
	if err != nil {
		utils.WriteErrorResponse(w, ctx, lc, err, "")
		return
	}

	err = reqDTO.Validate()
	if err != nil {
		utils.WriteErrorResponse(w, ctx, lc, err, "")
		return
	}

	kvModel := requestDTO.AddKeysReqToKVModels(reqDTO, key)
	keys, err := application.AddKeys(kvModel, isFlatten, rc.dic)
	if err != nil {
		utils.WriteErrorResponse(w, ctx, lc, err, "")
		return
	}
	KVModelBytes, encodeErr := json.Marshal(kvModel)
	if encodeErr != nil {
		utils.WriteErrorResponse(w, ctx, lc, errors.NewCommonEdgeX(errors.KindServerError, "KV struct encoding failed", encodeErr), "")
		return
	}

	// publish the key change event
	go application.PublishKeyChange(KVModelBytes, key, ctx, rc.dic)

	response := responseDTO.NewKeysResponse("", "", http.StatusOK, keys)
	utils.WriteHttpHeader(w, ctx, http.StatusOK)
	pkg.EncodeAndWriteResponse(response, w, lc)
}

func (rc *KVController) DeleteKeys(w http.ResponseWriter, r *http.Request) {
	lc := bootstrapContainer.LoggingClientFrom(rc.dic.Get)
	ctx := r.Context()

	// URL parameters
	vars := mux.Vars(r)
	key := vars[constants.Key]

	// parse URL query string for prefixMatch
	prefixMatch, err := kpContrUtils.ParseDeleteKeyRequestQueryString(r)
	if err != nil {
		utils.WriteErrorResponse(w, ctx, lc, err, "")
		return
	}

	resp, err := application.DeleteKeys(key, prefixMatch, rc.dic)
	if err != nil {
		utils.WriteErrorResponse(w, ctx, lc, err, "")
		return
	}

	response := responseDTO.NewKeysResponse("", "", http.StatusOK, resp)
	utils.WriteHttpHeader(w, ctx, http.StatusOK)
	pkg.EncodeAndWriteResponse(response, w, lc)
}

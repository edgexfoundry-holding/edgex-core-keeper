//
// Copyright (C) 2022-2023 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package requests

import (
	"github.com/edgexfoundry/edgex-go/internal/core/keeper/dtos"

	commonDTO "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
)

type AddRegistrationRequest struct {
	commonDTO.BaseRequest `json:",inline"`
	Registration          dtos.Registration `json:"registration"`
}

func (a *AddRegistrationRequest) Validate() error {
	err := a.Registration.Validate()
	if err != nil {
		return err
	}
	return nil
}

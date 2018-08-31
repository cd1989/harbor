// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"fmt"

	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/common/utils/log"
	"github.com/goharbor/harbor/src/ui/utils"
)

// RetagAPI retag an image
type RetagAPI struct {
	BaseController
}

func (r *RetagAPI) Retag() {
	if !r.SecurityCtx.IsAuthenticated() {
		r.HandleUnauthorized()
		return
	}

	request := models.RetagRequest{}
	r.DecodeJSONReq(&request)

	srcImage, err := models.ParseImage(request.SrcImage)
	if err != nil {
		r.HandleBadRequest(fmt.Sprintf("invalid src image string '%s', should in format '<project>/<repo>:<tag>'", request.SrcImage))
		return
	}
	destImage, err := models.ParseImage(request.DestImage)
	if err != nil {
		r.HandleBadRequest(fmt.Sprintf("invalid dest image string '%s', should in format '<project>/<repo>:<tag>'", request.DestImage))
		return
	}

	if !r.SecurityCtx.HasReadPerm(srcImage.Project) {
		log.Errorf("user has no read permission to project '%s'", srcImage.Project)
		r.HandleUnauthorized()
	}

	if !r.SecurityCtx.HasWritePerm(destImage.Project) {
		log.Errorf("user has no write permission to project '%s'", destImage.Project)
		r.HandleUnauthorized()
	}

	if err = utils.Retag(srcImage, destImage); err != nil {
		r.HandleInternalServerError(fmt.Sprintf("%v", err))
	}
}



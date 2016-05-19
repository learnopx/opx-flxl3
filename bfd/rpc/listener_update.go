//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __  
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  | 
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  | 
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   | 
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  | 
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__| 
//                                                                                                           

package rpc

import (
	"bfdd"
	"errors"
	"fmt"
	"l3/bfd/server"
)

func (h *BFDHandler) UpdateBfdGlobal(origConf *bfdd.BfdGlobal, newConf *bfdd.BfdGlobal, attrset []bool, op string) (bool, error) {
	h.logger.Info(fmt.Sprintln("Original global config attrs:", origConf))
	h.logger.Info(fmt.Sprintln("New global config attrs:", newConf))
	gConf := server.GlobalConfig{
		Enable: newConf.Enable,
	}
	h.server.GlobalConfigCh <- gConf
	return true, nil
}

func (h *BFDHandler) UpdateBfdSession(origConf *bfdd.BfdSession, newConf *bfdd.BfdSession, attrset []bool, op string) (bool, error) {
	if newConf == nil {
		err := errors.New("Invalid Session Configuration")
		return false, err
	}
	h.logger.Info(fmt.Sprintln("Update session config attrs:", newConf))
	return h.SendBfdSessionConfig(newConf), nil
}

func (h *BFDHandler) UpdateBfdSessionParam(origConf *bfdd.BfdSessionParam, newConf *bfdd.BfdSessionParam, attrset []bool, op string) (bool, error) {
	if newConf == nil {
		err := errors.New("Invalid Session Param Configuration")
		return false, err
	}
	h.logger.Info(fmt.Sprintln("Update session Param config attrs:", newConf))
	return h.SendBfdSessionParamConfig(newConf), nil
}

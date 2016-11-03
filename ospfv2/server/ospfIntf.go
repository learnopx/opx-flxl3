//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package server

import (
	"l3/ospfv2/objects"
)

func (server *OSPFV2Server) updateIntf(newCfg, oldCfg *objects.Ospfv2Intf, attrset []bool) (bool, error) {
	server.logger.Info("Intf configuration update")
	return true, nil
}

func (server *OSPFV2Server) createIntf(cfg *objects.Ospfv2Intf) (bool, error) {
	server.logger.Info("Intf configuration create")
	return true, nil
}

func (server *OSPFV2Server) deleteIntf(cfg *objects.Ospfv2Intf) (bool, error) {
	server.logger.Info("Intf configuration delete")
	return true, nil
}

func (server *OSPFV2Server) getIntfState(ipAddr, addressLessIfIdx uint32) (*objects.Ospfv2IntfState, error) {
	var retObj objects.Ospfv2IntfState
	return &retObj, nil
}

func (server *OSPFV2Server) getBulkIntfState(fromIdx, cnt int) (*objects.Ospfv2IntfStateGetInfo, error) {
	var retObj objects.Ospfv2IntfStateGetInfo
	return &retObj, nil
}

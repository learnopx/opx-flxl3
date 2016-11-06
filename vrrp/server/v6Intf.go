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
package server

import (
	"l3/vrrp/config"
	"l3/vrrp/debug"
	"net"
)

type V6Intf struct {
	Cfg     config.Ipv6Info // this is ipv6 interface created on the system config
	Vrrpkey *KeyInfo
}

func (intf *V6Intf) Init(obj *config.BaseIpInfo) {
	intf.Cfg.Info.IntfRef = obj.IntfRef
	intf.Cfg.Info.IfIndex = obj.IfIndex
	intf.Cfg.Info.OperState = obj.OperState
	ip, _, _ := net.ParseCIDR(obj.IpAddr)
	if ip.IsLinkLocalUnicast() {
		intf.Cfg.LinkScopeAddr = ip.String()
	} else {
		intf.Cfg.Info.IpAddr = ip.String()
	}
	intf.Vrrpkey = nil
	debug.Logger.Debug("v6 ip interface initialized:", intf.Cfg)
}

func (intf *V6Intf) Update(obj *config.BaseIpInfo) {
	// most likely update of OperState only
	intf.Cfg.Info.OperState = obj.OperState
}

func (intf *V6Intf) DeInit(obj *config.BaseIpInfo) {
}

func (intf *V6Intf) GetObjFromDb(l3Info *config.BaseIpInfo) {
	l3Info.IpAddr = intf.Cfg.Info.IpAddr
	l3Info.OperState = intf.Cfg.Info.OperState
}

func (intf *V6Intf) SetVrrpIntfKey(key KeyInfo) {
	intf.Vrrpkey = &key
}

func (intf *V6Intf) GetVrrpIntfKey() *KeyInfo {
	return intf.Vrrpkey
}

func (intf *V6Intf) GetIntfRef() string {
	return intf.Cfg.Info.IntfRef
}

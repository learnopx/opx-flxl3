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

package flexswitch

import (
	"l3/ndp/lib"
	"l3/ndp/lib/flexswitch"
	"sync"
	"utils/commonDefs"
)

// this is for notification messages
var ndpClientInst *flexswitch.NdpdClientStruct = nil
var ndpOnce sync.Once

func GetNdpInst() *flexswitch.NdpdClientStruct {
	ndpOnce.Do(func() {
		ndpClientInst = &flexswitch.NdpdClientStruct{}
	})
	return ndpClientInst
}

// this is initializing the NBMgr for ndp client
func InitNdpInst(plugin, paramsDir string, clntList []commonDefs.ClientJson, ndpHdl flexswitch.NdpdClientStruct) lib.NdpdClientIntf {
	return lib.NewNdpClient(plugin, paramsDir, clntList, ndpHdl)
}

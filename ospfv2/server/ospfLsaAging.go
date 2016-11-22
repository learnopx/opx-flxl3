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

import ()

func (server *OSPFV2Server) processLsdbAgeSelfOrigRouterLsa(lsdbKey LsdbKey, lsaKey LsaKey, lsa *RouterLsa) bool {
	//Increment LSA age
	if lsa.LsaMd.LSAge < MAX_AGE {
		lsa.LsaMd.LSAge++
	}
	//If Age = multiples of CheckAge compute checksum and verify if error raise an alarm
	if (lsa.LsaMd.LSAge % CHECK_AGE) == 0 {
		lsaEnc := encodeRouterLsa(*lsa, lsaKey)
		checksumOffset := uint16(14)
		cSum := computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		if cSum != 0 {
			server.logger.Err("Some serious problem, may be memory corruption")
			return false
		}
	}
	//TODO: If Age=LSRefreshTime Regenerate
	if lsa.LsaMd.LSAge == LS_REFRESH_TIME {
		msg := GenerateRouterLSAMsg{
			AreaId: lsdbKey.AreaId,
		}
		server.reGenerateRouterLSA(msg)
		return true
	}
	return false
}

func (server *OSPFV2Server) processLsdbAgeSelfOrigNetworkLsa(lsdbKey LsdbKey, lsaKey LsaKey, lsa *NetworkLsa) bool {
	//Increment LSA age
	if lsa.LsaMd.LSAge < MAX_AGE {
		lsa.LsaMd.LSAge++
	}
	//If Age = multiples of CheckAge compute checksum and verify if error raise an alarm
	if (lsa.LsaMd.LSAge % CHECK_AGE) == 0 {
		lsaEnc := encodeNetworkLsa(*lsa, lsaKey)
		checksumOffset := uint16(14)
		cSum := computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		if cSum != 0 {
			server.logger.Err("Some serious problem, may be memory corruption")
			return false
		}
	}
	//TODO: If Age=LSRefreshTime Regenerate
	return false
}

func (server *OSPFV2Server) processLsdbAgeSelfOrigSummary3Lsa(lsdbKey LsdbKey, lsaKey LsaKey, lsa *SummaryLsa) bool {
	//Increment LSA age
	if lsa.LsaMd.LSAge < MAX_AGE {
		lsa.LsaMd.LSAge++
	}
	//If Age = multiples of CheckAge compute checksum and verify if error raise an alarm
	if (lsa.LsaMd.LSAge % CHECK_AGE) == 0 {
		lsaEnc := encodeSummaryLsa(*lsa, lsaKey)
		checksumOffset := uint16(14)
		cSum := computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		if cSum != 0 {
			server.logger.Err("Some serious problem, may be memory corruption")
			return false
		}
	}
	//TODO: If Age=LSRefreshTime Regenerate
	return false
}

func (server *OSPFV2Server) processLsdbAgeSelfOrigSummary4Lsa(lsdbKey LsdbKey, lsaKey LsaKey, lsa *SummaryLsa) bool {
	//Increment LSA age
	if lsa.LsaMd.LSAge < MAX_AGE {
		lsa.LsaMd.LSAge++
	}
	//If Age = multiples of CheckAge compute checksum and verify if error raise an alarm
	if (lsa.LsaMd.LSAge % CHECK_AGE) == 0 {
		lsaEnc := encodeSummaryLsa(*lsa, lsaKey)
		checksumOffset := uint16(14)
		cSum := computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		if cSum != 0 {
			server.logger.Err("Some serious problem, may be memory corruption")
			return false
		}
	}
	//TODO: If Age=LSRefreshTime Regenerate
	return false
}

func (server *OSPFV2Server) processLsdbAgeSelfOrigASExternalLsa(lsdbKey LsdbKey, lsaKey LsaKey, lsa *ASExternalLsa) bool {
	//Increment LSA age
	if lsa.LsaMd.LSAge < MAX_AGE {
		lsa.LsaMd.LSAge++
	}
	//If Age = multiples of CheckAge compute checksum and verify if error raise an alarm
	if (lsa.LsaMd.LSAge % CHECK_AGE) == 0 {
		lsaEnc := encodeASExternalLsa(*lsa, lsaKey)
		checksumOffset := uint16(14)
		cSum := computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		if cSum != 0 {
			server.logger.Err("Some serious problem, may be memory corruption")
			return false
		}
	}
	//TODO: If Age=LSRefreshTime Regenerate
	if lsa.LsaMd.LSAge == LS_REFRESH_TIME {
		routeInfo := RouteInfo{
			NwAddr:  lsaKey.LSId,
			Netmask: lsa.Netmask,
			Metric:  lsa.Metric,
		}
		_, exist := server.LsdbData.ExtRouteInfoMap[routeInfo]
		if exist {
			server.reGenerateASExternalLSAForGivenArea(routeInfo, lsdbKey.AreaId)
			return true
		}
	}
	return false
}

func (server *OSPFV2Server) processLsdbAgeSelfOrigLsa(lsdbKey LsdbKey, lsaKey LsaKey, lsaEnt interface{}) bool {
	switch lsaKey.LSType {
	case RouterLSA:
		lsa, ok := lsaEnt.(*RouterLsa)
		if !ok {
			server.logger.Err("Unable to assert lsa")
			return false
		}
		return server.processLsdbAgeSelfOrigRouterLsa(lsdbKey, lsaKey, lsa)
	case NetworkLSA:
		lsa, ok := lsaEnt.(*NetworkLsa)
		if !ok {
			server.logger.Err("Unable to assert lsa")
			return false
		}
		return server.processLsdbAgeSelfOrigNetworkLsa(lsdbKey, lsaKey, lsa)
	case Summary3LSA:
		lsa, ok := lsaEnt.(*SummaryLsa)
		if !ok {
			server.logger.Err("Unable to assert lsa")
			return false
		}
		return server.processLsdbAgeSelfOrigSummary3Lsa(lsdbKey, lsaKey, lsa)
	case Summary4LSA:
		lsa, ok := lsaEnt.(*SummaryLsa)
		if !ok {
			server.logger.Err("Unable to assert lsa")
			return false
		}
		return server.processLsdbAgeSelfOrigSummary4Lsa(lsdbKey, lsaKey, lsa)
	case ASExternalLSA:
		lsa, ok := lsaEnt.(*ASExternalLsa)
		if !ok {
			server.logger.Err("Unable to assert lsa")
			return false
		}
		return server.processLsdbAgeSelfOrigASExternalLsa(lsdbKey, lsaKey, lsa)
	}
	return false
}

func (server *OSPFV2Server) processLsdbAgeNonSelfRouterLsa(lsdbKey LsdbKey, lsaKey LsaKey, lsa *RouterLsa) (LsdbToFloodLSAMsg, bool) {
	var msg LsdbToFloodLSAMsg
	//Increment LSA age
	if lsa.LsaMd.LSAge < MAX_AGE {
		lsa.LsaMd.LSAge++
	}
	//If Age = multiples of CheckAge compute checksum and verify if error raise an alarm
	if (lsa.LsaMd.LSAge % CHECK_AGE) == 0 {
		lsaEnc := encodeRouterLsa(*lsa, lsaKey)
		checksumOffset := uint16(14)
		cSum := computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		if cSum != 0 {
			server.logger.Err("Some serious problem, may be memory corruption")
			return msg, false
		}
	}
	if lsa.LsaMd.LSAge == MAX_AGE {
		msg.AreaId = lsdbKey.AreaId
		msg.LsaKey = lsaKey
		msg.LsaData = *lsa
		return msg, true
	}
	return msg, false
}

func (server *OSPFV2Server) processLsdbAgeNonSelfNetworkLsa(lsdbKey LsdbKey, lsaKey LsaKey, lsa *NetworkLsa) (LsdbToFloodLSAMsg, bool) {
	var msg LsdbToFloodLSAMsg
	//Increment LSA age
	if lsa.LsaMd.LSAge < MAX_AGE {
		lsa.LsaMd.LSAge++
	}
	//If Age = multiples of CheckAge compute checksum and verify if error raise an alarm
	if (lsa.LsaMd.LSAge % CHECK_AGE) == 0 {
		lsaEnc := encodeNetworkLsa(*lsa, lsaKey)
		checksumOffset := uint16(14)
		cSum := computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		if cSum != 0 {
			server.logger.Err("Some serious problem, may be memory corruption")
			return msg, false
		}
	}
	if lsa.LsaMd.LSAge == MAX_AGE {
		msg.AreaId = lsdbKey.AreaId
		msg.LsaKey = lsaKey
		msg.LsaData = *lsa
		return msg, true
	}
	return msg, false
}

func (server *OSPFV2Server) processLsdbAgeNonSelfSummary3Lsa(lsdbKey LsdbKey, lsaKey LsaKey, lsa *SummaryLsa) (LsdbToFloodLSAMsg, bool) {
	var msg LsdbToFloodLSAMsg
	//Increment LSA age
	if lsa.LsaMd.LSAge < MAX_AGE {
		lsa.LsaMd.LSAge++
	}
	//If Age = multiples of CheckAge compute checksum and verify if error raise an alarm
	if (lsa.LsaMd.LSAge % CHECK_AGE) == 0 {
		lsaEnc := encodeSummaryLsa(*lsa, lsaKey)
		checksumOffset := uint16(14)
		cSum := computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		if cSum != 0 {
			server.logger.Err("Some serious problem, may be memory corruption")
			return msg, false
		}
	}
	if lsa.LsaMd.LSAge == MAX_AGE {
		msg.AreaId = lsdbKey.AreaId
		msg.LsaKey = lsaKey
		msg.LsaData = *lsa
		return msg, true
	}
	return msg, false
}

func (server *OSPFV2Server) processLsdbAgeNonSelfSummary4Lsa(lsdbKey LsdbKey, lsaKey LsaKey, lsa *SummaryLsa) (LsdbToFloodLSAMsg, bool) {
	var msg LsdbToFloodLSAMsg
	//Increment LSA age
	if lsa.LsaMd.LSAge < MAX_AGE {
		lsa.LsaMd.LSAge++
	}
	//If Age = multiples of CheckAge compute checksum and verify if error raise an alarm
	if (lsa.LsaMd.LSAge % CHECK_AGE) == 0 {
		lsaEnc := encodeSummaryLsa(*lsa, lsaKey)
		checksumOffset := uint16(14)
		cSum := computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		if cSum != 0 {
			server.logger.Err("Some serious problem, may be memory corruption")
			return msg, false
		}
	}
	if lsa.LsaMd.LSAge == MAX_AGE {
		msg.AreaId = lsdbKey.AreaId
		msg.LsaKey = lsaKey
		msg.LsaData = *lsa
		return msg, true
	}
	return msg, false
}

func (server *OSPFV2Server) processLsdbAgeNonSelfASExternalLsa(lsdbKey LsdbKey, lsaKey LsaKey, lsa *ASExternalLsa) (LsdbToFloodLSAMsg, bool) {
	var msg LsdbToFloodLSAMsg
	//Increment LSA age
	if lsa.LsaMd.LSAge < MAX_AGE {
		lsa.LsaMd.LSAge++
	}
	//If Age = multiples of CheckAge compute checksum and verify if error raise an alarm
	if (lsa.LsaMd.LSAge % CHECK_AGE) == 0 {
		lsaEnc := encodeASExternalLsa(*lsa, lsaKey)
		checksumOffset := uint16(14)
		cSum := computeFletcherChecksum(lsaEnc[2:], checksumOffset)
		if cSum != 0 {
			server.logger.Err("Some serious problem, may be memory corruption")
			return msg, false
		}
	}
	if lsa.LsaMd.LSAge == MAX_AGE {
		msg.AreaId = lsdbKey.AreaId
		msg.LsaKey = lsaKey
		msg.LsaData = *lsa
		return msg, true
	}
	return msg, false
}

func (server *OSPFV2Server) processLsdbAgeNonSelfLsa(lsdbKey LsdbKey, lsaKey LsaKey, lsaEnt interface{}) (LsdbToFloodLSAMsg, bool) {
	var msg LsdbToFloodLSAMsg
	switch lsaKey.LSType {
	case RouterLSA:
		lsa, ok := lsaEnt.(*RouterLsa)
		if !ok {
			server.logger.Err("Unable to assert lsa")
			return msg, false
		}
		return server.processLsdbAgeNonSelfRouterLsa(lsdbKey, lsaKey, lsa)
	case NetworkLSA:
		lsa, ok := lsaEnt.(*NetworkLsa)
		if !ok {
			server.logger.Err("Unable to assert lsa")
			return msg, false
		}
		return server.processLsdbAgeNonSelfNetworkLsa(lsdbKey, lsaKey, lsa)
	case Summary3LSA:
		lsa, ok := lsaEnt.(*SummaryLsa)
		if !ok {
			server.logger.Err("Unable to assert lsa")
			return msg, false
		}
		return server.processLsdbAgeNonSelfSummary3Lsa(lsdbKey, lsaKey, lsa)
	case Summary4LSA:
		lsa, ok := lsaEnt.(*SummaryLsa)
		if !ok {
			server.logger.Err("Unable to assert lsa")
			return msg, false
		}
		return server.processLsdbAgeNonSelfSummary4Lsa(lsdbKey, lsaKey, lsa)
	case ASExternalLSA:
		lsa, ok := lsaEnt.(*ASExternalLsa)
		if !ok {
			server.logger.Err("Unable to assert lsa")
			return msg, false
		}
		return server.processLsdbAgeNonSelfASExternalLsa(lsdbKey, lsaKey, lsa)
	}
	return msg, false
}

func (server *OSPFV2Server) processLsdbAgingTicker() {
	var lsdbToFloodLSAMsgList []LsdbToFloodLSAMsg
	var needSPFCalcRouter bool
	var needSPFCalcNetwork bool
	var needSPFCalcSummary3 bool
	var needSPFCalcSummary4 bool
	var needSPFCalcASExternal bool
	for lsdbKey, lsdbEnt := range server.LsdbData.AreaLsdb {
		for lsaKey, lsaEnt := range lsdbEnt.RouterLsaMap {
			selfOrigEnt, exist := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
			if exist {
				_, exist := selfOrigEnt[lsaKey]
				if exist {
					needSPFCalcRouter = server.processLsdbAgeSelfOrigLsa(lsdbKey, lsaKey, &lsaEnt)
				} else {
					lsdbToFloodLSAMsg, flag := server.processLsdbAgeNonSelfLsa(lsdbKey, lsaKey, &lsaEnt)
					if flag == true {
						lsdbToFloodLSAMsgList = append(lsdbToFloodLSAMsgList, lsdbToFloodLSAMsg)
					}
				}
				if lsaEnt.LsaMd.LSAge == MAX_AGE {
					delete(server.LsdbData.AreaLsdb[lsdbKey].RouterLsaMap, lsaKey)
				} else {
					server.LsdbData.AreaLsdb[lsdbKey].RouterLsaMap[lsaKey] = lsaEnt
				}
			} else {
				server.logger.Err("This should Not happen some serious problem")
			}
		}
		for lsaKey, lsaEnt := range lsdbEnt.NetworkLsaMap {
			selfOrigEnt, exist := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
			if exist {
				_, exist := selfOrigEnt[lsaKey]
				if exist {
					needSPFCalcNetwork = server.processLsdbAgeSelfOrigLsa(lsdbKey, lsaKey, &lsaEnt)
				} else {
					lsdbToFloodLSAMsg, flag := server.processLsdbAgeNonSelfLsa(lsdbKey, lsaKey, &lsaEnt)
					if flag == true {
						lsdbToFloodLSAMsgList = append(lsdbToFloodLSAMsgList, lsdbToFloodLSAMsg)
					}
				}
				if lsaEnt.LsaMd.LSAge == MAX_AGE {
					delete(server.LsdbData.AreaLsdb[lsdbKey].NetworkLsaMap, lsaKey)
				} else {
					server.LsdbData.AreaLsdb[lsdbKey].NetworkLsaMap[lsaKey] = lsaEnt
				}
			} else {
				server.logger.Err("This should Not happen some serious problem")
			}
		}
		for lsaKey, lsaEnt := range lsdbEnt.Summary3LsaMap {
			selfOrigEnt, exist := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
			if exist {
				_, exist := selfOrigEnt[lsaKey]
				if exist {
					needSPFCalcSummary3 = server.processLsdbAgeSelfOrigLsa(lsdbKey, lsaKey, &lsaEnt)
				} else {
					lsdbToFloodLSAMsg, flag := server.processLsdbAgeNonSelfLsa(lsdbKey, lsaKey, &lsaEnt)
					if flag == true {
						lsdbToFloodLSAMsgList = append(lsdbToFloodLSAMsgList, lsdbToFloodLSAMsg)
					}
				}
				if lsaEnt.LsaMd.LSAge == MAX_AGE {
					delete(server.LsdbData.AreaLsdb[lsdbKey].Summary3LsaMap, lsaKey)
				} else {
					server.LsdbData.AreaLsdb[lsdbKey].Summary3LsaMap[lsaKey] = lsaEnt
				}
			} else {
				server.logger.Err("This should Not happen some serious problem")
			}
		}
		for lsaKey, lsaEnt := range lsdbEnt.Summary4LsaMap {
			selfOrigEnt, exist := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
			if exist {
				_, exist := selfOrigEnt[lsaKey]
				if exist {
					needSPFCalcSummary4 = server.processLsdbAgeSelfOrigLsa(lsdbKey, lsaKey, &lsaEnt)
				} else {
					lsdbToFloodLSAMsg, flag := server.processLsdbAgeNonSelfLsa(lsdbKey, lsaKey, &lsaEnt)
					if flag == true {
						lsdbToFloodLSAMsgList = append(lsdbToFloodLSAMsgList, lsdbToFloodLSAMsg)
					}
				}
				if lsaEnt.LsaMd.LSAge == MAX_AGE {
					delete(server.LsdbData.AreaLsdb[lsdbKey].Summary4LsaMap, lsaKey)
				} else {
					server.LsdbData.AreaLsdb[lsdbKey].Summary4LsaMap[lsaKey] = lsaEnt
				}
			} else {
				server.logger.Err("This should Not happen some serious problem")
			}
		}
		for lsaKey, lsaEnt := range lsdbEnt.ASExternalLsaMap {
			selfOrigEnt, exist := server.LsdbData.AreaSelfOrigLsa[lsdbKey]
			if exist {
				_, exist := selfOrigEnt[lsaKey]
				if exist {
					needSPFCalcASExternal = server.processLsdbAgeSelfOrigLsa(lsdbKey, lsaKey, &lsaEnt)
				} else {
					lsdbToFloodLSAMsg, flag := server.processLsdbAgeNonSelfLsa(lsdbKey, lsaKey, &lsaEnt)
					if flag == true {
						lsdbToFloodLSAMsgList = append(lsdbToFloodLSAMsgList, lsdbToFloodLSAMsg)
					}
				}
				if lsaEnt.LsaMd.LSAge == MAX_AGE {
					delete(server.LsdbData.AreaLsdb[lsdbKey].ASExternalLsaMap, lsaKey)
				} else {
					server.LsdbData.AreaLsdb[lsdbKey].ASExternalLsaMap[lsaKey] = lsaEnt
				}
			} else {
				server.logger.Err("This should Not happen some serious problem")
			}
		}
	}
	server.SendMsgFromLsdbToFloodLsa(lsdbToFloodLSAMsgList)
	if needSPFCalcRouter == true ||
		needSPFCalcNetwork == true ||
		needSPFCalcSummary3 == true ||
		needSPFCalcSummary4 == true ||
		needSPFCalcASExternal == true {
		server.CalcSPFAndRoutingTbl()
	}
}

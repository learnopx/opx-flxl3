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
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func (server *OSPFV2Server) InitFlooding() {
	server.logger.Debug("Flood: Init done.")
	server.FloodData.FloodCtrlCh = make(chan bool)
	server.FloodData.FloodCtrlReplyCh = make(chan bool)
}

func (server *OSPFV2Server) DeinitFlooding() {
	server.logger.Debug("Flooding deinitialised.:")
}

func (server *OSPFV2Server) StartFlooding() {
	server.InitFlooding()
	go server.ProcessFlooding()
}
func (server *OSPFV2Server) ProcessFlooding() {
	for {
		select {
		case lsdbToFloodArray := <-server.MessagingChData.LsdbToFloodChData.LsdbToFloodLSACh:
			server.logger.Debug("Flood: received self originated lsa.", lsdbToFloodArray)
			for _, lsdbToFloodData := range lsdbToFloodArray {
				server.ProcessLsdbToFloodMsg(lsdbToFloodData)
			}

		case nbrToFloodData := <-server.MessagingChData.NbrFSMToFloodChData.LsaFloodCh:
			server.logger.Debug("Flood: Received flood message from nbr ", nbrToFloodData.NbrKey)
			server.ProcessNbrToFloodMsg(nbrToFloodData)

		case stop := <-server.FloodData.FloodCtrlCh:
			server.logger.Debug("Flood: Stopping flood channel.", stop)
			server.DeinitFlooding()
			server.FloodData.FloodCtrlReplyCh <- true
			return
		}
	}
}

func (server *OSPFV2Server) StopFlooding() {
	server.FloodData.FloodCtrlCh <- true
	cnt := 0
	for {
		select {
		case _ = <-server.FloodData.FloodCtrlReplyCh:
			server.logger.Info("Successfully Stopped  flooding")
			return
		default:
			time.Sleep(time.Duration(10) * time.Millisecond)
			cnt = cnt + 1
			if cnt == 100 {
				server.logger.Err("Unable to stop the flooding routine")
				return
			}
		}
	}

}

func (server *OSPFV2Server) ProcessNbrToFloodMsg(msg NbrToFloodMsg) {
	switch msg.MsgType {
	case LSA_FLOOD_ALL:
		server.ProcessLsaFloodAll(msg.NbrKey, msg.LsaType, msg.LsaPkt)
	case LSA_FLOOD_INTF:
		server.ProcessLsaFloodIntf(msg.NbrKey, msg.LsaPkt)
	default:
		server.logger.Err("Flood: Invalid flood message type", msg.MsgType)
		return
	}
}

/* Flood incoming LSA to the appropriate interfaces. */
func (server OSPFV2Server) ProcessLsaFloodAll(nbrKey NbrConfKey, lsaType uint8, lsa_pkt []byte) {
	nbrConf := server.NbrConfMap[nbrKey]
	rxIntf := server.IntfConfMap[nbrConf.IntfKey]
	var lsaEncPkt []byte
	for key, intf := range server.IntfConfMap {
		areaid := intf.AreaId
		if intf.IpAddr == rxIntf.IpAddr || areaid != rxIntf.AreaId {
			server.logger.Info(fmt.Sprintln("LSA_FLOOD_ALL:Dont flood on rx intf ", rxIntf.IpAddr))
			continue // dont flood the LSA on the interface it is received.
		}
		send := server.nbrFloodCheck(nbrKey, key, intf, lsaType)
		if send {
			if lsa_pkt != nil {
				server.logger.Info(fmt.Sprintln("LSA_FLOOD_ALL: Unicast LSA interface ", intf.IpAddr))
				lsas_enc := make([]byte, 4)
				var no_lsa uint32
				no_lsa = 1
				binary.BigEndian.PutUint32(lsas_enc, no_lsa)
				lsaEncPkt = append(lsaEncPkt, lsas_enc...)
				lsaEncPkt = append(lsaEncPkt, lsa_pkt...)
				lsa_pkt_len := len(lsaEncPkt)
				destIp := net.ParseIP(convertUint32ToDotNotation(nbrConf.NbrIP))
				pkt := server.BuildLsaUpdPkt(key, intf,
					intf.IfMacAddr, destIp, lsa_pkt_len, lsaEncPkt)
				server.SendOspfPkt(key, pkt)
			}
		}
	}

}

/* For LSA request packets send LSA to that particular interface. */
func (server *OSPFV2Server) ProcessLsaFloodIntf(nbrKey NbrConfKey, lsa_pkt []byte) {
	nbrConf, exists := server.NbrConfMap[nbrKey]
	if !exists {
		server.logger.Info(fmt.Sprintln("Flood: LSA_FLOOD_INTF Neighbor doesnt exist . Dont flood.", nbrKey))
		return
	}
	intf, valid := server.IntfConfMap[nbrConf.IntfKey]
	if !valid {
		server.logger.Err("Flood: LSA_FLOOD_INTF intf not found ", nbrConf.IntfKey)
		return
	}
	var lsaEncPkt []byte
	if lsa_pkt != nil {
		lsas_enc := make([]byte, 4)
		var no_lsa uint32
		no_lsa = 1
		binary.BigEndian.PutUint32(lsas_enc, no_lsa)
		lsaEncPkt = append(lsaEncPkt, lsas_enc...)
		lsaEncPkt = append(lsaEncPkt, lsa_pkt...)
		lsa_pkt_len := len(lsaEncPkt)
		destIp := net.ParseIP(convertUint32ToDotNotation(nbrConf.NbrIP))
		pkt := server.BuildLsaUpdPkt(nbrConf.IntfKey, intf,
			intf.IfMacAddr, destIp, lsa_pkt_len, lsaEncPkt)
		server.logger.Info(fmt.Sprintln("LSA_FLOOD_INTF: Send  LSA to interface ", intf.IpAddr))
		server.SendOspfPkt(nbrConf.IntfKey, pkt)

	}

}

/* Flood self originated LSAs received from LSdb */
func (server *OSPFV2Server) ProcessLsdbToFloodMsg(msg LsdbToFloodLSAMsg) {
	dstMac := net.HardwareAddr{0x01, 0x00, 0x5e, 0x00, 0x00, 0x05}
	dstIp := net.IP{224, 0, 0, 5}

	var lsaEncPkt []byte
	for key, intf := range server.IntfConfMap {
		areaid := intf.AreaId
		if areaid != msg.AreaId {
			server.logger.Info(fmt.Sprintln("LSA_FLOOD_ALL:Dont flood on rx intf ", intf.IpAddr))
			continue // dont flood the LSA on the interface it is received.
		}
		nbrs, valid := server.NbrConfData.IntfToNbrMap[key]
		if !valid {
			server.logger.Debug("Flood: No nbrs exist for intf . No flood ", key)
			continue
		}
		for _, nbrKey := range nbrs {
			send := server.nbrFloodCheck(nbrKey, key, intf, msg.LsaKey.LSType)
			if send {
				_, valid := server.NbrConfMap[nbrKey]
				if !valid {
					server.logger.Debug("Flood: nbr conf does not exist.Dont flood ", nbrKey)
					continue
				}
				if msg.LsaData.([]byte) != nil {
					server.logger.Info(fmt.Sprintln("LSA_FLOOD_ALL: Unicast LSA interface ", intf.IpAddr))
					lsas_enc := make([]byte, 4)
					var no_lsa uint32
					no_lsa = 1
					binary.BigEndian.PutUint32(lsas_enc, no_lsa)
					lsaEncPkt = append(lsaEncPkt, lsas_enc...)
					lsaEncPkt = append(lsaEncPkt, msg.LsaData.([]byte)...)
					lsa_pkt_len := len(lsaEncPkt)
					pkt := server.BuildLsaUpdPkt(key, intf,
						dstMac, dstIp, lsa_pkt_len, lsaEncPkt)
					server.SendOspfPkt(key, pkt)
				}
			} // end of send packet
		} // end of nbrs / intf
	} // end of intf for
	server.logger.Debug("Flood: Lsdb to flood processing done..")
}

/*@fn sendRouterLsa
At the event of interface down need to flood
updated router LSA.
*/
/*@fn constructAndSendLsaAgeFlood
Flood LSAs which reached max age.
*/

/*
func (server *OSPFV2Server) constructAndSendLsaAgeFlood() {
	dstMac := net.HardwareAddr{0x01, 0x00, 0x5e, 0x00, 0x00, 0x05}
	dstIp := net.IP{224, 0, 0, 5}
	lsas_enc := make([]byte, 4)
	var lsaEncPkt []byte
	var lsasWithHeader []byte
	var no_lsa uint32
	no_lsa = 0
	total_len := 0
	for lsaKey, lsaPkt := range maxAgeLsaMap {
		if lsaPkt != nil {
			no_lsa++
			checksumOffset := uint16(14)
			checkSum := computeFletcherChecksum(lsaPkt[2:], checksumOffset)
			binary.BigEndian.PutUint16(lsaPkt[16:18], checkSum)
			pktLen := len(lsaPkt)
			binary.BigEndian.PutUint16(lsaPkt[18:20], uint16(pktLen))
			lsaEncPkt = append(lsaEncPkt, lsaPkt...)
			total_len += pktLen
			server.logger.Info(fmt.Sprintln("FLUSH: Added to flush list lsakey",
				lsaKey.AdvRouter, lsaKey.LSId, lsaKey.LSId))
		}
		msg := maxAgeLsaMsg{
			lsaKey:   lsaKey,
			msg_type: delMaxAgeLsa,
		}
		server.maxAgeLsaCh <- msg
	}
	lsa_pkt_len := total_len + OSPF_NO_OF_LSA_FIELD
	if lsa_pkt_len == OSPF_NO_OF_LSA_FIELD {
		return
	}
	binary.BigEndian.PutUint32(lsas_enc, no_lsa)
	lsasWithHeader = append(lsasWithHeader, lsas_enc...)
	lsasWithHeader = append(lsasWithHeader, lsaEncPkt...)

	for key, intConf := range server.IntfConfMap {
		server.logger.Info(fmt.Sprintln("FLUSH: Send flush message ", intConf.IpAddr))
		pkt := server.BuildLsaUpdPkt(key, intConf,
			dstMac, dstIp, lsa_pkt_len, lsasWithHeader)
		server.SendOspfPkt(key, pkt)
	}

}
*/
/* @fn interfaceFloodCheck
Check if we need to flood the LSA on the interface
*/
func (server *OSPFV2Server) nbrFloodCheck(nbrKey NbrConfKey, key IntfConfKey, intf IntfConf, lsType uint8) bool {
	/* Check neighbor state */
	flood_check := true
	nbrConf := server.NbrConfMap[nbrKey]
	intfConf, valid := server.IntfConfMap[nbrConf.IntfKey]
	if !valid {
		server.logger.Err("Nbr : Intf does not exist. Flood check failed. ", nbrConf.IntfKey)
		return false
	}
	//rtrid := convertIPv4ToUint32(server.globalData.RouterId)
	if nbrConf.IntfKey == key && nbrConf.NbrDR == intfConf.DRtrId && lsType != Summary3LSA && lsType != Summary4LSA {
		server.logger.Info(fmt.Sprintln("IF FLOOD: Nbr is DR/BDR.   flood on this interface . nbr - ", nbrKey.NbrIdentity, nbrConf.NbrIP))
		return false
	}
	flood_check = server.interfaceFloodCheck(key)
	return flood_check
}

func (server *OSPFV2Server) interfaceFloodCheck(key IntfConfKey) bool {
	flood_check := false
	nbrData, exist := server.NbrConfData.IntfToNbrMap[key]
	if !exist {
		server.logger.Info(fmt.Sprintln("FLOOD: Intf to nbr map doesnt exist.Dont flood."))
		return false
	}
	if nbrData != nil {
		for _, nbrId := range nbrData {
			nbrConf := server.NbrConfMap[nbrId]
			if nbrConf.State < NbrExchange {
				server.logger.Info(fmt.Sprintln("FLOOD: Nbr < exchange . ", nbrConf.NbrIP))
				flood_check = false
				continue
			}
			flood_check = true
			/* TODO - add check if nbrstate is loading - check its retransmission list
			   add LSA to the adjacency list of neighbor with FULL state.*/
		}
	} else {
		server.logger.Info(fmt.Sprintln("FLOOD: nbr list is null for interface ", key.IpAddr))
	}
	return flood_check
}

/*
@fn processSummaryLSAFlood
This API takes care of flooding new summary LSAs that is added in the LSDB
*/
func (server *OSPFV2Server) processSummaryLSAFlood(areaId uint32, lsaKey LsaKey) {
	var lsaEncPkt []byte
	LsaEnc := []byte{}

	server.logger.Info(fmt.Sprintln("Summary: Start flooding algorithm. Area ",
		areaId, " lsa ", lsaKey))
	LsaEnc = server.encodeSummaryLsa(areaId, lsaKey)
	no_lsas := uint32(1)
	lsas_enc := make([]byte, 4)
	binary.BigEndian.PutUint32(lsas_enc, no_lsas)
	lsaEncPkt = append(lsaEncPkt, lsas_enc...)
	lsaEncPkt = append(lsaEncPkt, LsaEnc...)
	lsid := lsaKey.LSId
	adv_router := lsaKey.AdvRouter
	server.logger.Info(fmt.Sprintln("SUMMARY: Send for flooding ",
		areaId, " adv_router ", adv_router, " lsid ",
		lsid))
	server.floodSummaryLsa(lsaEncPkt, areaId)
	server.logger.Info(fmt.Sprintln("SUMMARY: End flooding process. lsa", lsaKey))
}

func (server *OSPFV2Server) encodeSummaryLsa(areaid uint32, lsakey LsaKey) []byte {
	entry, ret := server.getSummaryLsaFromLsdb(areaid, lsakey)
	if ret == LsdbEntryNotFound {
		server.logger.Info(fmt.Sprintln("Summary LSA: Lsa not found . Area",
			areaid, " LSA key ", lsakey))
		return nil
	}
	LsaEnc := encodeSummaryLsa(entry, lsakey)
	pktLen := len(LsaEnc)
	checksumOffset := uint16(14)
	checkSum := computeFletcherChecksum(LsaEnc[2:], checksumOffset)
	binary.BigEndian.PutUint16(LsaEnc[16:18], checkSum)
	binary.BigEndian.PutUint16(LsaEnc[18:20], uint16(pktLen))
	return LsaEnc

}

func (server *OSPFV2Server) floodSummaryLsa(pkt []byte, areaid uint32) {
	dstMac := net.HardwareAddr{0x01, 0x00, 0x5e, 0x00, 0x00, 0x05}
	dstIp := net.IP{224, 0, 0, 5}
	for key, _ := range server.IntfConfMap {
		intf, ok := server.IntfConfMap[key]
		if !ok {
			continue
		}
		ifArea := intf.AreaId
		//isStub := server.isStubArea(areaid)
		if ifArea == areaid {
			// flood to your own area
			nbrMdata, ok := server.NbrConfData.IntfToNbrMap[key]
			if ok && len(nbrMdata) > 0 {
				send_pkt := server.BuildLsaUpdPkt(key, intf, dstMac, dstIp, len(pkt), pkt)
				server.logger.Info(fmt.Sprintln("SUMMARY: Send  LSA to interface ", intf.IpAddr, " area ", intf.AreaId))
				server.SendOspfPkt(key, send_pkt)
			}

		}
	}
}

/*
@fn processAsExternalLSAFlood
	This API takes care of flooding external routes through
	AS external LSA
*/
func (server *OSPFV2Server) processAsExternalLSAFlood(lsakey LsaKey) {
	areaId := uint32(0)
	for ent, _ := range server.AreaConfMap {
		areaId = ent
	}
	var lsaEncPkt []byte
	LsaEnc := []byte{}

	entry, ret := server.getASExternalLsaFromLsdb(areaId, lsakey)
	if ret == LsdbEntryNotFound {
		server.logger.Info(fmt.Sprintln("ASBR: Lsa not found . Area",
			areaId, " LSA key ", lsakey))
		return
	}
	LsaEnc = encodeASExternalLsa(entry, lsakey)
	pktLen := len(LsaEnc)
	checksumOffset := uint16(14)
	checkSum := computeFletcherChecksum(LsaEnc[2:], checksumOffset)
	binary.BigEndian.PutUint16(LsaEnc[16:18], checkSum)
	binary.BigEndian.PutUint16(LsaEnc[18:20], uint16(pktLen))

	no_lsas := uint32(1)
	lsas_enc := make([]byte, 4)
	binary.BigEndian.PutUint32(lsas_enc, no_lsas)
	lsaEncPkt = append(lsaEncPkt, lsas_enc...)
	lsaEncPkt = append(lsaEncPkt, LsaEnc...)
	lsid := lsakey.LSId
	adv_router := lsakey.AdvRouter
	server.logger.Info(fmt.Sprintln("ASBR: flood lsid ", lsid, " adv_router ", adv_router))
	server.floodASExternalLsa(lsaEncPkt)
}

func (server *OSPFV2Server) floodASExternalLsa(pkt []byte) {
	server.logger.Info(fmt.Sprintln("ASBR: Received for flood ", pkt))
	dstMac := net.HardwareAddr{0x01, 0x00, 0x5e, 0x00, 0x00, 0x05}
	dstIp := net.IP{224, 0, 0, 5}
	for key, _ := range server.IntfConfMap {
		intf, ok := server.IntfConfMap[key]
		if !ok {
			continue
		}
		isStub, _ := server.isStubArea(intf.AreaId)
		if isStub {
			server.logger.Info(fmt.Sprintln("ASBR: Dont flood AS external as area is stub ", intf.AreaId))
			continue
		}
		nbrMdata, ok := server.NbrConfData.IntfToNbrMap[key]
		if ok && len(nbrMdata) > 0 {
			send_pkt := server.BuildLsaUpdPkt(key, intf, dstMac, dstIp, len(pkt), pkt)
			server.logger.Info(fmt.Sprintln("ASBR: Send  LSA to interface ", intf.IpAddr))
			server.SendOspfPkt(key, send_pkt)
		}
	}
}

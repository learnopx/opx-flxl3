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
package packet

import (
	"asicdInt"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"l3/vrrp/config"
	_ "net"
	_ "time"
)

const (
	// ip/vrrp header Check Defines
	VRRP_TTL                        = 255
	VRRP_PKT_TYPE_ADVERTISEMENT     = 1 // Only one type is supported which is advertisement
	VRRP_RSVD                       = 0
	VRRP_HDR_CREATE_CHECKSUM        = 0
	VRRP_HEADER_SIZE_EXCLUDING_IPVX = 8 // 8 bytes...
	VRRP_IPV4_HEADER_MIN_SIZE       = 20
	VRRP_HEADER_MIN_SIZE            = 20
	VRRP_PROTOCOL_MAC               = "01:00:5e:00:00:12"
	VRRP_GROUP_IP                   = "224.0.0.18"

	// error message from Packet
	VRRP_CHECKSUM_ERR      = "VRRP checksum failure"
	VRRP_INCORRECT_VERSION = "Version is not correct for received VRRP Packet"
	VRRP_INCORRECT_FIELDS  = "Field like type/count ip addr/Advertisement Interval are not valid"
)

/*
Octet Offset--> 0                   1                   2                   3
 |		0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 |		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 V		|                    IPv4 Fields or IPv6 Fields                 |
		...                                                             ...
		|                                                               |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 0		|Version| Type  | Virtual Rtr ID|   Priority    |Count IPvX Addr|
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 4		|(rsvd) |     Max Adver Int     |          Checksum             |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 8		|                                                               |
		+                                                               +
12		|                       IPvX Address(es)                        |
		+                                                               +
..		+                                                               +
		+                                                               +
		+                                                               +
		|                                                               |
		+                                                               +
		|                                                               |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/

/*
 *  VRRP Packet INTERFACE
 */
type Packet interface {
	Decode(gopacket.Packet) *Header
	ValidateHeader(*Header, []byte) error
	Encode(*PacketInfo) []byte
}

type Header struct {
	Version      uint8
	Type         uint8
	VirtualRtrId uint8
	Priority     uint8
	CountIPAddr  uint8
	Rsvd         uint8
	MaxAdverInt  uint16
	CheckSum     uint16
	IpAddr       []net.IP
}

type PacketInfo struct {
	Version      string
	Vrid         uint8
	Priority     uint8
	AdvertiseInt uint16
	VirutalMac   string
	IpAddr       string
}

func Init() *PacketInfo {
	var pktInfo Packet
	pktInfo = &PacketInfo{}
	return pktInfo
}

func computeChecksum(version uint8, content []byte) uint16 {
	var csum uint32
	var rv uint16
	if version == config.VERSION2 {
		for i := 0; i < len(content); i += 2 {
			csum += uint32(content[i]) << 8
			csum += uint32(content[i+1])
		}
		rv = ^uint16((csum >> 16) + csum)
	} else if version == config.VERSION3 {
		//@TODO: .....
	}

	return rv
}

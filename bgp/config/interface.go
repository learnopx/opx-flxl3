package config

/*  Port/Interface state change manager.
 */
type IntfStateMgrIntf interface {
	Start()
	PortStateChange()
	GetIPv4Intfs() []*IntfStateInfo
	GetIPv4Information(ifIndex int32) (string, error)
	GetIfIndex(int, int) int32
}

/*  Adding routes to rib/switch/linux interface
 */
type RouteMgrIntf interface {
	Start()
	GetNextHopInfo(ipAddr string) (*NextHopInfo, error)
	CreateRoute(*RouteConfig)
	DeleteRoute(*RouteConfig)
	ApplyPolicy(protocol string, policy string, action string, conditions []*ConditionInfo)
	GetRoutes() ([]*RouteInfo, []*RouteInfo)
}

/*  Interface for handling policy related operations
 */
type PolicyMgrIntf interface {
	Start()
}

/*  Interface for handling bfd state notifications
 */
type BfdMgrIntf interface {
	Start()
	CreateBfdSession(ipAddr string, sessionParam string) (bool, error)
	DeleteBfdSession(ipAddr string) (bool, error)
}

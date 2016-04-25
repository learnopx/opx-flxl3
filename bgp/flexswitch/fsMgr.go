package FSMgr

import (
	"asicdServices"
	"bfdd"
	"errors"
	"fmt"
	nanomsg "github.com/op/go-nanomsg"
	"l3/bgp/rpc"
	"l3/bgp/server"
	"ribd"
	"utils/logging"
)

/*  Router manager will handle all the communication with ribd
 */
type FSRouteMgr struct {
	ribdClient *ribd.RIBDServicesClient
	plugin     string
	logger     *logging.Writer
}

/*  Interface manager will handle all the communication with asicd
 */
type FSIntfMgr struct {
	plugin               string
	logger               *logging.Writer
	Server               *server.BGPServer // copy of server pointer
	AsicdClient          *asicdServices.ASICDServicesClient
	asicdL3IntfSubSocket *nanomsg.SubSocket
}

/*  @FUTURE: this will be using in future if FlexSwitch is planning to support
 *	     daemon which is handling policy statments
 */
type FSPolicyMgr struct {
	plugin string
	logger *logging.Writer
}

/*  BFD manager will handle all the communication with bfd daemon
 */
type FSBfdMgr struct {
	plugin       string
	logger       *logging.Writer
	Server       *server.BGPServer // copy of server pointer
	bfddClient   *bfdd.BFDDServicesClient
	bfdSubSocket *nanomsg.SubSocket
}

/*  Init policy manager with specific needs
 */
func NewFSPolicyMgr(logger *logging.Writer, fileName string) *FSPolicyMgr {
	mgr := &FSPolicyMgr{
		plugin: "ovsdb",
		logger: logger,
	}

	return mgr
}

/*  Init route manager with ribd client as its core
 */
func NewFSRouteMgr(logger *logging.Writer, fileName string) (*FSRouteMgr, error) {
	var ribdClient *ribd.RIBDServicesClient = nil
	ribdClientChan := make(chan *ribd.RIBDServicesClient)

	logger.Info("Connecting to RIBd")
	go rpc.StartRibdClient(logger, fileName, ribdClientChan)
	ribdClient = <-ribdClientChan
	if ribdClient == nil {
		logger.Err("Failed to connect to RIBd\n")
		return nil, errors.New("Failed to connect to RIBd")
	} else {
		logger.Info("Connected to RIBd")
	}

	mgr := &FSRouteMgr{
		plugin:     "ovsdb",
		ribdClient: ribdClient,
		logger:     logger,
	}

	return mgr, nil
}

func (mgr *FSRouteMgr) CreateRoute() {
	fmt.Println("Create Route called in", mgr.plugin)
}

func (mgr *FSRouteMgr) DeleteRoute() {

}

func (mgr *FSPolicyMgr) AddPolicy() {

}

func (mgr *FSPolicyMgr) RemovePolicy() {

}

func (mgr *FSIntfMgr) PortStateChange() {

}

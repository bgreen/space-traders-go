package sthandler

import (
	"fmt"
	"os"

	"github.com/bgreen/space-traders-go/stapi"

	"context"
	"net/http"
	"reflect"
	"sync"
	"time"
)

type Server struct {
	apiClient *stapi.APIClient
	recvChan  chan Message
	done      chan bool
	limiter   sync.Mutex
	auth      string
	callbacks map[reflect.Type][]callbackInfo
	cbMutex   sync.RWMutex
}

/////////////////////////////
// Server Init
/////////////////////////////

func NewServer() *Server {
	configuration := stapi.NewConfiguration()
	s := Server{
		apiClient: stapi.NewAPIClient(configuration),
		recvChan:  make(chan Message, 10),
		done:      make(chan bool),
		callbacks: make(map[reflect.Type][]callbackInfo),
	}

	s.retrieveAuth()

	return &s
}

func (s *Server) Start() {
	go s.service()
}

func (s *Server) Stop() {
	s.done <- true
}

func (s *Server) retrieveAuth() {
	// Read Bearer token from token.txt
	token, err := os.ReadFile("token.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "No token file found")
		return
	}
	s.auth = string(token)

	// TODO: Check if name is registered
	// TODO: Register name if necessary
	// TODO: Write token.txt file
}

func (s *Server) getAuth() context.Context {
	return context.WithValue(context.Background(), stapi.ContextAccessToken, s.auth)
}

/////////////////////////////
//	Callbacks
/////////////////////////////

type Callback func(r Message)

type callbackInfo struct {
	f    Callback
	once bool
}

func (s *Server) RegisterCallback(r any, f Callback) {
	s.cbMutex.Lock()
	t := reflect.TypeOf(r)
	s.callbacks[t] = append(s.callbacks[t], callbackInfo{f, false})
	s.cbMutex.Unlock()
}

func (s *Server) RegisterCallbackOnce(r any, f Callback) {
	s.cbMutex.Lock()
	t := reflect.TypeOf(r)
	s.callbacks[t] = append(s.callbacks[t], callbackInfo{f, true})
	s.cbMutex.Unlock()
}

func (s *Server) UnregisterCallback(r any, f Callback) {
	/* TODO: can't unregister because functions aren't comparable, therefore can't find it in the list
	s.cbMutex.Lock()
	t := reflect.TypeOf(r)
	// Is there anything registered to this response type?
	if cbs, ok := s.callbacks[t]; ok {
		//s.callbacks[t] = append(s.callbacks[t], f)
		for i, cb := range cbs {
			// Is this the same callback?
			if cb == f {

			}
		}
	}
	s.cbMutex.Unlock()
	*/
}

func NewCallbackOnceChannel[T any](s *Server) chan Message {
	// A channel for replying to the requestor on
	replyC := make(chan Message, 1)

	// An anonymous function that will place a message onto the reply channel
	f := func(r Message) {

		// Ensure the Data type is the one requested
		if _, ok := r.Data.(T); ok {
			replyC <- r
			defer close(replyC)
		}

	}

	// A variable to place the reply in
	var t T

	// Register our anonymous callback function for whenever we get reply of type T
	s.RegisterCallbackOnce(t, f)

	return replyC
}

/////////////////////////////
// Timers
/////////////////////////////

func (s *Server) timerGive() {
	s.limiter.TryLock()
	s.limiter.Unlock()
}

func (s *Server) timerWaitToGive(d time.Duration) {
	time.Sleep(d)
	s.timerGive()
}

func (s *Server) timerTake() {
	s.limiter.Lock()
}

func (s *Server) service() {
	for {
		select {
		case m := <-s.recvChan:
			s.timerTake()
			s.handleMsg(m)
			go s.timerWaitToGive(time.Second / 2)

		case <-s.done:
			s.timerGive()
			return
		}
	}
}

// TODO: Optionally page over longer results
var limit int32 = 10

func (s *Server) handleMsg(m Message) {
	var resp any // The response to Execute
	var err error
	var http *http.Response

	switch req := m.Data.(type) {
	case stapi.ApiGetMyAgentRequest:
		resp, http, err = req.Execute()
	case stapi.ApiAcceptContractRequest:
		resp, http, err = req.Execute()
	case stapi.ApiDeliverContractRequest:
		resp, http, err = req.Execute()
	case stapi.ApiFulfillContractRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetContractRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetContractsRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetStatusRequest:
		resp, http, err = req.Execute()
	case stapi.ApiRegisterRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetFactionRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetFactionsRequest:
		resp, http, err = req.Execute()
	case stapi.ApiCreateChartRequest:
		resp, http, err = req.Execute()
	case stapi.ApiCreateShipShipScanRequest:
		resp, http, err = req.Execute()
	case stapi.ApiCreateShipSystemScanRequest:
		resp, http, err = req.Execute()
	case stapi.ApiCreateShipWaypointScanRequest:
		resp, http, err = req.Execute()
	case stapi.ApiCreateSurveyRequest:
		resp, http, err = req.Execute()
	case stapi.ApiDockShipRequest:
		resp, http, err = req.Execute()
	case stapi.ApiExtractResourcesRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetMountsRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetMyShipRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetMyShipCargoRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetMyShipsRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetShipCooldownRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetShipNavRequest:
		resp, http, err = req.Execute()
	case stapi.ApiInstallMountRequest:
		resp, http, err = req.Execute()
	case stapi.ApiJettisonRequest:
		resp, http, err = req.Execute()
	case stapi.ApiJumpShipRequest:
		resp, http, err = req.Execute()
	case stapi.ApiNavigateShipRequest:
		resp, http, err = req.Execute()
	case stapi.ApiNegotiateContractRequest:
		resp, http, err = req.Execute()
	case stapi.ApiOrbitShipRequest:
		resp, http, err = req.Execute()
	case stapi.ApiPatchShipNavRequest:
		resp, http, err = req.Execute()
	case stapi.ApiPurchaseCargoRequest:
		resp, http, err = req.Execute()
	case stapi.ApiPurchaseShipRequest:
		resp, http, err = req.Execute()
	case stapi.ApiRefuelShipRequest:
		resp, http, err = req.Execute()
	case stapi.ApiRemoveMountRequest:
		resp, http, err = req.Execute()
	case stapi.ApiSellCargoRequest:
		resp, http, err = req.Execute()
	case stapi.ApiShipRefineRequest:
		resp, http, err = req.Execute()
	case stapi.ApiTransferCargoRequest:
		resp, http, err = req.Execute()
	case stapi.ApiWarpShipRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetJumpGateRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetMarketRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetShipyardRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetSystemRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetSystemWaypointsRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetSystemsRequest:
		resp, http, err = req.Execute()
	case stapi.ApiGetWaypointRequest:
		resp, http, err = req.Execute()

	}
	m.Data = resp

	// Handle HTTP errors
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		switch err.Error() {
		case "401 Unauthorized":
			// TODO: Re-register
		}
	}

	// Decode API Errors
	if err != nil {
		apiErr, _ := DecodeApiError(http)
		m.Err = apiErr
	}

	// Execute callbacks
	s.cbMutex.Lock()
	t := reflect.TypeOf(m.Data)
	for i, v := range s.callbacks[t] {

		v.f(m)

		// If this is a run-once callback
		if v.once {
			// Cut the callback out of the slice
			count := len(s.callbacks[t])
			if count == 1 {
				s.callbacks[t] = []callbackInfo{}
			} else if i == count-1 {
				s.callbacks[t] = s.callbacks[t][:i]
			} else {
				s.callbacks[t] = append(s.callbacks[t][:i], s.callbacks[t][i+1:]...)
			}
		}
	}
	s.cbMutex.Unlock()
}

/////////////////////////////
// Request Methods
/////////////////////////////

func SendRequestGetReply[T any](s *Server, d any) Message {

	c := NewCallbackOnceChannel[T](s)
	s.recvChan <- NewMessage(d)

	m := <-c

	return m
}

func (s *Server) GetMyAgent() (stapi.Agent, error) {
	m := SendRequestGetReply[*stapi.GetMyAgent200Response](s, s.apiClient.AgentsApi.GetMyAgent(s.getAuth()))
	return m.Data.(*stapi.GetMyAgent200Response).GetData(), m.Err
}

func (s *Server) AcceptContract(contract string) (stapi.AcceptContract200ResponseData, error) {
	m := SendRequestGetReply[*stapi.AcceptContract200Response](s, s.apiClient.ContractsApi.AcceptContract(s.getAuth(), contract))
	return m.Data.(*stapi.AcceptContract200Response).GetData(), m.Err
}

func (s *Server) DeliverContract(contract string, ship string, trade stapi.TradeSymbol, units int32) (stapi.DeliverContract200ResponseData, error) {
	request := *stapi.NewDeliverContractRequest(ship, string(trade), units)
	m := SendRequestGetReply[*stapi.DeliverContract200Response](s, s.apiClient.ContractsApi.DeliverContract(s.getAuth(), contract).DeliverContractRequest(request))
	return m.Data.(*stapi.DeliverContract200Response).GetData(), m.Err
}

func (s *Server) FulfillContract(contract string) (stapi.AcceptContract200ResponseData, error) {
	m := SendRequestGetReply[*stapi.FulfillContract200Response](s, s.apiClient.ContractsApi.FulfillContract(s.getAuth(), contract))
	return m.Data.(*stapi.FulfillContract200Response).GetData(), m.Err
}

func (s *Server) GetContract(contract string) (stapi.Contract, error) {
	m := SendRequestGetReply[*stapi.GetContract200Response](s, s.apiClient.ContractsApi.GetContract(s.getAuth(), contract))
	return m.Data.(*stapi.GetContract200Response).GetData(), m.Err
}

func (s *Server) GetContracts() ([]stapi.Contract, error) {
	m := SendRequestGetReply[*stapi.GetContracts200Response](s, s.apiClient.ContractsApi.GetContracts(s.getAuth()).Limit(limit))
	return m.Data.(*stapi.GetContracts200Response).GetData(), m.Err
}

func (s *Server) GetStatus() (stapi.GetStatus200Response, error) {
	m := SendRequestGetReply[*stapi.GetStatus200Response](s, s.apiClient.DefaultApi.GetStatus(s.getAuth()))
	return *m.Data.(*stapi.GetStatus200Response), m.Err
}

func (s *Server) Register(name string, faction string) (stapi.Register201ResponseData, error) {
	request := *stapi.NewRegisterRequest(stapi.FactionSymbols(faction), name)
	m := SendRequestGetReply[*stapi.Register201Response](s, s.apiClient.DefaultApi.Register(s.getAuth()).RegisterRequest(request))
	return m.Data.(*stapi.Register201Response).GetData(), m.Err
}

func (s *Server) GetFaction(faction stapi.FactionSymbols) (stapi.Faction, error) {
	m := SendRequestGetReply[*stapi.GetFaction200Response](s, s.apiClient.FactionsApi.GetFaction(s.getAuth(), string(faction)))
	return m.Data.(*stapi.GetFaction200Response).GetData(), m.Err
}

func (s *Server) GetFactions() ([]stapi.Faction, error) {
	m := SendRequestGetReply[*stapi.GetFactions200Response](s, s.apiClient.FactionsApi.GetFactions(s.getAuth()).Limit(limit))
	return m.Data.(*stapi.GetFactions200Response).GetData(), m.Err
}

func (s *Server) CreateChart(ship string) (stapi.CreateChart201ResponseData, error) {
	m := SendRequestGetReply[*stapi.CreateChart201Response](s, s.apiClient.FleetApi.CreateChart(s.getAuth(), ship))
	return m.Data.(*stapi.CreateChart201Response).GetData(), m.Err
}

func (s *Server) CreateShipShipScan(ship string) (stapi.CreateShipShipScan201ResponseData, error) {
	m := SendRequestGetReply[*stapi.CreateShipShipScan201Response](s, s.apiClient.FleetApi.CreateShipShipScan(s.getAuth(), ship))
	return m.Data.(*stapi.CreateShipShipScan201Response).GetData(), m.Err
}

func (s *Server) CreateShipSystemScan(ship string) (stapi.CreateShipSystemScan201ResponseData, error) {
	m := SendRequestGetReply[*stapi.CreateShipSystemScan201Response](s, s.apiClient.FleetApi.CreateShipSystemScan(s.getAuth(), ship))
	return m.Data.(*stapi.CreateShipSystemScan201Response).GetData(), m.Err
}

func (s *Server) CreateShipWaypointScan(ship string) (stapi.CreateShipWaypointScan201ResponseData, error) {
	m := SendRequestGetReply[*stapi.CreateShipWaypointScan201Response](s, s.apiClient.FleetApi.CreateShipWaypointScan(s.getAuth(), ship))
	return m.Data.(*stapi.CreateShipWaypointScan201Response).GetData(), m.Err
}

func (s *Server) CreateSurvey(ship string) (stapi.CreateSurvey201ResponseData, error) {
	m := SendRequestGetReply[*stapi.CreateSurvey201Response](s, s.apiClient.FleetApi.CreateSurvey(s.getAuth(), ship))
	return m.Data.(*stapi.CreateSurvey201Response).GetData(), m.Err
}

func (s *Server) DockShip(ship string) (stapi.OrbitShip200ResponseData, error) {
	m := SendRequestGetReply[*stapi.DockShip200Response](s, s.apiClient.FleetApi.DockShip(s.getAuth(), ship))
	return m.Data.(*stapi.DockShip200Response).GetData(), m.Err
}

func (s *Server) ExtractResources(ship string) (stapi.ExtractResources201ResponseData, error) {
	request := *stapi.NewExtractResourcesRequest()
	m := SendRequestGetReply[*stapi.ExtractResources201Response](s, s.apiClient.FleetApi.ExtractResources(s.getAuth(), ship).ExtractResourcesRequest(request))
	return m.Data.(*stapi.ExtractResources201Response).GetData(), m.Err
}

func (s *Server) GetMounts(ship string) ([]stapi.ShipMount, error) {
	m := SendRequestGetReply[*stapi.GetMounts200Response](s, s.apiClient.FleetApi.GetMounts(s.getAuth(), ship))
	return m.Data.(*stapi.GetMounts200Response).GetData(), m.Err
}

func (s *Server) GetMyShip(ship string) (stapi.Ship, error) {
	m := SendRequestGetReply[*stapi.GetMyShip200Response](s, s.apiClient.FleetApi.GetMyShip(s.getAuth(), ship))
	return m.Data.(*stapi.GetMyShip200Response).GetData(), m.Err
}

func (s *Server) GetMyShipCargo(ship string) (stapi.ShipCargo, error) {
	m := SendRequestGetReply[*stapi.GetMyShipCargo200Response](s, s.apiClient.FleetApi.GetMyShipCargo(s.getAuth(), ship))
	return m.Data.(*stapi.GetMyShipCargo200Response).GetData(), m.Err
}

func (s *Server) GetMyShips() ([]stapi.Ship, error) {
	m := SendRequestGetReply[*stapi.GetMyShips200Response](s, s.apiClient.FleetApi.GetMyShips(s.getAuth()).Limit(limit))

	return m.Data.(*stapi.GetMyShips200Response).GetData(), m.Err
}

func (s *Server) GetShipCooldown(ship string) (stapi.Cooldown, error) {
	m := SendRequestGetReply[*stapi.GetShipCooldown200Response](s, s.apiClient.FleetApi.GetShipCooldown(s.getAuth(), ship))
	return m.Data.(*stapi.GetShipCooldown200Response).GetData(), m.Err
}

func (s *Server) GetShipNav(ship string) (stapi.ShipNav, error) {
	m := SendRequestGetReply[*stapi.GetShipNav200Response](s, s.apiClient.FleetApi.GetShipNav(s.getAuth(), ship))
	return m.Data.(*stapi.GetShipNav200Response).GetData(), m.Err
}

func (s *Server) InstallMount(ship string, mount string) (stapi.InstallMount201ResponseData, error) {
	request := *stapi.NewInstallMountRequest(mount)
	m := SendRequestGetReply[*stapi.InstallMount201Response](s, s.apiClient.FleetApi.InstallMount(s.getAuth(), ship).InstallMountRequest(request))
	return m.Data.(*stapi.InstallMount201Response).GetData(), m.Err
}

func (s *Server) Jettison(ship string, trade stapi.TradeSymbol, units int32) (stapi.Jettison200ResponseData, error) {
	request := *stapi.NewJettisonRequest(trade, units)
	m := SendRequestGetReply[*stapi.Jettison200Response](s, s.apiClient.FleetApi.Jettison(s.getAuth(), ship).JettisonRequest(request))
	return m.Data.(*stapi.Jettison200Response).GetData(), m.Err
}

func (s *Server) JumpShip(ship string, system string) (stapi.JumpShip200ResponseData, error) {
	request := *stapi.NewJumpShipRequest(system)
	m := SendRequestGetReply[*stapi.JumpShip200Response](s, s.apiClient.FleetApi.JumpShip(s.getAuth(), ship).JumpShipRequest(request))
	return m.Data.(*stapi.JumpShip200Response).GetData(), m.Err
}

func (s *Server) NavigateShip(ship string, waypoint string) (stapi.NavigateShip200ResponseData, error) {
	request := *stapi.NewNavigateShipRequest(waypoint)
	m := SendRequestGetReply[*stapi.NavigateShip200Response](s, s.apiClient.FleetApi.NavigateShip(s.getAuth(), ship).NavigateShipRequest(request))
	return m.Data.(*stapi.NavigateShip200Response).GetData(), m.Err
}

func (s *Server) NegotiateContract(ship string) (stapi.NegotiateContract200ResponseData, error) {
	m := SendRequestGetReply[*stapi.NegotiateContract200Response](s, s.apiClient.FleetApi.NegotiateContract(s.getAuth(), ship))
	return m.Data.(*stapi.NegotiateContract200Response).GetData(), m.Err
}

func (s *Server) OrbitShip(ship string) (stapi.OrbitShip200ResponseData, error) {
	m := SendRequestGetReply[*stapi.OrbitShip200Response](s, s.apiClient.FleetApi.OrbitShip(s.getAuth(), ship))
	return m.Data.(*stapi.OrbitShip200Response).GetData(), m.Err
}

func (s *Server) PatchShipNav(ship string) (stapi.ShipNav, error) {
	request := *stapi.NewPatchShipNavRequest()
	m := SendRequestGetReply[*stapi.GetShipNav200Response](s, s.apiClient.FleetApi.PatchShipNav(s.getAuth(), ship).PatchShipNavRequest(request))
	return m.Data.(*stapi.GetShipNav200Response).GetData(), m.Err
}

func (s *Server) PurchaseCargo(ship string, trade stapi.TradeSymbol, units int32) (stapi.SellCargo201ResponseData, error) {
	request := *stapi.NewPurchaseCargoRequest(trade, units)
	m := SendRequestGetReply[*stapi.PurchaseCargo201Response](s, s.apiClient.FleetApi.PurchaseCargo(s.getAuth(), ship).PurchaseCargoRequest(request))
	return m.Data.(*stapi.PurchaseCargo201Response).GetData(), m.Err
}

func (s *Server) PurchaseShip(ship stapi.ShipType, waypoint string) (stapi.PurchaseShip201ResponseData, error) {
	request := *stapi.NewPurchaseShipRequest(ship, waypoint)
	m := SendRequestGetReply[*stapi.PurchaseShip201Response](s, s.apiClient.FleetApi.PurchaseShip(s.getAuth()).PurchaseShipRequest(request))
	return m.Data.(*stapi.PurchaseShip201Response).GetData(), m.Err
}

func (s *Server) RefuelShip(ship string) (stapi.RefuelShip200ResponseData, error) {
	request := *stapi.NewRefuelShipRequest()
	m := SendRequestGetReply[*stapi.RefuelShip200Response](s, s.apiClient.FleetApi.RefuelShip(s.getAuth(), ship).RefuelShipRequest(request))
	return m.Data.(*stapi.RefuelShip200Response).GetData(), m.Err
}

func (s *Server) RemoveMount(ship string, mount string) (stapi.RemoveMount201ResponseData, error) {
	request := *stapi.NewRemoveMountRequest(mount)
	m := SendRequestGetReply[*stapi.RemoveMount201Response](s, s.apiClient.FleetApi.RemoveMount(s.getAuth(), ship).RemoveMountRequest(request))
	return m.Data.(*stapi.RemoveMount201Response).GetData(), m.Err
}

func (s *Server) SellCargo(ship string, trade stapi.TradeSymbol, units int32) (stapi.SellCargo201ResponseData, error) {
	request := *stapi.NewSellCargoRequest(trade, units)
	m := SendRequestGetReply[*stapi.SellCargo201Response](s, s.apiClient.FleetApi.SellCargo(s.getAuth(), ship).SellCargoRequest(request))
	return m.Data.(*stapi.SellCargo201Response).GetData(), m.Err
}

func (s *Server) ShipRefine(ship string, produce string) (stapi.ShipRefine201ResponseData, error) {
	request := *stapi.NewShipRefineRequest(produce)
	m := SendRequestGetReply[*stapi.ShipRefine201Response](s, s.apiClient.FleetApi.ShipRefine(s.getAuth(), ship).ShipRefineRequest(request))
	return m.Data.(*stapi.ShipRefine201Response).GetData(), m.Err
}

func (s *Server) TransferCargo(shipFrom string, trade stapi.TradeSymbol, units int32, shipTo string) (stapi.Jettison200ResponseData, error) {
	request := *stapi.NewTransferCargoRequest(trade, units, shipTo)
	m := SendRequestGetReply[*stapi.TransferCargo200Response](s, s.apiClient.FleetApi.TransferCargo(s.getAuth(), shipFrom).TransferCargoRequest(request))
	return m.Data.(*stapi.TransferCargo200Response).GetData(), m.Err
}

func (s *Server) WarpShip(ship string, waypoint string) (stapi.NavigateShip200ResponseData, error) {
	request := *stapi.NewNavigateShipRequest(waypoint)
	m := SendRequestGetReply[*stapi.NavigateShip200Response](s, s.apiClient.FleetApi.WarpShip(s.getAuth(), ship).NavigateShipRequest(request))
	return m.Data.(*stapi.NavigateShip200Response).GetData(), m.Err
}

func (s *Server) GetJumpGate(system string, waypoint string) (stapi.JumpGate, error) {
	m := SendRequestGetReply[*stapi.GetJumpGate200Response](s, s.apiClient.SystemsApi.GetJumpGate(s.getAuth(), system, waypoint))
	return m.Data.(*stapi.GetJumpGate200Response).GetData(), m.Err
}

func (s *Server) GetMarket(system string, waypoint string) (stapi.Market, error) {
	m := SendRequestGetReply[*stapi.GetMarket200Response](s, s.apiClient.SystemsApi.GetMarket(s.getAuth(), system, waypoint))
	return m.Data.(*stapi.GetMarket200Response).GetData(), m.Err
}

func (s *Server) GetShipyard(system string, waypoint string) (stapi.Shipyard, error) {
	m := SendRequestGetReply[*stapi.GetShipyard200Response](s, s.apiClient.SystemsApi.GetShipyard(s.getAuth(), system, waypoint))
	return m.Data.(*stapi.GetShipyard200Response).GetData(), m.Err
}

func (s *Server) GetSystem(system string) (stapi.System, error) {
	m := SendRequestGetReply[*stapi.GetSystem200Response](s, s.apiClient.SystemsApi.GetSystem(s.getAuth(), system))
	return m.Data.(*stapi.GetSystem200Response).GetData(), m.Err
}

func (s *Server) GetSystemWaypoints(system string) ([]stapi.Waypoint, error) {
	m := SendRequestGetReply[*stapi.GetSystemWaypoints200Response](s, s.apiClient.SystemsApi.GetSystemWaypoints(s.getAuth(), system).Limit(limit))
	return m.Data.(*stapi.GetSystemWaypoints200Response).GetData(), m.Err
}

func (s *Server) GetSystems() ([]stapi.System, error) {
	m := SendRequestGetReply[*stapi.GetSystems200Response](s, s.apiClient.SystemsApi.GetSystems(s.getAuth()).Limit(limit))
	return m.Data.(*stapi.GetSystems200Response).GetData(), m.Err
}

func (s *Server) GetWaypoint(system string, waypoint string) (stapi.Waypoint, error) {
	m := SendRequestGetReply[*stapi.GetWaypoint200Response](s, s.apiClient.SystemsApi.GetWaypoint(s.getAuth(), system, waypoint))
	return m.Data.(*stapi.GetWaypoint200Response).GetData(), m.Err
}

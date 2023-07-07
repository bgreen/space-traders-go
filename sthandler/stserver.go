package sthandler

import (
	"fmt"
	"os"

	"github.com/bgreen/space-traders-go/stapi"

	"context"
	"sync"
	"time"
)

type Server struct {
	apiClient *stapi.APIClient
	recvChan  chan Message
	done      chan bool
	limiter   sync.Mutex
	auth      string
	name      string
	faction   stapi.FactionSymbols
}

func NewServer(name string) *Server {
	configuration := stapi.NewConfiguration()
	s := Server{
		apiClient: stapi.NewAPIClient(configuration),
		recvChan:  make(chan Message, 10),
		done:      make(chan bool),
		name:      name,
		faction:   stapi.FactionSymbols("COSMIC"),
	}

	s.retrieveAuth()

	return &s
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

func (s *Server) Start() {
	go s.service()
}

func (s *Server) Stop() {
	s.done <- true
}

func (s *Server) timerGive() {
	s.limiter.TryLock()
	s.limiter.Unlock()
}

func (s *Server) timerTake() {
	s.limiter.Lock()
}

func (s *Server) service() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.timerGive()
		case <-s.done:
			return
		}
	}
}

func (s *Server) getAuth() context.Context {
	return context.WithValue(context.Background(), stapi.ContextAccessToken, s.auth)
}

// TODO: Optionally page over longer results
var limit int32 = 10

func (s *Server) GetMyAgent() (stapi.Agent, error) {
	s.timerTake()
	resp, _, err := s.apiClient.AgentsApi.GetMyAgent(s.getAuth()).Execute()
	return resp.GetData(), err
}

func (s *Server) AcceptContract(contract string) (stapi.AcceptContract200ResponseData, error) {
	s.timerTake()
	resp, _, err := s.apiClient.ContractsApi.AcceptContract(s.getAuth(), contract).Execute()
	return resp.GetData(), err
}

func (s *Server) DeliverContract(contract string, ship string, trade stapi.TradeSymbol, units int32) (stapi.DeliverContract200ResponseData, error) {
	s.timerTake()
	request := *stapi.NewDeliverContractRequest(ship, string(trade), units)
	resp, _, err := s.apiClient.ContractsApi.DeliverContract(s.getAuth(), contract).DeliverContractRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) FulfillContract(contract string) (stapi.AcceptContract200ResponseData, error) {
	s.timerTake()
	resp, _, err := s.apiClient.ContractsApi.FulfillContract(s.getAuth(), contract).Execute()
	return resp.GetData(), err
}

func (s *Server) GetContract(contract string) (stapi.Contract, error) {
	s.timerTake()
	resp, _, err := s.apiClient.ContractsApi.GetContract(s.getAuth(), contract).Execute()
	return resp.GetData(), err
}

func (s *Server) GetContracts() ([]stapi.Contract, error) {
	s.timerTake()
	resp, _, err := s.apiClient.ContractsApi.GetContracts(s.getAuth()).Limit(limit).Execute()
	return resp.GetData(), err
}

func (s *Server) GetStatus() (stapi.GetStatus200Response, error) {
	s.timerTake()
	resp, _, err := s.apiClient.DefaultApi.GetStatus(s.getAuth()).Execute()
	return *resp, err
}

func (s *Server) Register() (stapi.Register201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewRegisterRequest(s.faction, s.name)
	resp, _, err := s.apiClient.DefaultApi.Register(s.getAuth()).RegisterRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) GetFaction(faction stapi.FactionSymbols) (stapi.Faction, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FactionsApi.GetFaction(s.getAuth(), string(faction)).Execute()
	return resp.GetData(), err
}

func (s *Server) GetFactions() ([]stapi.Faction, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FactionsApi.GetFactions(s.getAuth()).Limit(limit).Execute()
	return resp.GetData(), err
}

func (s *Server) CreateChart(ship string) (stapi.CreateChart201ResponseData, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FleetApi.CreateChart(s.getAuth(), ship).Execute()
	return resp.GetData(), err
}

func (s *Server) CreateShipShipScan(ship string) (stapi.CreateShipShipScan201ResponseData, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FleetApi.CreateShipShipScan(s.getAuth(), ship).Execute()
	return resp.GetData(), err
}

func (s *Server) CreateShipSystemScan(ship string) (stapi.CreateShipSystemScan201ResponseData, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FleetApi.CreateShipSystemScan(s.getAuth(), ship).Execute()
	return resp.GetData(), err
}

func (s *Server) CreateShipWaypointScan(ship string) (stapi.CreateShipWaypointScan201ResponseData, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FleetApi.CreateShipWaypointScan(s.getAuth(), ship).Execute()
	return resp.GetData(), err
}

func (s *Server) CreateSurvey(ship string) (stapi.CreateSurvey201ResponseData, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FleetApi.CreateSurvey(s.getAuth(), ship).Execute()
	return resp.GetData(), err
}

func (s *Server) DockShip(ship string) (stapi.OrbitShip200ResponseData, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FleetApi.DockShip(s.getAuth(), ship).Execute()
	return resp.GetData(), err
}

func (s *Server) ExtractResources(ship string) (stapi.ExtractResources201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewExtractResourcesRequest()
	resp, _, err := s.apiClient.FleetApi.ExtractResources(s.getAuth(), ship).ExtractResourcesRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) GetMounts(ship string) ([]stapi.ShipMount, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FleetApi.GetMounts(s.getAuth(), ship).Execute()
	return resp.GetData(), err
}

func (s *Server) GetMyShip(ship string) (stapi.Ship, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FleetApi.GetMyShip(s.getAuth(), ship).Execute()
	return resp.GetData(), err
}

func (s *Server) GetMyShipCargo(ship string) (stapi.ShipCargo, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FleetApi.GetMyShipCargo(s.getAuth(), ship).Execute()
	return resp.GetData(), err
}

func (s *Server) GetMyShips() ([]stapi.Ship, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FleetApi.GetMyShips(s.getAuth()).Limit(limit).Execute()
	return resp.GetData(), err
}

func (s *Server) GetShipCooldown(ship string) (stapi.Cooldown, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FleetApi.GetShipCooldown(s.getAuth(), ship).Execute()
	return resp.GetData(), err
}

func (s *Server) GetShipNav(ship string) (stapi.ShipNav, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FleetApi.GetShipNav(s.getAuth(), ship).Execute()
	return resp.GetData(), err
}

func (s *Server) InstallMount(ship string, mount string) (stapi.InstallMount201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewInstallMountRequest(mount)
	resp, _, err := s.apiClient.FleetApi.InstallMount(s.getAuth(), ship).InstallMountRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) Jettison(ship string, trade stapi.TradeSymbol, units int32) (stapi.Jettison200ResponseData, error) {
	s.timerTake()
	request := *stapi.NewJettisonRequest(trade, units)
	resp, _, err := s.apiClient.FleetApi.Jettison(s.getAuth(), ship).JettisonRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) JumpShip(ship string, system string) (stapi.JumpShip200ResponseData, error) {
	s.timerTake()
	request := *stapi.NewJumpShipRequest(system)
	resp, _, err := s.apiClient.FleetApi.JumpShip(s.getAuth(), ship).JumpShipRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) NavigateShip(ship string, waypoint string) (stapi.NavigateShip200ResponseData, error) {
	s.timerTake()
	request := *stapi.NewNavigateShipRequest(waypoint)
	resp, _, err := s.apiClient.FleetApi.NavigateShip(s.getAuth(), ship).NavigateShipRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) NegotiateContract(ship string) (stapi.NegotiateContract200ResponseData, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FleetApi.NegotiateContract(s.getAuth(), ship).Execute()
	return resp.GetData(), err
}

func (s *Server) OrbitShip(ship string) (stapi.OrbitShip200ResponseData, error) {
	s.timerTake()
	resp, _, err := s.apiClient.FleetApi.OrbitShip(s.getAuth(), ship).Execute()
	return resp.GetData(), err
}

func (s *Server) PatchShipNav(ship string) (stapi.ShipNav, error) {
	s.timerTake()
	request := *stapi.NewPatchShipNavRequest()
	resp, _, err := s.apiClient.FleetApi.PatchShipNav(s.getAuth(), ship).PatchShipNavRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) PurchaseCargo(ship string, trade stapi.TradeSymbol, units int32) (stapi.SellCargo201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewPurchaseCargoRequest(trade, units)
	resp, _, err := s.apiClient.FleetApi.PurchaseCargo(s.getAuth(), ship).PurchaseCargoRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) PurchaseShip(ship stapi.ShipType, waypoint string) (stapi.PurchaseShip201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewPurchaseShipRequest(ship, waypoint)
	resp, _, err := s.apiClient.FleetApi.PurchaseShip(s.getAuth()).PurchaseShipRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) RefuelShip(ship string) (stapi.RefuelShip200ResponseData, error) {
	s.timerTake()
	request := *stapi.NewRefuelShipRequest()
	resp, _, err := s.apiClient.FleetApi.RefuelShip(s.getAuth(), ship).RefuelShipRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) RemoveMount(ship string, mount string) (stapi.RemoveMount201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewRemoveMountRequest(mount)
	resp, _, err := s.apiClient.FleetApi.RemoveMount(s.getAuth(), ship).RemoveMountRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) SellCargo(ship string, trade stapi.TradeSymbol, units int32) (stapi.SellCargo201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewSellCargoRequest(trade, units)
	resp, _, err := s.apiClient.FleetApi.SellCargo(s.getAuth(), ship).SellCargoRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) ShipRefine(ship string, produce string) (stapi.ShipRefine201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewShipRefineRequest(produce)
	resp, _, err := s.apiClient.FleetApi.ShipRefine(s.getAuth(), ship).ShipRefineRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) TransferCargo(shipFrom string, trade stapi.TradeSymbol, units int32, shipTo string) (stapi.Jettison200ResponseData, error) {
	s.timerTake()
	request := *stapi.NewTransferCargoRequest(trade, units, shipTo)
	resp, _, err := s.apiClient.FleetApi.TransferCargo(s.getAuth(), shipFrom).TransferCargoRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) WarpShip(ship string, waypoint string) (stapi.NavigateShip200ResponseData, error) {
	s.timerTake()
	request := *stapi.NewNavigateShipRequest(waypoint)
	resp, _, err := s.apiClient.FleetApi.WarpShip(s.getAuth(), ship).NavigateShipRequest(request).Execute()
	return resp.GetData(), err
}

func (s *Server) GetJumpGate(system string, waypoint string) (stapi.JumpGate, error) {
	s.timerTake()
	resp, _, err := s.apiClient.SystemsApi.GetJumpGate(s.getAuth(), system, waypoint).Execute()
	return resp.GetData(), err
}

func (s *Server) GetMarket(system string, waypoint string) (stapi.Market, error) {
	s.timerTake()
	resp, _, err := s.apiClient.SystemsApi.GetMarket(s.getAuth(), system, waypoint).Execute()
	return resp.GetData(), err
}

func (s *Server) GetShipyard(system string, waypoint string) (stapi.Shipyard, error) {
	s.timerTake()
	resp, _, err := s.apiClient.SystemsApi.GetShipyard(s.getAuth(), system, waypoint).Execute()
	return resp.GetData(), err
}

func (s *Server) GetSystem(system string) (stapi.System, error) {
	s.timerTake()
	resp, _, err := s.apiClient.SystemsApi.GetSystem(s.getAuth(), system).Execute()
	return resp.GetData(), err
}

func (s *Server) GetSystemWaypoints(system string) ([]stapi.Waypoint, error) {
	s.timerTake()
	resp, _, err := s.apiClient.SystemsApi.GetSystemWaypoints(s.getAuth(), system).Limit(limit).Execute()
	return resp.GetData(), err
}

func (s *Server) GetSystems() ([]stapi.System, error) {
	s.timerTake()
	resp, _, err := s.apiClient.SystemsApi.GetSystems(s.getAuth()).Limit(limit).Execute()
	return resp.GetData(), err
}

func (s *Server) GetWaypoint(system string, waypoint string) (stapi.Waypoint, error) {
	s.timerTake()
	resp, _, err := s.apiClient.SystemsApi.GetWaypoint(s.getAuth(), system, waypoint).Execute()
	return resp.GetData(), err
}
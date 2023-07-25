package st

import (
	"fmt"
	"math"
	"net/http"
	"os"

	"github.com/bgreen/space-traders-go/stapi"

	"context"
	"sync"
	"time"
)

type Client struct {
	apiClient *stapi.APIClient
	recvChan  chan Message
	done      chan bool
	limiter   sync.Mutex
	auth      string
}

func NewClient() *Client {
	configuration := stapi.NewConfiguration()
	s := Client{
		apiClient: stapi.NewAPIClient(configuration),
		recvChan:  make(chan Message, 10),
		done:      make(chan bool),
	}

	s.retrieveAuth()

	return &s
}

func (s *Client) retrieveAuth() {
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

func (s *Client) Start() {
	go s.service()
}

func (s *Client) Stop() {
	s.done <- true
}

func (s *Client) timerGive() {
	s.limiter.TryLock()
	s.limiter.Unlock()
}

func (s *Client) timerTake() {
	s.limiter.Lock()
}

func (s *Client) service() {
	ticker := time.NewTicker(time.Second / 2)
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

func (s *Client) getAuth() context.Context {
	return context.WithValue(context.Background(), stapi.ContextAccessToken, s.auth)
}

func handleErr(r *http.Response, err error) error {
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
		apiErr, _ := DecodeApiError(r)
		fmt.Fprintln(os.Stderr, apiErr.Message)
	}

	return err
}

// TODO: Optionally page over longer results
var pageSize int32 = 20

func (s *Client) GetMyAgent() (stapi.Agent, error) {
	s.timerTake()
	resp, r, err := s.apiClient.AgentsApi.GetMyAgent(s.getAuth()).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) AcceptContract(contract string) (stapi.AcceptContract200ResponseData, error) {
	s.timerTake()
	resp, r, err := s.apiClient.ContractsApi.AcceptContract(s.getAuth(), contract).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) DeliverContract(contract string, ship string, trade stapi.TradeSymbol, units int32) (stapi.DeliverContract200ResponseData, error) {
	s.timerTake()
	request := *stapi.NewDeliverContractRequest(ship, string(trade), units)
	resp, r, err := s.apiClient.ContractsApi.DeliverContract(s.getAuth(), contract).DeliverContractRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) FulfillContract(contract string) (stapi.AcceptContract200ResponseData, error) {
	s.timerTake()
	resp, r, err := s.apiClient.ContractsApi.FulfillContract(s.getAuth(), contract).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetContract(contract string) (stapi.Contract, error) {
	s.timerTake()
	resp, r, err := s.apiClient.ContractsApi.GetContract(s.getAuth(), contract).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetContracts() ([]stapi.Contract, error) {
	s.timerTake()
	// TODO: loop to get all of them
	resp, r, err := s.apiClient.ContractsApi.GetContracts(s.getAuth()).Limit(pageSize).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetStatus() (stapi.GetStatus200Response, error) {
	s.timerTake()
	resp, r, err := s.apiClient.DefaultApi.GetStatus(s.getAuth()).Execute()
	return *resp, handleErr(r, err)
}

func (s *Client) Register(name string, faction string) (stapi.Register201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewRegisterRequest(stapi.FactionSymbols(faction), name)
	resp, r, err := s.apiClient.DefaultApi.Register(context.Background()).RegisterRequest(request).Execute()
	if err != nil {
		fmt.Println(r)
	}
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetFaction(faction stapi.FactionSymbols) (stapi.Faction, error) {
	s.timerTake()
	resp, r, err := s.apiClient.FactionsApi.GetFaction(s.getAuth(), string(faction)).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetFactions() ([]stapi.Faction, error) {
	s.timerTake()
	// TODO: loop to get all of them
	resp, r, err := s.apiClient.FactionsApi.GetFactions(s.getAuth()).Limit(pageSize).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) CreateChart(ship string) (stapi.CreateChart201ResponseData, error) {
	s.timerTake()
	resp, r, err := s.apiClient.FleetApi.CreateChart(s.getAuth(), ship).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) CreateShipShipScan(ship string) (stapi.CreateShipShipScan201ResponseData, error) {
	s.timerTake()
	resp, r, err := s.apiClient.FleetApi.CreateShipShipScan(s.getAuth(), ship).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) CreateShipSystemScan(ship string) (stapi.CreateShipSystemScan201ResponseData, error) {
	s.timerTake()
	resp, r, err := s.apiClient.FleetApi.CreateShipSystemScan(s.getAuth(), ship).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) CreateShipWaypointScan(ship string) (stapi.CreateShipWaypointScan201ResponseData, error) {
	s.timerTake()
	resp, r, err := s.apiClient.FleetApi.CreateShipWaypointScan(s.getAuth(), ship).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) CreateSurvey(ship string) (stapi.CreateSurvey201ResponseData, error) {
	s.timerTake()
	resp, r, err := s.apiClient.FleetApi.CreateSurvey(s.getAuth(), ship).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) DockShip(ship string) (stapi.OrbitShip200ResponseData, error) {
	s.timerTake()
	resp, r, err := s.apiClient.FleetApi.DockShip(s.getAuth(), ship).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) ExtractResources(ship string) (stapi.ExtractResources201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewExtractResourcesRequest()
	resp, r, err := s.apiClient.FleetApi.ExtractResources(s.getAuth(), ship).ExtractResourcesRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) ExtractSurveyResources(ship string, survey stapi.Survey) (stapi.ExtractResources201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewExtractResourcesRequest()
	request.SetSurvey(survey)
	resp, r, err := s.apiClient.FleetApi.ExtractResources(s.getAuth(), ship).ExtractResourcesRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetMounts(ship string) ([]stapi.ShipMount, error) {
	s.timerTake()
	resp, r, err := s.apiClient.FleetApi.GetMounts(s.getAuth(), ship).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetMyShip(ship string) (stapi.Ship, error) {
	s.timerTake()
	resp, r, err := s.apiClient.FleetApi.GetMyShip(s.getAuth(), ship).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetMyShipCargo(ship string) (stapi.ShipCargo, error) {
	s.timerTake()
	resp, r, err := s.apiClient.FleetApi.GetMyShipCargo(s.getAuth(), ship).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetMyShips() ([]stapi.Ship, error) {
	s.timerTake()
	// TODO: loop to get all of them
	resp, r, err := s.apiClient.FleetApi.GetMyShips(s.getAuth()).Limit(pageSize).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetShipCooldown(ship string) (stapi.Cooldown, error) {
	s.timerTake()
	resp, r, err := s.apiClient.FleetApi.GetShipCooldown(s.getAuth(), ship).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetShipNav(ship string) (stapi.ShipNav, error) {
	s.timerTake()
	resp, r, err := s.apiClient.FleetApi.GetShipNav(s.getAuth(), ship).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) InstallMount(ship string, mount string) (stapi.InstallMount201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewInstallMountRequest(mount)
	resp, r, err := s.apiClient.FleetApi.InstallMount(s.getAuth(), ship).InstallMountRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) Jettison(ship string, trade stapi.TradeSymbol, units int32) (stapi.Jettison200ResponseData, error) {
	s.timerTake()
	request := *stapi.NewJettisonRequest(trade, units)
	resp, r, err := s.apiClient.FleetApi.Jettison(s.getAuth(), ship).JettisonRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) JumpShip(ship string, system string) (stapi.JumpShip200ResponseData, error) {
	s.timerTake()
	request := *stapi.NewJumpShipRequest(system)
	resp, r, err := s.apiClient.FleetApi.JumpShip(s.getAuth(), ship).JumpShipRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) NavigateShip(ship string, waypoint string) (stapi.NavigateShip200ResponseData, error) {
	s.timerTake()
	request := *stapi.NewNavigateShipRequest(waypoint)
	resp, r, err := s.apiClient.FleetApi.NavigateShip(s.getAuth(), ship).NavigateShipRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) NegotiateContract(ship string) (stapi.NegotiateContract200ResponseData, error) {
	s.timerTake()
	resp, r, err := s.apiClient.FleetApi.NegotiateContract(s.getAuth(), ship).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) OrbitShip(ship string) (stapi.OrbitShip200ResponseData, error) {
	s.timerTake()
	resp, r, err := s.apiClient.FleetApi.OrbitShip(s.getAuth(), ship).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) PatchShipNav(ship string, mode string) (stapi.ShipNav, error) {
	s.timerTake()
	request := *stapi.NewPatchShipNavRequest()
	request.SetFlightMode(stapi.ShipNavFlightMode(mode))
	resp, r, err := s.apiClient.FleetApi.PatchShipNav(s.getAuth(), ship).PatchShipNavRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) PurchaseCargo(ship string, trade stapi.TradeSymbol, units int32) (stapi.SellCargo201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewPurchaseCargoRequest(trade, units)
	resp, r, err := s.apiClient.FleetApi.PurchaseCargo(s.getAuth(), ship).PurchaseCargoRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) PurchaseShip(ship stapi.ShipType, waypoint string) (stapi.PurchaseShip201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewPurchaseShipRequest(ship, waypoint)
	resp, r, err := s.apiClient.FleetApi.PurchaseShip(s.getAuth()).PurchaseShipRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) RefuelShip(ship string) (stapi.RefuelShip200ResponseData, error) {
	s.timerTake()
	request := *stapi.NewRefuelShipRequest()
	resp, r, err := s.apiClient.FleetApi.RefuelShip(s.getAuth(), ship).RefuelShipRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) RemoveMount(ship string, mount string) (stapi.RemoveMount201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewRemoveMountRequest(mount)
	resp, r, err := s.apiClient.FleetApi.RemoveMount(s.getAuth(), ship).RemoveMountRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) SellCargo(ship string, trade stapi.TradeSymbol, units int32) (stapi.SellCargo201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewSellCargoRequest(trade, units)
	resp, r, err := s.apiClient.FleetApi.SellCargo(s.getAuth(), ship).SellCargoRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) ShipRefine(ship string, produce string) (stapi.ShipRefine201ResponseData, error) {
	s.timerTake()
	request := *stapi.NewShipRefineRequest(produce)
	resp, r, err := s.apiClient.FleetApi.ShipRefine(s.getAuth(), ship).ShipRefineRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) TransferCargo(shipFrom string, trade stapi.TradeSymbol, units int32, shipTo string) (stapi.Jettison200ResponseData, error) {
	s.timerTake()
	request := *stapi.NewTransferCargoRequest(trade, units, shipTo)
	resp, r, err := s.apiClient.FleetApi.TransferCargo(s.getAuth(), shipFrom).TransferCargoRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) WarpShip(ship string, waypoint string) (stapi.NavigateShip200ResponseData, error) {
	s.timerTake()
	request := *stapi.NewNavigateShipRequest(waypoint)
	resp, r, err := s.apiClient.FleetApi.WarpShip(s.getAuth(), ship).NavigateShipRequest(request).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetJumpGate(system string, waypoint string) (stapi.JumpGate, error) {
	s.timerTake()
	resp, r, err := s.apiClient.SystemsApi.GetJumpGate(s.getAuth(), system, waypoint).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetMarket(system string, waypoint string) (stapi.Market, error) {
	s.timerTake()
	resp, r, err := s.apiClient.SystemsApi.GetMarket(s.getAuth(), system, waypoint).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetShipyard(system string, waypoint string) (stapi.Shipyard, error) {
	s.timerTake()
	resp, r, err := s.apiClient.SystemsApi.GetShipyard(s.getAuth(), system, waypoint).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetSystem(system string) (stapi.System, error) {
	s.timerTake()
	resp, r, err := s.apiClient.SystemsApi.GetSystem(s.getAuth(), system).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetSystemWaypoints(system string) ([]stapi.Waypoint, error) {
	s.timerTake()
	// TODO: loop to get all of them
	resp, r, err := s.apiClient.SystemsApi.GetSystemWaypoints(s.getAuth(), system).Limit(pageSize).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetSystems() ([]stapi.System, error) {
	s.timerTake()
	resp, r, err := s.apiClient.SystemsApi.GetSystems(s.getAuth()).Limit(pageSize).Execute()
	return resp.GetData(), handleErr(r, err)
}

func (s *Client) GetMoreSystems(start int, count int) ([]stapi.System, error) {
	pageStart := (int32(start) / pageSize) + 1
	pageEnd := int32(math.Ceil(float64(start+count) / float64(pageSize)))
	var result []stapi.System
	for i := pageStart; i <= pageEnd; i++ {
		s.timerTake()
		resp, r, err := s.apiClient.SystemsApi.GetSystems(s.getAuth()).Limit(pageSize).Page(i).Execute()
		result = append(result, resp.GetData()...)
		if (err != nil) || (len(resp.GetData()) < int(pageSize)) {
			return result, handleErr(r, err)
		}
	}

	return result, nil
}

func (s *Client) GetWaypoint(system string, waypoint string) (stapi.Waypoint, error) {
	s.timerTake()
	resp, r, err := s.apiClient.SystemsApi.GetWaypoint(s.getAuth(), system, waypoint).Execute()
	return resp.GetData(), handleErr(r, err)
}

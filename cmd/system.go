/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/bgreen/space-traders-go/stapi"
	"github.com/spf13/cobra"
)

// systemCmd represents the system command
var systemCmd = &cobra.Command{
	Use:   "system [symbol]",
	Short: "Print system info",
	Long:  `Print system info`,
	Run: func(cmd *cobra.Command, args []string) {
		var systems []stapi.System
		c, _ := cmd.Flags().GetInt("count")

		if len(args) == 0 {
			resp, err := client.GetMoreSystems(0, c)
			if err != nil {
				return
			}
			systems = append(systems, resp...)

		} else if len(args) == 1 {
			resp, err := client.GetSystem(args[0])
			if err != nil {
				return
			}
			systems = append(systems, resp)
		}

		for i, s := range systems {
			if i == c {
				return
			}
			l, _ := cmd.Flags().GetBool("long")
			if l {
				wps, _ := client.GetSystemWaypoints(s.Symbol)
				fmt.Print(systemInfoLong(s, wps) + "\n")
			} else {
				fmt.Print(systemInfoShort(s))
			}
		}
	},
}

func systemInfoShort(sys stapi.System) string {
	return fmt.Sprintf("%-7v:\t%v\tX:%6v Y:%6v\n", sys.Symbol, sys.Type, sys.X, sys.Y)
}

func systemInfoLong(sys stapi.System, wps []stapi.Waypoint) string {
	var s string
	s += fmt.Sprintf("%-7v:\t%v\tX:%6v Y:%6v\n", sys.Symbol, sys.Type, sys.X, sys.Y)
	s += fmt.Sprintln("Waypoints:")
	for _, v := range wps {
		s += fmt.Sprintf("%-14v:\t%-14v\tX:%3v Y:%3v\t", v.Symbol, v.Type, v.X, v.Y)
		for _, t := range v.Traits {
			s += fmt.Sprintf("%v, ", t.Name)
		}
		s += "\n"
	}

	return s
}

var systemMarketCmd = &cobra.Command{
	Use:   "market <symbol>",
	Short: "Print market info",
	Long:  `Print market info`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var markets []stapi.Market

		if isSystemSymbol(args[0]) {
			resp, err := client.GetSystemWaypoints(args[0])
			if err != nil {
				return
			}
			for _, v := range resp {
				if isWaypointMarket(v) {
					m, _ := client.GetMarket(v.SystemSymbol, v.Symbol)
					markets = append(markets, m)
				}
			}

		} else if isWaypointSymbol(args[0]) {
			resp, err := client.GetMarket(waypointSymbolToSystemSymbol(args[0]), args[0])
			if err != nil {
				return
			}
			markets = append(markets, resp)
		}

		for _, m := range markets {
			line := "==========================================================="
			fmt.Println(line)
			fmt.Print(marketInfoLong(m))
			fmt.Println(line)
		}
	},
}

var systemShipyardCmd = &cobra.Command{
	Use:   "shipyard <symbol>",
	Short: "Print shipyard info",
	Long:  `Print shipyard info`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var markets []stapi.Shipyard

		if isSystemSymbol(args[0]) {
			resp, err := client.GetSystemWaypoints(args[0])
			if err != nil {
				return
			}
			for _, v := range resp {
				if isWaypointShipyard(v) {
					m, _ := client.GetShipyard(v.SystemSymbol, v.Symbol)
					markets = append(markets, m)
				}
			}

		} else if isWaypointSymbol(args[0]) {
			resp, err := client.GetShipyard(waypointSymbolToSystemSymbol(args[0]), args[0])
			if err != nil {
				return
			}
			markets = append(markets, resp)
		}

		for _, m := range markets {
			line := "==========================================================="
			fmt.Println(line)
			fmt.Print(shipyardInfoLong(m))
			fmt.Println(line)
		}
	},
}

var systemGateCmd = &cobra.Command{
	Use:   "gate <symbol>",
	Short: "Print Jump Gate info",
	Long:  `Print Jump Gate info`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var gates []stapi.JumpGate

		if isSystemSymbol(args[0]) {
			resp, err := client.GetSystemWaypoints(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, v := range resp {
				if isWaypointJumpGate(v) {
					g, _ := client.GetJumpGate(v.SystemSymbol, v.Symbol)
					gates = append(gates, g)
				}
			}

		} else if isWaypointSymbol(args[0]) {
			resp, err := client.GetJumpGate(waypointSymbolToSystemSymbol(args[0]), args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			gates = append(gates, resp)
		}

		for _, m := range gates {
			line := "==========================================================="
			fmt.Println(line)
			fmt.Print(gateInfoLong(m))
			fmt.Println(line)
		}
	},
}

func isSectorSymbol(s string) bool {
	return strings.Count(s, "-") == 0
}

func isSystemSymbol(s string) bool {
	return strings.Count(s, "-") == 1
}

func isWaypointSymbol(s string) bool {
	return strings.Count(s, "-") == 2
}

func isWaypointMarket(w stapi.Waypoint) bool {
	for _, v := range w.Traits {
		if v.Symbol == "MARKETPLACE" {
			return true
		}
	}
	return false
}

func isWaypointShipyard(w stapi.Waypoint) bool {
	for _, v := range w.Traits {
		if v.Symbol == "SHIPYARD" {
			return true
		}
	}
	return false
}

func isWaypointJumpGate(w stapi.Waypoint) bool {
	return w.Type == stapi.WaypointType("JUMP_GATE")
}

func waypointSymbolToSystemSymbol(w string) string {
	ss := strings.Split(w, "-")
	return strings.Join(ss[:2], "-")
}

func marketInfoLong(m stapi.Market) string {
	var s string
	s += m.Symbol + "\n"
	s += "Exchange: "
	for _, v := range m.Exchange {
		s += fmt.Sprintf("%v ", v.Symbol)
	}
	s += "\n"
	s += "Exports: "
	for _, v := range m.Exports {
		s += fmt.Sprintf("%v ", v.Symbol)
	}
	s += "\n"
	s += "Imports: "
	for _, v := range m.Imports {
		s += fmt.Sprintf("%v ", v.Symbol)
	}
	s += "\n"
	s += "Trade Goods\n"
	for _, v := range m.TradeGoods {
		s += fmt.Sprintf("%-24v:\tBuy:%5v\tSell:%5v\tVol:%5v\tSupply:%v\n", v.Symbol, v.PurchasePrice, v.SellPrice, v.TradeVolume, v.Supply)
	}
	/*
		s += "Recent Transactions\n"
		for _, v := range m.Transactions {
			s += fmt.Sprintf("%v: %-24v\tPrice:%5v\tUnits:%3v\n", v.Timestamp.Format(time.Stamp), v.TradeSymbol, v.PricePerUnit, v.Units)
		}
	*/
	return s
}

func shipyardInfoLong(y stapi.Shipyard) string {
	var s string
	s += y.Symbol + "\n"
	s += "Ship Types: "
	for _, v := range y.ShipTypes {
		s += fmt.Sprintf("%v ", *v.Type)
	}
	s += "\n"
	s += "Ships\n"
	for _, v := range y.Ships {
		s += fmt.Sprintf("%-20v: %-20v\t%v\n", *v.Type, v.Frame.Name[6:], v.PurchasePrice)
	}
	s += "\n"
	/*
		s += "Recent Transactions\n"
		for _, v := range y.Transactions {
			s += fmt.Sprintf("%v: Agent:%-20v\tShip:%v\tPrice:%5v\n", v.Timestamp.Format(time.Stamp), v.AgentSymbol, v.ShipSymbol, v.Price)
		}
	*/
	return s
}

func gateInfoLong(g stapi.JumpGate) string {
	var s string
	s += fmt.Sprintf("Faction:%v\tRange:%v\n", *g.FactionSymbol, g.JumpRange)
	s += "Connected Systems:\n"
	for _, v := range g.ConnectedSystems {
		s += fmt.Sprintf("%-7v: %v\tX:%v Y:%v\tDistance:%v\t\n", v.Symbol, v.Type, v.X, v.Y, v.Distance)
	}
	return s
}

var systemShipyardBuyCmd = &cobra.Command{
	Use:   "buy <waypoint> <ship type>",
	Short: "Buy a ship",
	Long:  `Buy a ship`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.PurchaseShip(stapi.ShipType(args[1]), args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", resp)
	},
}

func init() {
	rootCmd.AddCommand(systemCmd)

	systemCmd.AddCommand(systemMarketCmd)
	systemCmd.AddCommand(systemShipyardCmd)
	systemCmd.AddCommand(systemGateCmd)

	systemShipyardCmd.AddCommand(systemShipyardBuyCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// systemCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// systemCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	systemCmd.Flags().IntP("count", "c", 10, "Maximum number of systems to list")
}

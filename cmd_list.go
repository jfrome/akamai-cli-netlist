package main

import (
	common "github.com/apiheat/akamai-cli-common"
	"github.com/urfave/cli"
)

// cmdlistNetLists is used by cli to execute action of listing all network lists
//
// cmd_list
func cmdlistNetLists(c *cli.Context) error {
	return listNetLists(c)
}

// cmdlistNetListId is used by cli to execute action of listing specific network list
//
// cmd_list
func cmdlistNetListID(c *cli.Context) error {
	return listNetListbyID(c)
}

// cmdlistNetListName is used by cli to execute action of listing specific network list
//
// cmd_list
func cmdlistNetListName(c *cli.Context) error {
	return listNetListbyName(c)
}

// listNetLists execute client API call to get all network lists
//
// cmd_list
func listNetLists(c *cli.Context) error {

	// Since we are listing all we do not filter results
	listNetListOptsv2.Search = ""

	// List all IP based
	listNetListOptsv2.TypeOflist = "IP"
	netListsIP, _, netlistErr := apiClient.NetworkListsv2.ListNetworkLists(listNetListOptsv2)
	if netlistErr != nil {
		return netlistErr
	}

	// List all GEO based
	listNetListOptsv2.TypeOflist = "GEO"
	netListsGEO, _, netlistErr := apiClient.NetworkListsv2.ListNetworkLists(listNetListOptsv2)
	if netlistErr != nil {
		return netlistErr
	}

	netListsResult := append(*netListsIP, *netListsGEO...)
	common.OutputJSON(netListsResult)

	return nil
}

// listNetLists execute client API call to get specific network list
//
// cmd_list
func listNetListbyID(c *cli.Context) error {
	common.VerifyArgumentByName(c, "id")

	netList, _, netlistErr := apiClient.NetworkListsv2.GetNetworkList(listID, listNetListOptsv2)
	if netlistErr != nil {
		return netlistErr
	}

	common.OutputJSON(netList)

	return nil
}

// listNetLists execute client API call to get specific network list
//
// cmd_list
func listNetListbyName(c *cli.Context) error {
	common.VerifyArgumentByName(c, "name")

	listNetListOptsv2.Search = c.String("name")
	netList, _, netlistErr := apiClient.NetworkListsv2.ListNetworkLists(listNetListOptsv2)
	if netlistErr != nil {
		return netlistErr
	}

	common.OutputJSON(netList)

	return nil
}

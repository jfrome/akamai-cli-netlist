package main

import (
	"fmt"
	"os"
	"sort"

	common "github.com/apiheat/akamai-cli-common"
	edgegrid "github.com/apiheat/go-edgegrid"
	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

var (
	apiClient       *edgegrid.Client
	appVer, appName string
)

func main() {
	app := common.CreateNewApp(appName, "A CLI to interact with Akamai network lists", appVer)
	app.Flags = common.CreateFlags()
	app.Before = func(c *cli.Context) error {
		var err error

		// Provide struct details needed for apiClient init
		apiClientOpts := &edgegrid.ClientOptions{}
		apiClientOpts.ConfigPath = c.GlobalString("config")
		apiClientOpts.ConfigSection = c.GlobalString("section")
		apiClientOpts.DebugLevel = c.GlobalString("debug")
		apiClientOpts.AccountSwitchKey = c.GlobalString("ask")

		apiClient, err = common.EdgeClientInit(apiClientOpts)

		if err != nil {
			log.Fatalln(err)
			return cli.NewExitError(err, 1)
		}

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "get",
			Usage: "List network lists objects",
			Subcommands: []cli.Command{
				{
					Name:      "all",
					Usage:     "Gets all network list in the account",
					UsageText: fmt.Sprintf("%s get all [command options]", appName),
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "extended",
							Usage: "returns more verbose data such as creation date and activation status",
						},
						cli.BoolFlag{
							Name:  "includeElements",
							Usage: "includes the full list of IP or GEO elements",
						},
						cli.StringFlag{
							Name:  "listType",
							Value: "ANY",
							Usage: "filters by the network list type [ IP | GEO | ANY ]",
						},
					},
					Action: cmdlistNetLists,
				},
				{
					Name:      "by-id",
					Usage:     "Gets a network list by unique-id",
					UsageText: fmt.Sprintf("%s get by-id --id UNIQUE-ID [command options]", appName),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "id",
							Usage: "list unique-id",
						},
						cli.BoolFlag{
							Name:  "extended",
							Usage: "returns more verbose data such as creation date and activation status",
						},
						cli.BoolFlag{
							Name:  "includeElements",
							Usage: "includes the full list of IP or GEO elements",
						},
					},
					Action: cmdlistNetListID,
				},
				{
					Name:      "by-name",
					Usage:     "Gets a network list by name",
					UsageText: fmt.Sprintf("%s get by-id --name NAME [command options]", appName),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name",
							Usage: "list name",
						},
						cli.BoolFlag{
							Name:  "extended",
							Usage: "returns more verbose data such as creation date and activation status",
						},
						cli.BoolFlag{
							Name:  "includeElements",
							Usage: "includes the full list of IP or GEO elements",
						},
						cli.StringFlag{
							Name:  "listType",
							Value: "IP",
							Usage: "filters by the network list type [ IP | GEO ]",
						},
					},
					Action: cmdlistNetListName,
				},
			},
		},
		{
			Name:      "search",
			Usage:     "Finds all network lists that match specific expression ( either name or network element )",
			UsageText: fmt.Sprintf("%s search --searchPattern SEARCH-ELEMENT [command options]", appName),
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "extended",
					Usage: "returns more verbose data such as creation date and activation status",
				},
				cli.StringFlag{
					Name:  "searchPattern",
					Usage: "includes network lists that match search pattern",
				},
				cli.StringFlag{
					Name:  "listType",
					Value: "ANY",
					Usage: "filters by the network list type [ IP | GEO | ANY ]",
				},
			},
			Action: cmdSearchNetLists,
		},
		{
			Name:      "sync",
			Usage:     "Synchronizes items from source list into destination list ( without activation )",
			UsageText: fmt.Sprintf("%s sync-items --id-src SOURCE-LIST-ID --id-dst TARGET-LIST-ID [command options]", appName),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id-src",
					Usage: "Source list ID to take items from",
				},
				cli.StringFlag{
					Name:  "id-dst",
					Usage: "Target list ID to which items should be added",
				},
			},
			Action: cmdSyncNetListID,
		},
		{
			Name:  "items",
			Usage: "Manages items in network lists",
			Subcommands: []cli.Command{
				{
					Name:      "add",
					Usage:     "Adds network list element to provided network list",
					UsageText: fmt.Sprintf("%s items add --id UNIQUE-ID --items ITEM1,ITEM2,ITEM3", appName),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "id",
							Usage: "list unique-id",
						},
						cli.StringSliceFlag{
							Name:  "items",
							Usage: "items to be included",
						},
					},
					Action: cmdAddItemsToNetlist,
				},
				{
					Name:      "remove",
					Usage:     "Removes network list element from provided network list",
					UsageText: fmt.Sprintf("%s items remove --id UNIQUE-ID --element ELEMENT", appName),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "id",
							Usage: "list unique-id",
						},
						cli.StringFlag{
							Name:  "element",
							Usage: "element to be removed",
						},
					},
					Action: cmdRemoveItemFromNetlist,
				},
			},
		},
		{
			Name:      "create",
			Usage:     "Creates new network list",
			UsageText: fmt.Sprintf("%s create --name NETWORK-LIST-NAME [command options]", appName),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Value: "",
					Usage: "name for the new list",
				},
				cli.StringFlag{
					Name:  "description",
					Value: "created via akamai-cli-networklist",
					Usage: "description for the new list",
				},
				cli.StringFlag{
					Name:  "type",
					Value: "IP",
					Usage: "defines type of list for creation (IP/GEO)",
				},
			},
			Action: cmdCreateNetList,
		},
		{
			Name:  "activate",
			Usage: "Manages network list activation/status",
			Subcommands: []cli.Command{
				{
					Name:      "list",
					Usage:     "Activates network list on given network",
					UsageText: fmt.Sprintf("%s activate list --id UNIQUE-ID [command options]", appName),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "id",
							Usage: "list unique-id",
						},
						cli.StringFlag{
							Name:  "comments",
							Value: "activated via akamai-cli",
							Usage: "comments",
						},
						cli.StringSliceFlag{
							Name:  "notificationRecipients",
							Usage: "recipients of notification",
						},
						cli.BoolFlag{
							Name:  "fast",
							Usage: "n/a",
						},
						cli.BoolFlag{
							Name:  "prd",
							Usage: "activate on production",
						},
					},
					Action: cmdActivateNetList,
				},
				{
					Name:      "status",
					Usage:     "Displays activation status for given network list",
					UsageText: fmt.Sprintf("%s activate status --id UNIQUE-ID [command options]", appName),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "id",
							Usage: "list unique-id",
						},
						cli.BoolFlag{
							Name:  "prd",
							Usage: "activate on production",
						},
					},
					Action: cmdActivateNetListStatus,
				},
			},
		},
		{
			Name:      "delete",
			Usage:     "Deletes network list ( ** REQUIRES LIST TO BE DEACTIVATED ON BOTH NETWORKS ** )",
			UsageText: fmt.Sprintf("%s delete --id UNIQUE-ID", appName),
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "list unique-id to remove",
				},
			},
			Action: cmdRemoveNetlist,
		},
		{
			Name:      "notification",
			Usage:     "Manages network list subscription notifications ( SUBSCRIBE by default ) ",
			UsageText: fmt.Sprintf("%s notification status --id UNIQUE-ID --notificationRecipients RECIPIENTS [command options]", appName),
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "networkListsIDs",
					Usage: "recipients of notification",
				},
				cli.StringSliceFlag{
					Name:  "notificationRecipients",
					Usage: "recipients of notification",
				},
				cli.BoolFlag{
					Name:  "unsubscribe",
					Usage: "Unsubscribe from notifications",
				},
			},
			Action: cmdNotificationManagement,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Action = func(c *cli.Context) error {

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

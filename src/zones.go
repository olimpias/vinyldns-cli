/*
Copyright 2018 Comcast Cable Communications Management, LLC
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"

	clitable "github.com/crackcomm/go-clitable"
	"github.com/urfave/cli"
	"github.com/vinyldns/go-vinyldns/vinyldns"
)

func zones(c *cli.Context) error {
	client := client(c)
	zones, err := client.ZonesListAll(vinyldns.ListFilter{})
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(zones)
	}

	data := [][]string{}
	for _, z := range zones {
		data = append(data, []string{
			z.Name,
			z.ID,
		})
	}

	if len(data) != 0 {
		printTableWithHeaders([]string{"Name", "ID"}, data)
	} else {
		fmt.Printf("No zones found")
	}

	return nil
}

func zone(c *cli.Context) error {
	client := client(c)
	name := c.String("zone-name")
	id := c.String("zone-id")
	z, err := getZone(client, name, id)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(z)
	}

	data := [][]string{
		{"Name", z.Name},
		{"ID", z.ID},
		{"Status", z.Status},
	}

	printBasicTable(data)

	return nil
}

func zoneDelete(c *cli.Context) error {
	id := c.String("zone-id")
	client := client(c)
	deleted, err := client.ZoneDelete(id)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(deleted)
	}

	fmt.Printf("Deleted zone %s\n", id)

	return nil
}

func zoneCreate(c *cli.Context) error {
	client := client(c)
	id, err := getAdminGroupID(client, c.String("admin-group-id"), c.String("admin-group-name"))
	if err != nil {
		return err
	}
	connection := &vinyldns.ZoneConnection{
		Key:           c.String("zone-connection-key"),
		KeyName:       c.String("zone-connection-key-name"),
		Name:          c.String("zone-connection-key-name"),
		PrimaryServer: c.String("zone-connection-primary-server"),
	}
	tConnection := &vinyldns.ZoneConnection{
		Key:           c.String("transfer-connection-key"),
		KeyName:       c.String("transfer-connection-key-name"),
		Name:          c.String("transfer-connection-key-name"),
		PrimaryServer: c.String("transfer-connection-primary-server"),
	}
	z := &vinyldns.Zone{
		Name:         c.String("name"),
		Email:        c.String("email"),
		AdminGroupID: id,
	}

	zc, err := validateConnection("zone", connection)
	if err != nil {
		return err
	}
	if zc {
		z.Connection = connection
	}

	tc, err := validateConnection("transfer", tConnection)
	if err != nil {
		return err
	}
	if tc {
		z.TransferConnection = tConnection
	}

	created, err := client.ZoneCreate(z)
	if err != nil {
		return err
	}

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(created)
	}

	fmt.Printf("Created zone %s\n", created.Zone.Name)

	return nil
}

func zoneConnection(c *cli.Context) error {
	client := client(c)
	id := c.String("zone-id")
	z, err := client.Zone(id)
	if err != nil {
		return err
	}
	con := z.Connection

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(con)
	}

	if con == nil {
		fmt.Printf("No zone connection found for zone %s\n", id)

		return nil
	}

	data := [][]string{
		{"Name", con.Name},
		{"KeyName", con.KeyName},
		{"Key", con.Key},
		{"PrimaryServer", con.PrimaryServer},
	}

	printBasicTable(data)

	return nil
}

func zoneChanges(c *cli.Context) error {
	client := client(c)
	zh, err := client.ZoneHistory(c.String("zone-id"))
	if err != nil {
		return err
	}
	cs := zh.ZoneChanges

	if c.GlobalString(outputFlag) == "json" {
		return printJSON(cs)
	}

	for _, c := range cs {
		clitable.PrintHorizontal(map[string]interface{}{
			"Zone":       c.Zone.Name,
			"ZoneID":     c.Zone.ID,
			"UserID":     c.UserID,
			"ChangeType": c.ChangeType,
			"Status":     c.Status,
			"Created":    c.Created,
			"ID":         c.ID,
		})
	}

	return nil
}

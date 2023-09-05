package main

import (
	"context"
	"fmt"

	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	ip       = "labvm1.myad.lab"
	user     = "user"
	password = "password"
)

func main() {
	u := &url.URL{
		Scheme: "https",
		Host:   ip,
		Path:   "/sdk",
	}
	ctx := context.Background()
	u.User = url.UserPassword(user, password)
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		panic(err)
	}

	m := view.NewManager(client.Client)

	v, err := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		panic(err)
	}
	defer v.Destroy(ctx)

	var hosts []mo.HostSystem
	err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"config"}, &hosts)
	if err != nil {
		panic(err)
	}

	for _, host := range hosts {
		for _, service := range host.Config.Service.Service {
			if service.Key == "TSM-SSH" {
				if service.Running {
					fmt.Printf("%s is running.\n", service.Key)
				} else {
					fmt.Printf("%s is not running.\n", service.Key)
				}

			}

		}
	}
	ref := types.ManagedObjectReference{
		Type:  hosts[0].Reference().Type,
		Value: hosts[0].Reference().Value,
	}
	var hostcm = object.NewHostConfigManager(client.Client, ref)

	s, err := hostcm.ServiceSystem(ctx)
	if err != nil {
		panic(err)
	}
	s.Stop(ctx, "TSM-SSH")
}

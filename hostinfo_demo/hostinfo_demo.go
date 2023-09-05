package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"text/tabwriter"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const (
		ip       = "192.168.10.12"
		user     = "root"
		password = "P@ssw0rd"
	)

	vURL := &url.URL{
		Scheme: "https",
		Host:   ip,
		Path:   "/sdk",
	}

	vURL.User = url.UserPassword(user, password)
	client, err := govmomi.NewClient(ctx, vURL, true)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Logout(ctx)

	m := view.NewManager(client.Client)
	v, err := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		fmt.Println(err)
	}
	defer v.Destroy(ctx)

	var hss []mo.HostSystem
	err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"summary"}, &hss)
	if err != nil {
		fmt.Println(err)
	}

	tw := tabwriter.NewWriter(os.Stdout, 2, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Name:\tUsed CPU:\tTotal CPU:\tFree CPU:\tUsed Memory:\tTotal Memory:\tFree Memory:\t\n")

	for _, hs := range hss {
		totalCPU := int64(hs.Summary.Hardware.CpuMhz) * int64(hs.Summary.Hardware.NumCpuCores)
		freeCPU := int64(totalCPU) - int64(hs.Summary.QuickStats.OverallCpuUsage)
		freeMemory := int64(hs.Summary.Hardware.MemorySize) - (int64(hs.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024)
		fmt.Fprintf(tw, "%s\t", hs.Summary.Config.Name)
		fmt.Fprintf(tw, "%d\t", hs.Summary.QuickStats.OverallCpuUsage)
		fmt.Fprintf(tw, "%d\t", totalCPU)
		fmt.Fprintf(tw, "%d\t", freeCPU)
		fmt.Fprintf(tw, "%s\t", (units.ByteSize(hs.Summary.QuickStats.OverallMemoryUsage))*1024*1024)
		fmt.Fprintf(tw, "%s\t", units.ByteSize(hs.Summary.Hardware.MemorySize))
		fmt.Fprintf(tw, "%d\t", freeMemory)
		fmt.Fprintf(tw, "\n")
	}
	_ = tw.Flush()
}

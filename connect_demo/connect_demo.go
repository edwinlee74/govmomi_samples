package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/vmware/govmomi"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const (
		ip       = "vcenter6.myad.lab"
		user     = "administrator@vsphere.local"
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
	fmt.Println(client.ServiceContent.About.Name)
}

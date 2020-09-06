package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-provider-squadcast/squadcast"
)

func main() {

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: squadcast.Provider})

}

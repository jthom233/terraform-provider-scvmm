package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return &schema.Provider{
				Schema: map[string]*schema.Schema{
					"vmm_server": {
						Type:        schema.TypeString,
						Required:    true,
						DefaultFunc: schema.EnvDefaultFunc("VMM_SERVER", nil),
						Description: "SCVMM server address",
					},
					"username": {
						Type:        schema.TypeString,
						Required:    true,
						DefaultFunc: schema.EnvDefaultFunc("VMM_USERNAME", nil),
						Description: "SCVMM username",
					},
					"password": {
						Type:        schema.TypeString,
						Required:    true,
						Sensitive:   true,
						DefaultFunc: schema.EnvDefaultFunc("VMM_PASSWORD", nil),
						Description: "SCVMM password",
					},
				},
				ResourcesMap: map[string]*schema.Resource{
					"scvmm_vm": resourceVM(),
				},
				ConfigureFunc: providerConfigure,
			}
		},
	})
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	vmmServer := d.Get("vmm_server").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Here you can set up a client to interact with SCVMM
	return &SCVMMClient{
		VMMServer: vmmServer,
		Username:  username,
		Password:  password,
	}, nil
}

type SCVMMClient struct {
	VMMServer string
	Username  string
	Password  string
}

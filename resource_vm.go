package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceVMCreate,
		Read:   resourceVMRead,
		Update: resourceVMUpdate,
		Delete: resourceVMDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template": {
				Type:     schema.TypeString,
				Required: true,
			},
			"computer_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cloud_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func runPowerShell(script string) (string, error) {
	cmd := exec.Command("powershell", "-Command", script)
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func resourceVMCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*SCVMMClient)

	name := d.Get("name").(string)
	template := d.Get("template").(string)
	computerName := d.Get("computer_name").(string)
	description := d.Get("description").(string)
	cloudName := d.Get("cloud_name").(string)

	script := fmt.Sprintf(`
        param (
            [string]$Template,
            [string]$ComputerName,
            [string]$VMName,
            [string]$Description,
            [string]$CloudName,
            [string]$VMMServer,
            [PSCredential]$Credential
        )

        Import-Module virtualmachinemanager

        Get-SCVMMServer -ComputerName $VMMServer -Credential $Credential
        $Cloud = Get-SCCloud -Name $CloudName
        $vmTemplate = Get-SCVMTemplate -Name $Template
        $vmConfig = New-SCVMConfiguration -VMTemplate $vmTemplate -Cloud $Cloud -Name $VMName -Description $Description
        $vmConfig | Set-SCVMConfiguration -Name $VMName -ComputerName $ComputerName
        Update-SCVMConfiguration $vmConfig
        New-SCVirtualMachine -Name $VMName -VMConfiguration $vmConfig -Cloud $Cloud -Description $Description -RunAsynchronously
        Start-Sleep -Seconds 60
    `, template, computerName, name, description, cloudName, client.VMMServer, fmt.Sprintf("%s:%s", client.Username, client.Password))

	_, err := runPowerShell(script)
	if err != nil {
		return err
	}

	d.SetId(name)
	return resourceVMRead(d, m)
}

func resourceVMRead(d *schema.ResourceData, m interface{}) error {
	// Implement the read logic
	return nil
}

func resourceVMUpdate(d *schema.ResourceData, m interface{}) error {
	// Implement the update logic
	return nil

}

func resourceVMDelete(d *schema.ResourceData, m interface{}) error {
	// Implement the delete logic
	return nil
}

package main

import (
	"fmt"

	"github.com/Cloudbase/go-powershell"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func resourceVMCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*SCVMMClient)

	name := d.Get("name").(string)
	template := d.Get("template").(string)
	computerName := d.Get("computer_name").(string)
	description := d.Get("description").(string)
	cloudName := d.Get("cloud_name").(string)

	// Use go-powershell to execute the PowerShell script
	ps, err := powershell.New(&powershell.DefaultBackend{})
	if err != nil {
		return err
	}
	defer ps.Exit()

	script := `
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
    `

	// Set parameters for the script
	params := map[string]interface{}{
		"Template":     template,
		"ComputerName": computerName,
		"VMName":       name,
		"Description":  description,
		"CloudName":    cloudName,
		"VMMServer":    client.VMMServer,
		"Credential":   fmt.Sprintf("%s:%s", client.Username, client.Password),
	}

	_, err = ps.Execute(script, params)
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

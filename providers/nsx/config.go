package main

import (
        "git.devops.int.ovp.bskyb.com/paas/gonsx/client"
	"log"
)

// Config is a struct for containing the provider parameters.
type Config struct {
	Debug         bool
	Insecure      bool
	NSXUserName   string
	NSXPassword   string
	NSXServer     string
}

// Client returns a new client for accessing VMWare vSphere.
func (c *Config) Client() (*client.NSXClient, error) {
	log.Printf("[INFO] VMWare NSX Client configured for URL: %s", c.NSXServer)
        nsxclient := client.NewNSXClient("https://"+c.NSXServer, c.NSXUserName, c.NSXPassword, c.Insecure, c.Debug)
	return nsxclient, nil
}

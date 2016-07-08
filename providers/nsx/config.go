package main

import (
	"fmt"
)

type Config struct {
	Debug         bool
	Insecure      bool
	NSXUserName   string
	NSXPassword   string
	NSXServer     string
}

// Client() returns a new client for accessing VMWare vSphere.
func (c *Config) Client() (*Config, error) {
	fmt.Printf("[INFO] VMWare NSX Client configured for URL: %s", c.NSXServer)
	return c, nil
}

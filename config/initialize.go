package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func Initialize(cmd *cobra.Command, args []string) {
	fmt.Println("Initializing configuration...\n")
	initialConfig := []byte(`stimmtausch:
  # This is a basic Stimmtausch configuration file! It contains some handy
  # default settings, as well as some stubbed-out data you'll want to change
  # to suit your own needs. Everything that comes after a '#' is a comment
  # and will be ignored (like these lines!), so you can feel free to delete
  # these or add your own for future reference.
  #
  # All options in the configuration file come as key/value pairs, separated
  # by a ':'. However, you'll notice that these can be nested within each
  # other, rather than just a simple string of characters. This format is
  # called YAML and you can read more about it here: https://yaml.org
  #
  # If YAML isn't your favorite, you can use anything that works with the
  # config file reader we use, Viper, which are listed here:
  # https://github.com/spf13/viper
  
  # The configuration file version. As more features are added to 'st' down
  # the line, the format of this file may change, so it's important to specify
  # what version the program should expect!
  version: 1

  # The main difference in various server types is the connection string.
  # An example server type for TinyMUCK style servers is show below, but
  # feel free to add your own.
  server_types:
    MUCK:
      connect_string: "connect $user $password"

  # You can specify the default server type to use for the serveres below.
  default_server_type: MUCK
  servers:
    spr:
      host: muck.sprmuck.org
      port: 23551
      ssl: true
    furrymuck:
      host: furrymuck.com
      port: 8899
      ssl: true
    tapestries:
      host: tapestries.fur.com
      port: 6699
      ssl: true
    spindizzy:
      host: muck.spindizzy.org
      port: 7073
      ssl: true
      # If you want to use a different server type for one of the
      # servers, specify it like this.
      type: MUCK

  # Worlds are a cobination of a server and a character (think "user", in
  # modern web parlance). The worlds below are just examples. You'll probably
  # want to change them :)
  worlds:
    taps_foxface:
      server: tapestries
      username: foxface
      password: 12345
    taps_rudderbutt:
      # You can associate as many characters as you want with a server.
      server: tapestries
      username: rudderbutt
      password: 67890
    furry_foxface:
      server: furrymuck
      username: foxface
      password: 12345

  # The default world is what is connected to when you run Stimmtausch
  # without specifying a world on the command line.
  default_world: taps_foxface

  # You can specify the ways the client acts below.
  client:
    # ...whether or not to show the system log in a tiny pane.
    show_syslog: true
    # ...what log level to show (TRACE, DEBUG, INFO, WARNING, ERROR, CRITICAL)
    log_level: INFO
`)
	home, err := HomeDir()
	if err != nil {
		log.Criticalf("unable to get homedir, bailing... %v", err)
		os.Exit(4)
	}
	if _, err = os.Stat(home); err != nil {
		fmt.Printf("Looks like %s doesn't exist yet. Creating that for you...\n\n", home)
		if err = os.MkdirAll(home, 0755); err != nil {
			log.Criticalf("unable to ensure config directory! %v", err)
			os.Exit(4)
		}
	}
	configFile := filepath.Join(home, "config.yaml")
	if _, err = os.Stat(configFile); err == nil {
		fmt.Println("Yikes! you already have a config.yaml file! I won't overwite that!")
		fmt.Println("Bailing early. If you want to run this again, maybe stash that file somewhere safe first.")
		os.Exit(4)
	}
	err = ioutil.WriteFile(configFile, initialConfig, 0644)
	if err != nil {
		log.Criticalf("uanble to write config file! %v", err)
		os.Exit(4)
	}
	fmt.Printf(`A new configuration file has been written for you in %s!

This file contains a bunch of sensible defaults, but also some stubbed-out data,
so please take a moment to edit it and make sure it contains all of your info!

Welcome to Stimmtausch :)
`, configFile)
}

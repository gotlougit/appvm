package main

import (
	"io/ioutil"
)

// Builtin VMs

type app struct {
	Name string
	Nix  []byte
}

var builtin_brave_nix = app{
	Name: "brave",
	Nix: []byte(`
{pkgs, ...}:
let
  application = "${pkgs.brave}/bin/brave";
  appRunner = pkgs.writeShellScriptBin "app" ''
    ARGS_FILE=/home/user/.args
    ARGS=$(cat $ARGS_FILE)
    rm $ARGS_FILE

    ${application} $ARGS
    systemctl poweroff
  '';
in {
  imports = [
    ./base.nix
  ];

  services.xserver.displayManager.sessionCommands = "${appRunner}/bin/app &";
}
`),
}

func writeBuiltinApps(path string) (err error) {
	for _, f := range []app{
		builtin_brave_nix,
	} {
		err = ioutil.WriteFile(configDir+"/"+f.Name+".nix", f.Nix, 0644)
		if err != nil {
			return
		}
	}

	return
}

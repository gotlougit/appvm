package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var template = `
{pkgs, ...}:
let
  application = "${pkgs.%s}/bin/%s";
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
`

func isPackageExists(name string) bool {
	cmd := exec.Command("nix", "eval", "--impure", "--expr", fmt.Sprintf("(import <nixpkgs> {}).%s or null", name))
	return cmd.Run() == nil
}

func nixPath(name string) (path string, err error) {
	command := exec.Command("nix", "eval", "--impure", "--raw", "--expr", fmt.Sprintf("(import <nixpkgs> {}).%s", name))
	bytes, err := command.Output()
	if err != nil {
		return
	}
	path = string(bytes)
	return
}

func filterDotfiles(files []os.FileInfo) (notHiddenFiles []os.FileInfo) {
	for _, f := range files {
		if !strings.HasPrefix(f.Name(), ".") {
			notHiddenFiles = append(notHiddenFiles, f)
		}
	}
	return
}

func generate(pkg, bin, vmname string, build bool) (err error) {
	var name string

	// Remove all channel-related logic
	name = pkg
	log.Println("Using package:", name)

	if !isPackageExists(name) {
		s := "Package " + name + " does not exists"
		err = errors.New(s)
		log.Println(s)
		return
	}

	path, err := nixPath(name)
	if err != nil {
		log.Println("Cannot find nix path")
		return
	}

	path = strings.TrimSpace(path)

	files, err := ioutil.ReadDir(path + "/bin/")
	if err != nil {
		log.Println(err)
		return
	}

	if bin == "" && len(files) != 1 {
		fmt.Println("There's more than one binary in */bin")
		fmt.Println("Files in", path+"/bin/:")
		for _, f := range files {
			fmt.Println("\t", f.Name())
		}

		log.Println("Trying to guess binary")
		var found bool = false

		notHiddenFiles := filterDotfiles(files)
		if len(notHiddenFiles) == 1 {
			log.Println("Use", notHiddenFiles[0].Name())
			bin = notHiddenFiles[0].Name()
			found = true
		}

		if !found {
			for _, f := range files {
				parts := strings.Split(pkg, ".")
				if f.Name() == parts[len(parts)-1] {
					log.Println("Use", f.Name())
					bin = f.Name()
					found = true
				}
			}
		}

		if !found {
			log.Println("Cannot guess in */bin, " +
				"you should specify one of them explicitly")
			return
		}
	}

	if bin != "" {
		var found bool = false
		for _, f := range files {
			if bin == f.Name() {
				found = true
			}
		}
		if !found {
			log.Println("There's no such file in */bin")
			return
		}
	} else {
		bin = files[0].Name()
	}

	var appFilename string
	var finalAppName string
	if vmname != "" {
		appFilename = configDir + "/" + vmname + ".nix"
		finalAppName = vmname
	} else {
		appFilename = configDir + "/" + name + ".nix"
		finalAppName = name
	}

	appNixConfig := fmt.Sprintf(template, name, bin)

	err = ioutil.WriteFile(appFilename, []byte(appNixConfig), 0600)
	if err != nil {
		log.Println(err)
		return
	}

	// Update flake.nix for this app
	err = updateFlakeForApp(configDir, finalAppName)
	if err != nil {
		log.Println("Failed to update flake.nix:", err)
		return
	}

	fmt.Print(appNixConfig + "\n")
	log.Println("Configuration file is saved to", appFilename)

	if build {
		_, _, _, err = generateVM(configDir, finalAppName, true)
		if err != nil {
			return
		}
	}

	return
}
func updateFlakeForApp(appvmPath, appName string) error {
	flakeContent := `{
  description = "AppVM Flake";

  inputs = { nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable"; };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";

      # Function to create a NixOS configuration for an app
      makeAppVM = name:
        nixpkgs.lib.nixosSystem {
          inherit system;
          modules = [
            ./${name}.nix
            "${nixpkgs}/nixos/modules/virtualisation/qemu-vm.nix"
            {
              virtualisation.graphics = false;
              virtualisation.qemu.options = [ "-enable-kvm" "-m 8192" ];
            }
          ];
        };
      # Get all .nix files except base.nix and local.nix
      nixFiles = builtins.filter (name:
        builtins.match ".*.nix$" name != null && name != "base.nix" && name
        != "local.nix") (builtins.attrNames (builtins.readDir ./.));

      # Convert filenames to app names (remove .nix extension)
      appNames =
        map (name: builtins.substring 0 (builtins.stringLength name - 4) name)
        nixFiles;

    in {
      nixosConfigurations = builtins.listToAttrs (map (name: {
        name = name;
        value = makeAppVM name;
      }) appNames);
    };
}
`

	return ioutil.WriteFile(appvmPath+"/flake.nix", []byte(flakeContent), 0644)
}

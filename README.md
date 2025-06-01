# Nix application VMs: security through virtualization

Simple application VMs (hypervisor-based sandbox) based on Nix package manager.

Uses one **read-only** /nix directory for all appvms. So creating a new appvm (but not first) is just about one minute.

The home directory of each appvm is inside ~/appvm, so you can easily share
files between the two as and when needed

![appvm screenshot](https://gateway.ipfs.io/ipfs/QmetVp2LRwcy3baxuAjDgBPwv5ych5kRfXeULoNpQAFsaP)

## Installation and Usage

1. Clone this repo.

2. Run `go build` to build the program

3. Run `./appvm generate brave` to generate the config files for Brave.

Note: If you use flakes for NixOS, nix channels probably won't appear when
you run this as a normal user. Run the command as root in this case (it just
needs to create configs, it doesn't do anything else)

4. Run `./appvm start brave` to launch Brave inside the VM. (this can be done
as your user only, you don't need to use root!)

Right now as a proof of concept only Brave is contained inside the config,
however the code can be extended to allow any program you want here.

You can customize local settings in **~/.config/appvm/nix/local.nix**.

Default hotkey to release cursor: ctrl+alt.

### Shared directory

    $ ls appvm/chromium
    foo.tar.gz
    bar.tar.gz

### Close VM

    $ appvm stop chromium

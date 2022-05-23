[![Documentation Status](https://readthedocs.org/projects/appvm/badge/?version=latest)](https://appvm.readthedocs.io/en/latest/?badge=latest)

# Nix application VMs: security through virtualization

Simple application VMs (hypervisor-based sandbox) based on Nix package manager.

Uses one **read-only** /nix directory for all appvms. So creating a new appvm (but not first) is just about one minute.

![appvm screenshot](https://gateway.ipfs.io/ipfs/QmetVp2LRwcy3baxuAjDgBPwv5ych5kRfXeULoNpQAFsaP)

## Installation

See [related documentation](https://appvm.readthedocs.io/en/latest/installation.html).

## Usage

### Search for applications

    $ appvm search chromium

### Run application

    $ appvm start chromium
    $ # ... long wait for first time, because we need to collect a lot of packages

### Synchronize remote repos for applications

    $ appvm sync

You can customize local settings in **~/.config/appvm/nix/local.nix**.

Default hotkey to release cursor: ctrl+alt.

### Shared directory

    $ ls appvm/chromium
    foo.tar.gz
    bar.tar.gz

### Close VM

    $ appvm stop chromium

### Automatic ballooning

Add this command:

    $ appvm autoballoon

to crontab like that:

    $ crontab -l
    * * * * * /home/user/dev/go/bin/appvm autoballoon

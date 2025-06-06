package main

import "fmt"

// You may think that you want to rewrite to proper golang structures.
// Believe me, you shouldn't.

func generateXML(vmName string, network networkModel, gui bool, sound bool,
	vmNixPath, reginfo, img, sharedDir string) string {

	devices := ""

	if gui {
		devices = guiDevices
	}

	if sound {
		devices += soundDevices
	}

	qemuParams := qemuParamsDefault

	if network == networkQemu {
		qemuParams = qemuParamsWithNetwork
	} else if network == networkLibvirt {
		devices += netDevices
	}

	return fmt.Sprintf(xmlTmpl, vmName, vmNixPath, vmNixPath, vmNixPath,
		reginfo, img, sharedDir, sharedDir, sharedDir, devices, qemuParams)
}

var qemuParamsDefault = `
  <qemu:commandline>
    <qemu:arg value='-snapshot'/>
  </qemu:commandline>
`

var qemuParamsWithNetwork = `
  <qemu:commandline>
    <qemu:arg value='-device'/>
    <qemu:arg value='e1000,netdev=net0,bus=pci.0,addr=0x10'/>
    <qemu:arg value='-netdev'/>
    <qemu:arg value='user,id=net0'/>
    <qemu:arg value='-snapshot'/>
  </qemu:commandline>
`

var netDevices = `
    <interface type='network'>
      <source network='default'/>
    </interface>
`

var guiDevices = `
    <!-- Graphical console -->
    <graphics type='spice'>
      <listen type='socket' socket='/tmp/spice.sock'/>
      <gl enable='yes' rendernode='/dev/dri/renderD128'/>
    </graphics>
    <!-- Guest additionals support -->
    <channel type='spicevmc'>
      <target type='virtio' name='com.redhat.spice.0'/>
    </channel>
    <video>
      <model type='virtio' heads='1' primary='yes'>
        <acceleration accel3d='yes'/>
      </model>
    </video>
`

var soundDevices = `
    <sound model='ich9'>
      <codec type='duplex'/>
    </sound>
`

// <!-- Graphical console -->
// <graphics type='spice' autoport='yes'>
//   <listen type='none'/>
//   <gl enable='yes' rendernode='/dev/dri/renderD128'/>
// </graphics>
// <!-- Guest additionals support -->
// <channel type='spicevmc'>
//   <target type='virtio' name='com.redhat.spice.0'/>
// </channel>
// <video>
//   <model type='virtio' heads='1' primary='yes'>
//     <acceleration accel3d='yes'/>
//   </model>
// </video>

var xmlTmpl = `
<domain type='kvm' xmlns:qemu='http://libvirt.org/schemas/domain/qemu/1.0'>
  <name>%s</name>
  <memory unit='GiB'>8</memory>
  <currentMemory unit='GiB'>4</currentMemory>
  <vcpu>8</vcpu>
  <os>
    <type arch='x86_64'>hvm</type>
    <kernel>%s/kernel</kernel>
    <initrd>%s/initrd</initrd>
    <cmdline>loglevel=4 init=%s/init %s</cmdline>
  </os>
  <features>
    <acpi></acpi>
  </features>
  <clock offset='utc'/>
  <on_poweroff>destroy</on_poweroff>
  <on_reboot>restart</on_reboot>
  <on_crash>destroy</on_crash>
  <devices>
    <!-- Fake (because -snapshot) writeback image -->
    <disk type='file' device='disk'>
      <driver name='qemu' type='qcow2' cache='writeback' error_policy='report'/>
      <source file='%s'/>
      <target dev='vda' bus='virtio'/>
    </disk>
    <!-- filesystems -->
    <filesystem type='mount' accessmode='passthrough'>
      <source dir='/nix/store'/>
      <target dir='nix-store'/>
      <readonly/>
    </filesystem>
    <filesystem type='mount' accessmode='mapped'>
      <source dir='%s'/>
      <target dir='xchg'/> <!-- workaround for nixpkgs/nixos/modules/virtualisation/qemu-vm.nix -->
    </filesystem>
    <filesystem type='mount' accessmode='mapped'>
      <source dir='%s'/>
      <target dir='shared'/> <!-- workaround for nixpkgs/nixos/modules/virtualisation/qemu-vm.nix -->
    </filesystem>
    <filesystem type='mount' accessmode='mapped'>
      <source dir='%s'/>
      <target dir='home'/>
    </filesystem>
    %s
  </devices>
  %s
</domain>
`

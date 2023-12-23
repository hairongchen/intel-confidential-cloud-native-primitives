# Confidential VM Customization Tool

This tool is used to customize the confidential VM guest including guest image,
config, OVMF firmware etc.

## 1. Overview

The confidential VM guest can be customized including follows:

![](/docs/cvm-customizations.png)

| Name | Type/Scope | Description |
| ---- | ---------- | ----------- |
| Launch Identity | Config | MROwner, MRConfig, MROwnerConfig |
| VM Configuration | Config | vCPU, memory, network config |
| Secure Boot Key | OVMF | the PK/DB/KEK for secure boot or Linux MoK |
| Config Variable | OVMF | the configurations in variable |
| Grub | Boot Loader | Grub kernel command, Grub modules |
| initrd | Boot Loader | Customize build-in binaries |
| IMA Policy | OS | Policy via loading systemd |
| Root File System | OS | RootFS customization |

## 2. Design

It is based on the [cloud-init](https://cloudinit.readthedocs.io/en/latest/)
framework, and the whole flow was divided into three stages:

- **Pre Stage**: prepare to run cloud-init. It will collect the files for target
  image, meta-data/x-shellscript/user-data for cloud-init's input.
- **Cloud-Init Stage**: it will run cloud init in sequences of
  - Generate meta files via `cloud-init make-mime`
  - Generate `ciiso.iso` via `genisoimage`
  - Run cloud-init via `virt-install`
- **Post Stage**: clean up and run post check

![](/docs/cvm-image-rewriter-flow.png)

## 2. Run

### 2.1 Customize

```
$ ./run.sh -h
Usage: run.sh [OPTION]...
Required
  -i <guest image>          Specify initial guest image file
Optional
  -t <number of minutes>    Specify the timeout of rewriting, 3 minutes default,
                            If enabling ima, recommend timeout >6 minutes
  -s <connection socket>    Default is connection URI is qemu:///system,
                            if install libvirt, you can specify to "/var/run/libvirt/libvirt-sock"
                            then the corresponding URI is "qemu+unix:///system?socket=/var/run/libvirt/libvirt-sock"
```

**_NOTE_**:

1. If want to skip to run specific plugins at `pre-stage` directory, please create
a file named as `NOT_RUN` at the plugin directory. For example:
    ```
    touch pre-stage/01-resize-image/NOT_RUN
    ```

2. Please make sure read permission on `/boot/vmlinuz-$(uname-r)`.

3. It can run without installing `virt-daemon/libvirt`, by default the URI of `virt-install`
   is `qemu:///system`

4. If it is running with `libvirt/virt-daemon` hypervisor, then:

  - In file `/etc/libvirt/qemu.conf`, make sure `user` and `group` is `root` or
    current user.
  - If need customize the connection URL, you can specify via `-s` like `-s /var/run/libvirt/libvirt-sock`,
    please make sure current user belong to libvirt group via following commands:
    ```
    sudo usermod -aG libvirt $USER
    sudo systemctl daemon-reload
    sudo systemctl restart libvirtd
    ```


### 2.2 Run Test

```
$ ./qemu-test.sh -h
Usage: qemu-test.sh [OPTION]...
Required
  -i <guest image>          Specify initial guest image file
```

## 3. Plugin

### 3.1 Existing Plugins

There are following customization plugins in Pre-Stage:

| Name | Descriptions |
| ---- | ------------ |
| 01-resize-image | Resize the input qcow2 image |
| 02-motd-welcome | Customize the login welcome message |
| 03-netplan | Customize the netplan.yaml |
| 60-initrd-update | Update the initrd image |
| 98-ima-enable-simple | Enable IMA (Integrity Measurement Architecture) feature |

### 3.1 Design a new plugin

TBD
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

packer {
  required_plugins {
    kubernetes = {
      version = ">= 1.0.0"
      source = "github.com/robbert229/kubernetes"
    }

    qemu = {
      source  = "github.com/hashicorp/qemu"
      version = "~> 1"
    }
  }
}

data "kubernetes-secret" "credentials" {
  name = "example-credentials"
  namespace = "default"
}

locals {
  community_repo = "http://dl-cdn.alpinelinux.org/alpine/v3.21/community"
  cpus = "1"
  disk_size = "10G"
  iso_checksum = "f28171c35bbf623aa3cbaec4b8b29297f13095b892c1a283b15970f7eb490f2d"
  iso_checksum_type = "sha256"
  iso_download_url = "https://dl-cdn.alpinelinux.org/alpine/v3.21/releases/x86_64/alpine-virt-3.21.3-x86_64.iso"
  iso_local_url = "../../iso/alpine-virt-3.21.3-x86_64.iso"
  memory = "1024"
  ssh_username     = data.kubernetes-secret.credentials.data["ssh_username"]
  ssh_password     = data.kubernetes-secret.credentials.data["ssh_password"]
  root_password    = data.kubernetes-secret.credentials.data["root_password"]
  vm_name = "alpine-3.21.3-x86_64"
}

source "qemu" "alpine" {
  boot_command = [
    "root<enter><wait>",
    "ifconfig eth0 up && udhcpc -i eth0<enter><wait5>",
    "wget http://{{ .HTTPIP }}:{{ .HTTPPort }}/answers<enter><wait>",
    "export ERASE_DISKS=/dev/vda<enter>",
    "export USEROPTS='-a -u -g audio,video,netdev ${local.ssh_username}'<enter>",
    "export USERSSHKEY='http://{{ .HTTPIP }}:{{ .HTTPPort }}/ssh.keys'<enter>",
    "setup-alpine -f $PWD/answers<enter><wait5>",
    "${local.root_password}<enter><wait>",
    "${local.root_password}<enter><wait30>",
    "mount /dev/vda3 /mnt<enter>",
    "chroot /mnt<enter>",
    "echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config<enter>",
    "su ${local.ssh_username}<enter>",
    "passwd<enter>${local.ssh_password}<enter>${local.ssh_password}<enter>",
    "exit<enter>",
    "umount /mnt<enter>",
    "reboot<enter>"
  ]
  boot_wait = "10s"
  communicator = "ssh"
  disk_size = local.disk_size
  format = "qcow2"
  headless = false
  http_directory = "http"
  iso_checksum = "${local.iso_checksum_type}:${local.iso_checksum}"
  iso_urls = [
    local.iso_local_url,
    local.iso_download_url,
  ]
  shutdown_command = "/sbin/poweroff"
  ssh_timeout = "60m"
  ssh_username = "root"
  ssh_password = local.root_password
  vm_name = local.vm_name
}

build {
  sources = ["source.qemu.alpine"]

  provisioner "shell" {
    inline = [
      "echo ${local.community_repo} >> /etc/apk/repositories",
      "apk update",
      "apk upgrade",
      "apk add sudo",
      "echo '${local.ssh_username} ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers.d/${local.ssh_username}",
      "adduser ${local.ssh_username} wheel",
      "sed -i '/^PermitRootLogin yes$/d' /etc/ssh/sshd_config",
    ]
  }
}
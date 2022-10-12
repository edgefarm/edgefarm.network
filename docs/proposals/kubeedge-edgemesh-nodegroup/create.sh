#!/bin/bash
PROVISIONER=${1:-docker}
sudo --preserve-env=HOME /home/armin/bin/talosctl cluster create --kubernetes-version 1.22.9 --masters 2 --workers 2 --provisioner ${PROVISIONER}


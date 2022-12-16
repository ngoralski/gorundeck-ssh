#!/usr/bin/env bash

#go build cmd/gorundeck-ssh/main.go

scp contents/build_linux root@rundeck:/var/lib/rundeck/libext/cache/gorundeck-ssh-1.0-plugin/go_build_main_go_linux

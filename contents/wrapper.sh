#!/usr/bin/env bash

plugin_dir=${BASH_SOURCE[0]%/*}

binary=${plugin_dir}/go_build_main_go_linux

chmod +x ${binary}
${binary}
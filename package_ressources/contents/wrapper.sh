#!/usr/bin/env bash

plugin_dir=${BASH_SOURCE[0]%/*}

binary=${plugin_dir}/gorundeck-ssh

if [[ ! -x "${binary}" ]]
then
  chmod +x "${binary}"
fi

${binary}
# GORundeck-ssh

## Description
It's a plugin for rundeck to run ssh command via a bastion.\
The original software is https://github.com/rundeck-plugins/openssh-bastion-node-execution/. \
The main issue I encounter using it is it does not support sudo interaction.\
If you are using it and a sudo command not defined as password less is scheduled in rundeck it will wait for
password input until timeout but the remote command remain on the server.   

The main goal of this plugin is to reproduce almost the same functionality to run remote command via a bastion and if
a sudo command with a password request is detected to exit properly from the remote host and return an error message.

## Installation
Deploy the package like others rundeck plugins.


## Usage
Fill the form with all expected values


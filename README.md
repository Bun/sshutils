# sshutils

Parallel SSH utilities with support for Ansible-like inventory files.

- pssh: execute commands
- prun: execute local scripts/binaries remotely (requires SFTP support)

Requirements:

- A running ssh-agent with the proper keys loaded
- The public key of the server is in `~/.ssh/known_hosts`


## Examples

Run an apt update and check the effects:

    pssh -i hosts.yaml all sudo env DEBIAN_FRONTEND=noninteractive apt-get update
    pssh -i hosts.yaml all sudo env DEBIAN_FRONTEND=noninteractive apt-get --simulate dist-upgrade

Run a local script (or binary) remotely:

    $ cat myscript.sh
    #!/bin/sh
    echo "$@"
    $ prun -i hosts.yaml all ./myscript.sh hello world
    cat> hello world
    dog> hello world


## TODO

- Additional documentation
- pscp: parallel transfer
- Load hostnames from `~/.ssh/config`
- (Non-)inventory support
    - Groups
    - Assume hostname if name not in inventory and no glob char


## Ansible inventories

Ansible YAML-based configurations are supported, for example:

    all:
      hosts:
        # Set required variables per host:
        dog: {ansible_host: dog.example.com}
        cat: {ansible_host: cat.example.org}
      vars:
        # Variables can also be set globally
        ansible_user: admin

Supported variables:

- `ansible_host`
- `ansible_port`
- `ansible_user`

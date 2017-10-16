# sshutils

Parallel SSH utilities with support for Ansible-like inventory files.

- prun: execute local scripts/binaries remotely
- pssh: execute commands

Requirements:
- A running ssh-agent with the proper keys loaded
- The public key of the server is in ``~/.ssh/known_hosts``


## TODO

- Documentation, examples
- Work in progress: porting Paramiko/Python version to Go
- pscp: parallel transfer
- ~/.ssh/config support
- (Non-)inventory support
    - Groups
    - Assume host if name not in inventory and no glob char

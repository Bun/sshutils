# sshutils

Parallel SSH utilities with support for Ansible-like inventory files.

- prun: execute local scripts/binaries remotely
- pssh: execute commands

TODO:
- Semi-realtime output
- Documentation, examples
- Flags
- Work in progress: porting Paramiko/Python version to Go
- pscp: parallel transfer
- ~/.ssh/config support
- (Non-)inventory support
    - Groups
    - Assume host if name not in inventory and no glob char

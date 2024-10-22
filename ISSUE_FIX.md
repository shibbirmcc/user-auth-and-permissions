### **TestContainer test execution in a podman environment**

if the tests contains testcontainer and the machine has podman in it, then there is a common problem occurs :

```
Failed to start container: create container: container start: Error response from daemon: error configuring network namespace for container 9158929e7f717f72159279dedcbf7591aa0582e08f378ae45fb1d12b36a486ca: CNI network "reaper_default" not found: could not start container: creating reaper failed
```

Follow the below procedure to solve the issues:

```
podman network ls
podman network create reaper_default
```

Podman uses CNI (Container Network Interface) configurations located in directories like /etc/cni/net.d/ or ~/.config/cni/net.d/. You may need to ensure that the configuration file for reaper_default is present and correctly set up. check if the version is different if the file already exists, then just change the version to 0.4.0

```json
{
  "cniVersion": "0.4.0",
  "name": "reaper_default",
  "plugins": [
    {
      "type": "bridge",
      "bridge": "cni0",
      "isGateway": true,
      "ipMasq": true,
      "ipam": {
        "type": "host-local",
        "subnet": "10.88.0.0/16"
      }
    },
    {
      "type": "firewall"
    },
    {
      "type": "portmap",
      "capabilities": {
        "portMappings": true
      }
    }
  ]
}

```

After creating or editing the CNI configuration, restart Podman to apply the changes:

```
systemctl --user restart podman.socket
```

If you continue to have issues, it may help to disable the Ryuk reaper, which is used for resource cleanup but may not be fully compatible with Podman:
Set the environment variable before running your tests
```
export TESTCONTAINERS_RYUK_DISABLED=true
```

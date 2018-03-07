## installation

```bash
$ go get -u github.com/carterjones/instanceinfo/...
```

## usage

```bash
$ instanceinfo <search-term>
```

This tool searches against the following pieces of information found on each AWS instance.

- ImageID
- InstanceID
- InstanceType
- LaunchTime
- PrivateIPAddress
- PublicIPAddress
- Name

When any matches are found, instance information is printed.

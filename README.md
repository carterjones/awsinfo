## installation

```bash
$ go get -u github.com/carterjones/instanceinfo/...
```

## usage

### instanceinfo

```bash
$ instanceinfo <search-term>
```

This tool searches against the following pieces of information found on each AWS instance.

- Image ID
- Instance ID
- Instance Type
- Launch Time
- Private IP Address
- Public IP Address
- Name

When any matches are found, instance information is printed.

### elbinfo

```bash
$ elbinfo <search-term>
```

This tool searches against the following pieces of information found on each classic ELB.

- Name
- DNS name
- Any IPs resolved from the DNS name

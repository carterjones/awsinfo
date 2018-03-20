## installation

```bash
$ go get -u github.com/carterjones/awsinfo/...
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

### r53info

```bash
$ r53info <search-term>
```

This tool searches against the following pieces of information found in Route53.

- Zone
- Name
- Value(s)
- Alias Target

# vouch4cluster 🎟☁️ 

vouch4cluster is a tool for running voucher against all of the images running in a cluster or deployment.

## Installing

Install using:

```
$ go get -u github.com/Shopify/vouch4cluster
``` 
## Configuration

vouch4cluster is configured using either json, yaml, or toml. By default, vouch4cluster loads from `~/.vouch4cluster.{json,yaml,toml}`, but you can also specify the configuration to read from with the `--config` flag.

| Group        | Key           | Description                                        |
| :----------- | :------------ | :------------------------------------------------- |
| `voucher`    | `hostname`    | The address of the Voucher instance to connect to. |
| `voucher`    | `username`    | The username to connect as.                        |
| `voucher`    | `password`    | The password to authenticate with.                 |

For example:

```json
{
   "voucher": {
       "hostname": "https://<voucher address>",
       "username": "<username>", 
       "password": "<password>"
   }
}
```

## Using

### Attest all images in the current Kubernetes context

To attest all images in the current Kubernetes context, use: 

```
$ vouch4cluster kube
```

This will query kubernetes for all of the active images, and then run each image through voucher.

### Attest all images from a file

If you have a list of images that need to be attested, you can put them in a file, newline separated,
and pass that file to vouch4cluster. By default, vouch4cluster will read from standard input.

```
$ vouch4cluster reader < input
```

This will iterate through each line in the file, and run each image through voucher.


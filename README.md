# vouch4cluster 🎟☁️ 

vouch4cluster is a tool for running voucher against all of the images running in a cluster or deployment.

## Installing

Install using:

```
$ go get -u github.com/Shopify/vouch4cluster
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


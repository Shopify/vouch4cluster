# vouch4cluster 🎟☁️ 

vouch4cluster is a tool for running voucher against all of the images running in a cluster or deployment.

## Installing

Install using:

```
$ go get -u github.com/Shopify/vouch4cluster
``` 

## Using

Currently, vouch4cluster only reads newline separated image references from standard input. 

```
$ vouch4cluster < input.txt
```


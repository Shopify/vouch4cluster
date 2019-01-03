# How to contribute

Thank you for considering contributing to vouch4cluster!

## Getting Started

- Review this document and the [Code of Conduct](CODE_OF_CONDUCT.md).

- Setup a [Go development environment](https://golang.org/doc/install#install) 
if you haven't already.

- Get the Shopify version the project by using Go get:

```
$ go get -u github.com/Shopify/vouch4cluster
```

- Fork this project on GitHub. :octocat:

- Setup your fork as a remote for your project:

``` 
$ cd $GOPATH/src/github.com/Shopify/vouch4cluster
$ git remote add <your username> <your fork's remote path>
```

## Work on your feature

- Create your feature branch based off of the `master` branch. (It might be 
worth doing a `git pull` if you haven't done one in a while.)

```
$ git checkout master
$ git pull
$ git checkout -b <the name of your branch>
```

- Code! :keyboard:

    - Please run `go fmt` and `golint` while you work on your change, to clean 
up your formatting/check for issues.

- Push your changes to your fork's remote:

```
$ git push -u <your username> <the name of your branch>
```

## Send in your changes

- Sign the [Contributor License Agreement](https://cla.shopify.com).

- Open a PR against Shopify/vouch4cluster!

TransIp Terraform Provider
==================

[![Build Status](https://travis-ci.org/demeesterdev/terraform-provider-transip.svg?branch=master)](https://travis-ci.org/demeesterdev/terraform-provider-transip)

General Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.12.x (to build the provider plugin)

Windows Specific Requirements
-----------------------------
- [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)
- [Git Bash for Windows](https://git-scm.com/download/win)

For *GNU32 Make*, make sure its bin path is added to PATH environment variable.*

For *Git Bash for Windows*, at the step of "Adjusting your PATH environment", please choose "Use Git and optional Unix tools from Windows Command Prompt".*

Building The Provider
---------------------

build the provider
$ make build

Using The Provider
------------------

```
# configure the TransIP Provider
provider "transip" {
    account_name = "..."
    private_key_path = "..."
}

# retrieve domain data
data "transip_domain" "example" {
  name = "thijsdemeester.nl"
}

# use name_servers
output "name_servers" {
  value = data.transip_domain.example.name_servers
}
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.12+ is **required**). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-transip
...
```

In order to run the unit tests for the provider, you can run:

```sh
$ make test
```

The majority of tests in the provider are Acceptance Tests - which manage real resources at TransIP. It's possible to run the entire acceptance test suite by running `make testacc` - however it's likely you'll want to run a subset, which you can do using a prefix, by running:

```
make testacc TESTARGS='-run=TestAccDataSourceDomain'
```

The following Environment Variables must be set in your shell prior to running acceptance tests:

- `TRANSIP_ACCOUNT_NAME`
- `TRANSIP_ACCOUNT_KEY`
- `TRANSIP_TEST_DOMAIN`

**Note:** Acceptance tests edit real resources at TransIP. Make sure the domain you set with `TRANSIP_TEST_DOMAIN` is safe to edit.


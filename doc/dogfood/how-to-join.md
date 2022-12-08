# How to join our governance community

Welcome and thank you for your interest in gov4git!

Our community is governed by a dogfood release of our own software — gov4git.

We welcome anyone to join and take our governance system for a spin.

In order to join, please follow these steps:

## Quick setup

## Prepare your home repositories

- Create two empty repositories, which will host your public and private identity for the purposes of participating in community governance.

  Your public repository should be publically readable. For instance, mine is hosted on GitHub and its URL is [https://github.com/petar/gov4git.public.git](https://github.com/petar/gov4git.public.git).

  Your private repository should NOT be publically readable. For instance, mine is [https://github.com/petar/gov4git.private](https://github.com/petar/gov4git.private).

  Decide what is the name of the branch (inside your repos) that will host your gov4git-related data. Using `main` is a good convention.

### Install the client

The client is a Go application. Start by making sure the [Go language is installed](https://go.dev/doc/install) on your machine. Then use the standard workflow for installing Go applications:
```sh
go install github.com/gov4git/gov4git/gov4git@latest
```

Verify that you have the command `gov4git` in your environment:
```sh
gov4git -h
```

### Configure the client

Before you can use the client, you need to configure it with information about your own identity, as well as the community that you plan to interact this.

You will need the following information handy:
- The git URL of the community's public repo, which is `https://github.com/gov4git/governance.git`
- The git branch in the community's public repo, where governance proceedings live, which in our case is `main`
- The git URL of your identity's public repo, which in my case is `git@github.com:petar/gov4git.public.git` Make sure you use a writable URL, since your client will perform write operations to your own public repo.
- The git branch in your identity's public repo where your identity resides. In my case, this is `main`.
- The git URL of your identity's private repo, which in my case is `git@github.com:petar/gov4git.private.git` Make sure you use a writable URL here as well.
- The git branch in your identity's private repo. In my case, this is `main`.
- The location of your SSH private keys, in my case `/Users/petar/.ssh/id_rsa`, which are necessary so that the client can write to your identity's public and private repos.

Place the following configuration in a file at path `~/.gov4git/config.json`, making sure to replace the example parameters with the ones that apply in your case:

```json
{
     "auth" : {
          "git@github.com:petar/gov4git.public.git": {
               "ssh_private_keys_file": "/Users/petar/.ssh/id_rsa"
          },
          "git@github.com:petar/gov4git.private.git": {
               "ssh_private_keys_file": "/Users/petar/.ssh/id_rsa"
          }
     },
     "gov_public_url": "https://github.com/gov4git/governance.git",
	"gov_public_branch": "gov",
	"member_public_url": "git@github.com:petar/gov4git.public.git",
	"member_public_branch": "main",
	"member_private_url": "git@github.com:petar/gov4git.private.git",
	"member_private_branch": "main"
}
```

You are now ready to use the `gov4git` client.

### Initialize your identity

Prior to engaging with the community, you need to initialize your digital identity, so we can ensure no one can impersonate you. Run this command:
```sh
gov4git init-id
```
This will generate a fresh pair of [ED25519](https://ed25519.cr.yp.to/) public and private signature keys and stores them in your public and private repos, respectively.

You should be all set to participate in governance.

## Request to join

- Request to join the gov4git community using [this form](https://github.com/gov4git/gov4git/issues/new?assignees=petar&labels=community&template=join.yml&title=I%27d+like+to+join+this+project%27s+community).

   You will receive a confirmation once you've been added to the community members. 

   At this point, you can start participating in our governance by following [these instructions](how-to-participate.md).

If you run into problems, [let us know](https://github.com/gov4git/gov4git/issues/new/choose).

# Welcome to the gov4git community!

You are here because you [requested to join](https://github.com/gov4git/gov4git/issues/new?assignees=petar&labels=community&template=join.yml&title=I%27d+like+to+join+this+project%27s+community) our community, and you have notified that you've been added to our list of members. Awesome!

You can now take part in the community governance of the gov4git project. You will be able to cast votes in prioritizing outstanding issues and pull requests, as well as donate voting credits to other users you want to empower.

All governance operations and proceedings of our community are recorded in a transparent and verifiable manner in our [GitHub repository](https://github.com/gov4git/gov4git) in a [dedicated branch](https://github.com/gov4git/gov4git/tree/gov).

You can interact with our community governance and exercise your member rights (such as voting) using our command-line client application. A slick web UI is in the making, but for now we will assume that you are not the patient kind, and you are willing to wrestle the terminal.

## Quick setup

Let's first set you up by installing the gov4git command-line tool (the client).

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
- The git URL of the community's public repo, which is `https://github.com/gov4git/gov4git.git`
- The git branch in the community's public repo, where governance proceedings live, which in our case is `gov`
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
     "gov_public_url": "https://github.com/gov4git/gov4git.git",
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

## Take part in governance

As a community member, you can take part in prioritizing issues by spending your voting credits for or against any open issue.

First, let's list all issues that are open for voting. Run the command:
```sh
gov4git ballot list
```

Suppose we are interested in issue #6. Let's take a look at how many votes this issue has collected so far. Run the command:

```sh
gov4git ballot show-open --name=issue/6
```

Suppose that we would like to promote this issue. 

XXX

```sh
gov4git ballot vote --name issue/6 --choices issue-6 --strengths=+3.0
```
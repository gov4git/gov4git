
# Deployment guide

This guide will help you deploy governance for a GitHub project repository.

## Create a governance identity for your project

1. Begin by creating two empty repositories — one public, one private — for the governance system. Suppose `GOV_PUB_REPO` and `GOV_PRIV_REPO` are the HTTPS URLs of the public and private repositories, respectively.

2. Download and install `gov4git` on your local machine:


3. From your local machine, initialize the governance system:

     ```bash
     gov4git init-gov
     ```



## Create a governance environment

Add a new environment to your GitHub project repository, named `governance`.

Using the GitHub UI, add an environment variable `GOV4GIT_RELEASE` pointing to the desired release of gov4git. For instance,

```GOV4GIT_RELEASE=v1.1.4```


# Welcome to the gov4git community!

You are here because you [requested to join](https://github.com/gov4git/gov4git/issues/new?assignees=petar&labels=community&template=join.yml&title=I%27d+like+to+join+this+project%27s+community) our community, and you have received a confirmation that you are a member. Great!

You can now take part in the community governance of the gov4git project. For instance, you could cast votes in prioritizing outstanding issues, as well as donate voting credits to other users you want to empower.

All governance operations and proceedings of our community are recorded in a transparent and verifiable manner in our [governance repository](https://github.com/gov4git/governance).

You can interact with our community governance and exercise your member rights (such as voting) using a dedicated command-line client application. A slick web UI is in the making, but for now we assume that you are not the patient kind, and are willing to wrestle the terminal.

## Take part in governance

As a community member, you can take part in prioritizing issues by spending your voting credits for or against any open issue.

First, let's list all issues that are open for voting. Run the command:
```sh
gov4git ballot list-open
```

Suppose we are interested in issue #6. Let's take a look at how many votes this issue has collected so far. Run the command:

```sh
gov4git ballot show-open --name=issue/6
```

Say, we would like to promote this issue. Every community member, like yourself, has a balance of "voting credits" which can be spent towards up- or down-voting individual issues. Your balance of voting credits is replenished when join and periodically afterwards, currently by the community organizer.

Let's find out how many voting credits we have at our disposal. Run this command:

```sh
gov4git balance get --user=petar --key=voting_credits
```

Note that your community username (in this example, mine is `petar`) will be your GitHub username.

You can now spend your available credits by up-voting or down-voting various open issues, like so:

```sh
gov4git ballot vote --name issue/6 --choices issue-6 --strengths=+3.0
gov4git ballot vote --name issue/8 --choices issue-8 --strengths=-1.0
```

There is no limit to how many votes you cast, as long as you don't exceed your voting credit allowance. Individual votes are cumulative.

Your votes will not take effect immediately. They will be incorporated in the official tallies periodically by the community organizer.

You can always check the current tally for any open (or closed) issue. Use this command:

```sh
gov4git ballot show-open --name=issue/6
```

This will display the current tally for issue #6, along with a list of everyone's votes. (In this case, we are using a voting strategy whereby everyone's votes are public.) When your vote is "counted", it will appear in the tally.

## What next?

We are dogfood-ing a brand new product which is exploring a few different frontiers at once:
- What is the right model and UX for governance applications?
- Can decentralized social applications be built entirely on the git protocol?
- Can a blockchain-like immutable and verifiable history of governance be implemented by a networked community of git repos?

Our governance tool `gov4git` supports an additional range of functionalities which we plan to document soon. In the meantime, we would like to start with small steps and become comfortable with a few core workflows, such as the ones described above.

Please, don't hesitate to submit feedback in the form of [issues](https://github.com/gov4git/gov4git/issues/new/choose). When submitting bugs, be sure to attach the relevant logs, which can be obtained by using the `-v` command-line option with any command, i.e.
```sh
gov4git -v ...
```

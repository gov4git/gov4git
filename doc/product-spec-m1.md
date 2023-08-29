## Product spec (Milestone 1)

The baseline application features (and mental model) are as follows:

- The community's governance state is kept in a repo owned and operated by the community organizer
- The community is associated with a source code repo (e.g. the Plurality Book repo), which is distinct from the governance repo, also owned and operated by the organizer

- The organizer can add new GitHub users to the community, making them "members"
	- To request membership, GitHub users fill out a custom GitHub issue form (found in the GitHub source code repo of the community)
	- Membership can either be granted automatically (using a GitHub action) or can be deferred for manual approval by the organizer via the command-line client (are there better UI approaches here?)

- Each member is associated with a persistent balance of "voting credits", initially zero
- The organizer can modify (add/remove) any member's voting credit balance, using the command line client at any time

- Prioritizing issues:
	- The organizer can create issues — in the source code GitHub repo of the community — which are designated for prioritization polling
		- To create an issue for prioritization, the organizer creates a regular GitHub issue that is labelled "prioritize"
	- Members can vote on any open issue, continuously (at any time), multiple-times using their voting credits via the desktop app
		- When a vote is processed, voting credits are withdrawn from the voter's account
	- The current tally (priority of the issue) is computed as follows. Let $p_u^t$ be voting credits user $u$ spends on an upvote at time $t$. Similarly, let $n_u^t$ be the voting credits user $u$ spends on a downvote at time $t$. The priority of the issue being voted on equals:
     $$\sum_u \sqrt{ \sum_t p_u^t - n_u^t }$$
	- The organizer can close or cancel an issue (with prioritization polling) at any time
		- Closing an issue rejects any further votes, and records the issue, its tally and its voting history
			- To close an issue, the organizer simply closes the respective issue on GitHub
		- Cancelling an issue returns all voting credits captured in the issue to the original voters
			- To cancel an issue, the organizer ... closes the GitHub issue with a comment "cancelled"?
	- The organizer can freeze/unfreeze an issue as long as it is not closed. While an issue is frozen, votes on it are rejected.
	- Once an issue is closed, it cannot be re-opened, cancelled, frozen or unfrozen.
	- All users (members and organizers) can view currently open and closed issues, as well as their current voting tallies in the desktop app.

- Community members use the desktop app to:
	- Create gov4git identities for themselves
		- Identity creation requires the user to input their: GitHub Access Token and GitHub username?
			- Or simply requires them to authenticate the desktop app as a GitHub App? (I think the former. Please, correct me.)
	- View their voting credit balance on every screen of the desktop app
	- View open/closed issues, ordered by current priority (... what else?), filtered by ... (GitHub issue labels?)
	- Cast a vote on any open issue

Deployment integrations:

- To deploy governance to an existing source code repo on GitHub (e.g. the Plurality Book), the deployer/organizer must:
	- Provision a public and a priavte git repos for the community's governance state
	- Install a package of GitHub actions in the source code repo, for short the "automation actions"
		- Automation actions capture GitHub events (issue create, issue close, issue label change) and trigger corresponding actions in governance (create a poll, close a poll, etc.)
		- Automation actions process requests to join the community, in the form of GitHub issues with a specific template
		- Automation actions also run an hourly (?) "sweep" job which collects votes from members (from their public repos, specifically) that were cast since the last sweep. Collected votes are recorded in the community state and reflected where they apply.

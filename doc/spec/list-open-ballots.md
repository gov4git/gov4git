
## SUMMARY

- function "list open ballots"
- implemented in [ballot/list.go](../../proto/ballot/ballot/list.go)

## CONTEXT

At any point in time, there can be zero or more ballots that are "open" for voting.

Each ballot is specified by:
- _path_: a valid file-system path, e.g. `issue/3`, which uniquely names this ballot
- _title_ and _description_: human-readable strings
- _choices_: a set of strings, representing abstract voting choices
- _strategy_: a string that specifies the type of voting algorithm used, e.g. `priority_poll`
- _participants_: the group of community users that can participate in the ballot. Groups are named sets of community members, defined by the community organizer. Group `everybody` contains all community members by default. The organizer can define custom groups, like `english-translators`, for instance.

When a new ballot is created, its specification (above) is placed in an "advertisement" file in JSON format, and the file is placed inside the repo on the path:

     `ballots/open/BALLOT_PATH/ballot_ad.json`

Where `BALLOT_PATH` is the unique path that names this ballot. Here is an example ballot advertisement:

```json
{
     "community":{
          "Repo": "git@github.com:gov4git/gov4git.git",
          "Branch":"gov",
     },
     "path":["issue","8"],
     "title":"Issue 8",
     "description":"https://github.com/gov4git/gov4git/issues/8",
     "choices":["issue-8"],
     "strategy":"priority_poll",
     "participants_group":"everybody","parent_commit":"1b07688fe71d5808b2bf1474f69284e0f817ea8e"
}
```

Generally, we will be interested in the fields `path`, `title`, `description`, `choices`, `strategy` and `participants_group`. (The other fields are of no interest to UI considerations.)

## PSEUDOCODE

Listing open ballots entails returning a list of ballot advertisements for all open ballots.

1. clone the governance repo of the community
2. find all files whose path meets the pattern `ballots/open/**/ballot_ad.json`
3. parse the found JSON advertisement files and return them as a list

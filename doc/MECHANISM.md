
## Global

M_TOTAL is the number of credits within a global matching fund

K is a global parameter used for computing contribution payouts

K = (SUM over issues J of M(J)) / M_TOTAL

L is a global parameter for adjusting the relative cost of bets, set by admins


## Issues

Every user I can cast a "prioritization vote", W(I, J), on every issue, J

P(I, J) = credits spent by user I on issue J = W(I, J)^2

Only positive prioritization votes and (partial or complete) vote withdrawals are allowed (no negative votes)

Credits spent on prioritization are placed in a "bounty" pool, specific to the issue

CAP(J) = credits escrowed in the bounty of issue J = SUM over users I of P(I, J)

QP(J) = quadratic priority of issue J = (SUM over users I of W(I, J))^2

## PRs

Every user I can cast a "contribution vote", V(I, J), and make an optional "bet" on every PR, J:

Q(I, J) = credits spent by user I on PR addressing issue J

Q(I, J) = L * V(I, J)^2 - V, with bet (shouldn't this be a "+" sign?)

Q(I, J) = L * V(I, J)^2, without bet

Q(I, J) is always positive

S(I, J) = direction of vote, +1/-1 for "for"/"against", respectively

PR J is ACCEPTED if VERDICT(J) > 0, and otherwise REJECTED

VERDICT(J) = SUM over users I of S(I, J) * SQRT(Q(I, J))

CC(J) = all credits spent on contribution votes = SUM over users I of Q(I, J)

Credits spent on contribution votes, CC(J), are placed into the global matching fund, M_TOTAL, after paying out winning bets

A winning bet, associated with a vote V(I, J), is rewarded 2 * V(I, J)

## Contribution payout

If a PR J is accepted, the contributor is rewarded a "contribution payout", CP(J), sourced from CC(J)

CP(J) = contribution payout for issue J = CAP(J) + K * M(J)

M(J) = QP(J) - CAP(J)

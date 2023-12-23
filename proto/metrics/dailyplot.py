import matplotlib.pyplot as plt
import matplotlib.dates as mdates
import numpy as np
from datetime import date
# from matplotlib.gridspec import GridSpec

x = [
    date(2004, 11, 1),
    date(2004, 11, 2),
    date(2004, 11, 3),
    date(2004, 11, 4),
    date(2004, 11, 5),
    date(2004, 11, 6),
    date(2004, 11, 7),
    date(2004, 11, 8),
    date(2004, 11, 9),
    date(2004, 11, 10),
    date(2004, 11, 11),
    date(2004, 11, 12),
    date(2004, 11, 13),
    date(2004, 11, 14),
    date(2004, 11, 15),
    date(2004, 11, 16),
    date(2004, 11, 17),
    date(2004, 11, 18),
    date(2004, 11, 19),
    date(2004, 11, 20),
    date(2004, 11, 21),
    date(2004, 11, 22),
    # date(2004, 11, 23),
    # date(2004, 11, 24),
    # date(2004, 11, 25),
    # date(2004, 11, 26),
    # date(2004, 11, 27),
    # date(2004, 11, 28),
    # date(2004, 11, 29),
    # date(2004, 11, 30),
    # date(2004, 12, 1),
]
n = len(x)

y1 = np.random.normal(size=n) + 2
y2 = np.random.normal(size=n) + 2
y3 = np.random.normal(size=n) + 2

# gs = GridSpec(1, 1)
fig, ax = plt.subplots(figsize=(9, 5))
# fig.autofmt_xdate()

ax.bar(x, y1, color='#55cc88')
ax.bar(x, y2, bottom=y1, color='#eeaa77')
ax.bar(x, y3, bottom=y1+y2, color='#cccccc')

ax.set_xlabel("Days")
ax.set_ylabel("Count")
ax.legend(["Opened", "Closed", "Cancelled", ])
ax.set_title("Daily issues and PRs for the past month")
# ax.fmt_xdata = mdates.DateFormatter('%b %d')
ax.set_xticks(x[0::5])
# 30: 6 skips
# 22: 5 skips
# 15: 3 skips
# 7: 2 skips
# 3: 2 skips

fig.savefig('dailyplot.png', dpi=200, bbox_inches = 'tight')

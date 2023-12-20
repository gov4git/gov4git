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
]
y1 = np.array([10, 20, 10, 30, 17, 1,2,3, 11,5, 4, 3])
y2 = np.array([20, 25, 15, 25, 12, 1,2,3,11,5, 4, 3])
y3 = np.array([12, 15, 19, 6, 11, 1,2,3,11,5, 4, 3])

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
ax.set_xticks(x[0::2])

fig.savefig('dailyplot.png', dpi=200, bbox_inches = 'tight')

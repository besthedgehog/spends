import io
from collections import defaultdict
from datetime import datetime

import matplotlib.pyplot as plt


def parse_data_by_day(data: list[dict]):
    totals = defaultdict(float)

    for item in data:
        date = datetime.fromisoformat(item["CreatedAt"]).date()
        totals[date] += item["Amount"]

    dates = sorted(totals.keys())
    labels = [d.strftime("%d.%m") for d in dates]
    values = [totals[d] for d in dates]

    return labels, values


def parse_cumulative_by_day(data: list[dict]):
    totals = defaultdict(float)

    for item in data:
        date = datetime.fromisoformat(item["CreatedAt"]).date()
        totals[date] += item["Amount"]

    dates_sorted = sorted(totals.keys())

    labels = [d.strftime("%d.%m") for d in dates_sorted]

    cumulative = []
    running_sum = 0.0
    for d in dates_sorted:
        running_sum += totals[d]
        cumulative.append(running_sum)

    return labels, cumulative


def parse_data_for_bars(data: list[dict]):
    totals = defaultdict(float)

    for item in data:
        totals[item["Priority"]] += item["Amount"]

    labels = [str(k) for k in totals.keys()]
    values = list(totals.values())
    return labels, values


def parse_data_for_category(data: list[dict]):
    totals = defaultdict(float)

    # for item in data:
    #     totals[item["Category"]] += item["Amount"]
    for item in data:
        category = item.get("Category")
        amount = item.get("Amount")

        if category is None or amount is None:
            continue  # пропускаем плохие данные

        totals[category] += float(amount)

    categories = list(totals.keys())
    values = list(totals.values())
    return categories, values


def parse_data_for_scatter(data: list[dict]):
    priorities = [item["Priority"] for item in data]
    amounts = [item["Amount"] for item in data]
    return priorities, amounts


def parse_data_for_pie(data):
    totals = defaultdict(float)

    # for item in data:
    #     totals[item["Priority"]] += item["Amount"]

    for item in data:
        priority = item.get("Priority")
        amount = item.get("Amount")

        if priority is None or amount is None:
            continue

        totals[priority] += amount

    labels = [str(cat) for cat in totals.keys()]
    values = list(totals.values())
    return labels, values


def fig_to_bytes(fig) -> io.BytesIO:
    """
    Из изображения matplotlib создаёт набор байтов
    """
    buf = io.BytesIO()
    fig.savefig(buf, format="png", dpi=150, bbox_inches="tight")
    plt.close(fig)  # Освобождаем память
    buf.seek(0)
    return buf

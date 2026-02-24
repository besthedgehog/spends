import matplotlib.pyplot as plt
import mplcyberpunk

from .utils import (
    parse_data_by_day,
    parse_data_for_bars,
    parse_data_for_category,
    parse_data_for_pie,
    parse_data_for_scatter,
)

plt.style.use("cyberpunk")


# --------------------
# Plots
# --------------------


def plot_pie(data: list[dict]):
    labels, values = parse_data_for_pie(data)
    plt.style.use("cyberpunk")

    fig, ax = plt.subplots(figsize=(7, 4))

    colors = ["C0", "C2", "#5555ff"]

    plt.pie(
        values,
        labels=labels,
        autopct="%1.1f%%",
        colors=colors[
            : len(values)
        ],  # важно обрезать если значений меньше,  # важно обрезать если значений меньше
        startangle=90,
        textprops={
            "color": "#FF0090",  # неоново-голубой
            "fontsize": 12,
            "weight": "bold",
        },
    )

    plt.title("Expenses by category")
    plt.tight_layout()
    return fig


def plot_by_day(data: list[dict]):
    labels, values = parse_data_by_day(data)

    fig, ax = plt.subplots(figsize=(7, 4))

    bars = ax.bar(
        labels,
        values,
        color="#00ffff",
        zorder=2,
    )

    ax.set_title("Expenses by day")
    ax.set_xlabel("Date")
    ax.set_ylabel("Amount")

    mplcyberpunk.add_bar_gradient(bars=bars)
    mplcyberpunk.add_glow_effects()

    return fig


# По приорететам
def plot_bars(data: list[dict]):
    """
    По приорететам
    """

    labels, values = parse_data_for_bars(data)

    plt.style.use("cyberpunk")
    fig, ax = plt.subplots(figsize=(6, 4))

    colors = ["C0", "C2", "C1"]

    bars = ax.bar(labels, values, color=colors, zorder=2)

    ax.set_title("Expenses by priority")

    mplcyberpunk.add_bar_gradient(bars=bars)
    mplcyberpunk.add_glow_effects()

    plt.tight_layout()
    return fig


# Столбики по категориям
def plot_by_category(data: list[dict]):
    """
    Столбики по категориям
    """
    categories, values = parse_data_for_category(data)
    plt.style.use("cyberpunk")
    fig, ax = plt.subplots(figsize=(7, 4))

    colors = [f"C{i}" for i in range(len(categories))]

    bars = ax.bar(
        categories,
        values,
        color=colors,
        zorder=2,
    )

    ax.set_title("Expenses by category")
    ax.set_ylabel("Amount")

    mplcyberpunk.add_bar_gradient(bars=bars)
    mplcyberpunk.add_glow_effects()

    plt.tight_layout()

    return fig


def plot_cumulative(data: list[dict]):
    labels, values = parse_data_by_day(data)

    cumulative = []
    total = 0.0
    for v in values:
        total += v
        cumulative.append(total)

    fig, ax = plt.subplots(figsize=(7, 4))

    ax.plot(
        labels,
        cumulative,
        marker="o",
        linewidth=2,
    )

    ax.set_title("Cumulative expenses")
    ax.set_xlabel("Date")
    ax.set_ylabel("Total")

    mplcyberpunk.add_glow_effects()

    return fig


def plot_scatter(data: list[dict]):
    priorities, mounts = parse_data_for_scatter(data)
    priorities = [item["Priority"] for item in data]
    amounts = [item["Amount"] for item in data]

    fig, ax = plt.subplots(figsize=(6, 4))

    ax.scatter(
        priorities,
        amounts,
        s=70,
        color="#00ffff",
        zorder=3,
    )

    ax.set_title("Priority vs Amount")
    ax.set_xlabel("Priority")
    ax.set_ylabel("Amount")

    mplcyberpunk.make_scatter_glow()

    return fig

from flask import jsonify, request, send_file

from .plots import (
    plot_bars,
    plot_by_category,
    plot_by_day,
    plot_cumulative,
    plot_pie,
    plot_scatter,
)
from .utils import fig_to_bytes


def register_routes(app):
    # Столбики по приорететам
    @app.post("/plot/by-day")
    def plot_by_day_route():
        data = request.json
        fig = plot_by_day(data)
        return send_file(fig_to_bytes(fig), mimetype="image/png")

    # Столбики по приорететам (всего три столбика)
    @app.post("/plot/priority")
    def plot_bars_route():
        data = request.json
        fig = plot_bars(data)
        return send_file(fig_to_bytes(fig), mimetype="image/png")

    # Столбики по категориям
    @app.post("/plot/category")
    def plot_by_category_route():
        data = request.json
        fig = plot_by_category(data)
        return send_file(fig_to_bytes(fig), mimetype="image/png")

    # График с накоплением
    @app.post("/plot/cumulative")
    def plot_cumulative_route():
        data = request.json
        fig = plot_cumulative(data)
        return send_file(fig_to_bytes(fig), mimetype="image/png")

    # Точечная диаграмма
    @app.post("/plot/scatter")
    def plot_scatter_route():
        data = request.json
        fig = plot_scatter(data)
        return send_file(fig_to_bytes(fig), mimetype="image/png")

    # Пирог по категориям
    @app.post("/plot/pie")
    def plot_pie_route():
        data = request.json
        fig = plot_pie(data)
        return send_file(fig_to_bytes(fig), mimetype="image/png")

    @app.get("/health")
    def health():
        return jsonify({"status": "ok"})

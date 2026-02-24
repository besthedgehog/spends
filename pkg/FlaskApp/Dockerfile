FROM python:3.13-slim

WORKDIR /app

# Копируем зависимости
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Копируем код
COPY . .

# Запускаем через Gunicorn (PRODUCTION!)
CMD ["gunicorn", "-w", "4", "-b", "0.0.0.0:5005", "main:app"]

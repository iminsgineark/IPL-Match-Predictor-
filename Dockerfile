FROM python:3.9

RUN apt-get update && apt-get install -y python3-venv

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY . .

EXPOSE 8501

CMD ["streamlit", "run", "app.py", "--server.port=8501", "--server.address=127.0.0.1"]



#replace 0.0.0.0 with localhost in your browser
FROM python:3.9

LABEL maintainer="Utkrist Ark"
LABEL description="a dockerfile for ML model"

RUN apt-get update && \
    apt-get install -y --no-install-recommends curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* 
#&& apt-get install -y python3-venv

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt && \
    groupadd -g 1001 arkgroup && \
    useradd -u 1001 -g arkgroup -m ark && \
    chown -R ark:arkgroup /home/ark

USER ark

WORKDIR /app

COPY . .

EXPOSE 8501

CMD ["streamlit", "run", "app.py", "--server.port=8501", "--server.address=0.0.0.0"]

HEALTHCHECK CMD curl --fail http://localhost:8501/ || exit 1

#replace 0.0.0.0 with localhost in your browser

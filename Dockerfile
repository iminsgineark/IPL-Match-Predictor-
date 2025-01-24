FROM python:3.9-slim

LABEL maintainer="Utkrist Ark"
LABEL description="a dockerfile for ML model"

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends curl && \
    apt-get clean && \ 
    rm -rf /var/lib/apt/lists/* 

COPY requirements.txt /app/

RUN pip install --no-cache-dir -r requirements.txt && \
    pip install streamlit && \
    groupadd -g 10001 arkgroup && \
    useradd -u 10001 -g arkgroup -m ark && \
    chown -R ark:arkgroup /home/ark

USER ark

WORKDIR /app

COPY app.py /app/
COPY pipe.pkl /app/

EXPOSE 8501

CMD ["streamlit", "run", "app.py", "--server.port=8501", "--server.address=0.0.0.0", "--server.baseUrlPath=/model"]

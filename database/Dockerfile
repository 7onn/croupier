FROM postgres:11-alpine

RUN cp /usr/share/zoneinfo/UTC /etc/localtime && \
	echo UTC > /etc/timezone


WORKDIR /app
ENV WORKDIR /app

COPY . .

COPY ./scripts/docker-entrypoint-initdb.sh /docker-entrypoint-initdb.d

EXPOSE 5432

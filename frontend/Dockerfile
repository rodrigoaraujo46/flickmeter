FROM oven/bun:latest

RUN useradd -ms /bin/sh -u 1001 app
USER app

WORKDIR /app
COPY ./package.json ./bun.lock ./
RUN bun install

COPY --chown=app:app . /app

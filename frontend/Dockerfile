FROM node:latest
RUN npm install -g pnpm

RUN useradd -ms /bin/sh -u 1001 app
USER app

WORKDIR /app
COPY ./package.json ./pnpm-lock.yaml ./
RUN pnpm install

COPY --chown=app:app . /app

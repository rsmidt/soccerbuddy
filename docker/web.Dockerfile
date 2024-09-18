FROM node:22-slim AS base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

FROM base AS prod

RUN mkdir /app
COPY pnpm-lock.yaml /app
COPY package.json /app
WORKDIR /app
RUN pnpm install

COPY . .
RUN pnpm run build

FROM base
COPY --from=prod /app/node_modules node_modules
COPY --from=prod /app/build build
COPY --from=prod /app/package.json package.json

EXPOSE 3000/tcp
ENTRYPOINT [ "node", "build" ]

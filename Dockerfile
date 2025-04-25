FROM node:23-alpine AS frontend-deps

ARG SKIP_FRONTEND_TEST

WORKDIR /srv/frontend/

COPY ./frontend/package.json ./frontend/pnpm-lock.yaml ./frontend/pnpm-workspace.yaml /srv/frontend/
COPY ./frontend/app/package.json /srv/frontend/app/

RUN apk add --no-cache --update git && \
    npm i -g pnpm@10.8.0;

RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile


FROM frontend-deps AS build-frontend

WORKDIR /srv/frontend/app/

COPY ./frontend/app/ /srv/frontend/app/

RUN \
  if [ -z "$SKIP_FRONTEND_TEST" ]; then \
    pnpm checks-all; \
  else \
    echo 'Skip frontend test'; \
  fi

RUN pnpm build


FROM umputun/baseimage:buildgo-v1.15.0 AS build-backend

ARG SKIP_BACKEND_TEST

RUN apk --no-cache add gcc libc-dev

ADD ./backend /build/backend

COPY --from=build-frontend /srv/frontend/app/dist/ /build/backend/app/cmd/web/

WORKDIR /build/backend

RUN echo go version: `go version`

# run tests
RUN \
  	if [ -z "$SKIP_BACKEND_TEST" ] ; then \
				CGO_ENABLED=1 go test -race -p 1 -timeout="300s" -covermode=atomic -coverprofile=/profile.cov_tmp ./... && \
				cat /profile.cov_tmp | grep -v "_mock.go" > /profile.cov && \
				golangci-lint run --config .golangci.yml ./... ; \
  	else \
    		echo 'Skip backend test'; \
  	fi

RUN \
    version="$(/script/version.sh)" && \
    echo "version=$version" && \
    go build -o neskuchka -ldflags "-X main.revision=${version} -s -w" ./app


FROM umputun/baseimage:app-v1.15.0

WORKDIR /srv

COPY --from=build-backend /build/backend/neskuchka /srv/neskuchka

RUN chown -R app:app /srv
RUN ln -s /srv/neskuchka /usr/bin/neskuchka

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD curl --fail http://localhost:8080/ping || exit 1

CMD ["/srv/neskuchka", "server"]
package buildscript

const script = `#!/bin/bash
Name="{{ .Name }}"
MainPath="{{ .Main }}"
Org="{{ .Org }}"

# shellcheck disable=SC2046
Version=$(git describe --tags $(git rev-list --tags --max-count=1))
# shellcheck disable=SC2154
GitCommit=$(git log --pretty=format:"%h" -1)
BuildTime=$(date +%FT%T%z)
# shellcheck disable=SC2154
Compiler=${go version}

build_application(){
  git checkout "${Version}"
  docker build -t "${Org}/${Name}:${Version}" \
  --build-arg NAME="${Name}" \
  --build-arg VERSION="${Version}" \
  --build-arg BUILD_TIME="${BuildTime}" \
  --build-arg COMMIT="${GitCommit}" \
  --build-arg COMPILER="${Compiler}" \
  --build-arg MAIN_PATH="${MainPath}" .
}

print_app_info(){
  echo "****************************************"
  echo "App:${Org}:${Name}"
  echo "Version:${Version}"
  echo "Commit:${GitCommit}"
  echo "Build:${BuildTime}"
  echo "Compiler:${Compiler}"
  echo "Main_Path:${MainPath}"
  echo "****************************************"
  echo ""
}

print_app_info
build_application
`

const dockerFile = `
FROM node:18.4.0 as ui
ARG NAME
ARG VERSION
WORKDIR /ui_build
ADD ui .
RUN npm install && npm run build

FROM golang:1.18 as build
ARG NAME
ARG VERSION
ARG COMMIT
ARG BUILD_TIME
ARG MAIN_PATH
ARG COMPILER
ENV GOPROXY=https://goproxy.cn,direct
ENV LDFLAGS=" \
    -X 'github.com/lishimeng/app-starter/version.AppName=${NAME}' \
    -X 'github.com/lishimeng/app-starter/version.Version=${VERSION}' \
    -X 'github.com/lishimeng/app-starter/version.Commit=${COMMIT}' \
    -X 'github.com/lishimeng/app-starter/version.Build=${BUILD_TIME}' \
    -X 'github.com/lishimeng/app-starter/version.Compiler=${COMPILER}' \
    "
WORKDIR /release
ADD . .
COPY --from=ui /ui_build/dist/ static/
RUN go mod download && go mod verify
RUN go build -v --ldflags "${LDFLAGS}" -o ${NAME} ${MAIN_PATH}

FROM ubuntu:22.04 as prod
ARG NAME
EXPOSE 80/tcp
WORKDIR /
COPY --from=build /release/${NAME} /
RUN ln -s /${NAME} /app
CMD [ "/app"]
`

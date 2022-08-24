package app

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

const (
	DockerfileName = "Dockerfile"
	BuildFileName  = "build.sh"
)

const (
	buildTpl = `
#!/bin/bash
NAME="{{ .name }}"
MAIN_PATH="{{ .mainPath }}"
ORG="{{ .dockerOrganization }}"

# shellcheck disable=SC2046
VERSION=$(git describe --tags $(git rev-list --tags --max-count=1))
# shellcheck disable=SC2154
COMMIT=$(git log --pretty=format:"%h" -1)
BUILD_TIME=$(date +%FT%T%z)

build_application(){
  git checkout "${VERSION}"
  docker build -t "${ORG}/${NAME}:${VERSION}" \
  --build-arg NAME="${NAME}" \
  --build-arg VERSION="${VERSION}" \
  --build-arg COMMIT="${COMMIT}" \
  --build-arg BUILD_TIME="${BUILD_TIME}" \
  --build-arg MAIN_PATH="${MAIN_PATH}" .
}

print_app_info(){
  echo "****************************************"
  echo "App:${NAME}"
  echo "Version:${VERSION}"
  echo "Commit:${COMMIT}"
  echo "BuildTime:${BUILD_TIME}"
  echo "Main_Path:${MAIN_PATH}"
  echo "****************************************"
  echo ""
}

print_app_info
build_application
`
	dockerTpl = `
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
ARG TAG
ARG BUILD_TIME
ARG MAIN_PATH
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /release
ADD . .
COPY --from=ui /ui_build/dist/ static/
RUN go mod download && go mod verify
RUN go build -v --ldflags "-X cmd.AppName=${NAME} -X cmd.Version=${VERSION} -X cmd.Commit=${COMMIT} -X cmd.Build=${BUILD_TIME}" -o ${NAME} ${MAIN_PATH}

FROM ubuntu:22.04 as prod
ARG NAME
EXPOSE 80/tcp
WORKDIR /
COPY --from=build /release/${NAME} /
RUN ln -s /${NAME} /app
CMD [ "/app"]
`
)

type ApplicationConfig struct {
	Name               string `json:"name,omitempty"`
	MainPath           string `json:"mainPath,omitempty"`
	DockerOrganization string `json:"dockerOrganization,omitempty"`
}

func GenerateDockerfile(config ApplicationConfig) (err error) {

	bs, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	var param map[string]string
	err = json.Unmarshal(bs, &param)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = genBuildFiles(BuildFileName, buildTpl, param)
	if err != nil {
		fmt.Println("gen build script failure")
		fmt.Println(err)
		return
	}

	err = genBuildFiles(DockerfileName, dockerTpl, param)
	if err != nil {
		fmt.Println("gen docker script failure")
		fmt.Println(err)
		return
	}

	return
}

func genBuildFiles(output string, tplContent string, data interface{}) (err error) {
	tpl, err := template.New("_").Parse(tplContent)
	if err != nil {
		return
	}
	fw, err := os.Create(output)
	if err != nil {
		return
	}
	defer func() {
		_ = fw.Close()
	}()
	err = tpl.Execute(fw, data)
	if err != nil {
		return
	}
	return
}

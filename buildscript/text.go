package buildscript

const localBuildScript = `#!/bin/bash

Version=''
GitCommit=''
BuildTime=''

Base="github.com/lishimeng/app-starter/version"
GOPROXY=https://goproxy.cn,direct
if [ -n "${GOOS}" ]; then
  echo "custom GOOS"
else
  GOOS=linux
fi
echo "GOOS:${GOOS}"

if [ -n "${GOARCH}" ]; then
  echo "custom GOARCH"
else
  GOARCH=amd64
fi
echo "GOARCH:${GOARCH}"

meta(){
  # shellcheck disable=SC2046
  Version=$(git describe --tags $(git rev-list --tags --max-count=1))
  # shellcheck disable=SC2154
  GitCommit=$(git log --pretty=format:"%h" -1)
  BuildTime=$(date +%FT%T%z)
}

meta_dev(){
  # shellcheck disable=SC2046
  Version="snapshot_"$(git log --pretty=format:"%ct" -1)
  # shellcheck disable=SC2154
  GitCommit=$(git log --pretty=format:"%h" -1)
  BuildTime=$(date +%FT%T%z)
}

checkout_tag(){
  git checkout "${Version}"
}

common(){
  echo ""
}

print_app_info(){
  local Name=$1
  local AppPath=$2
  echo "****************************************"
  echo "App:${Name}"
  echo "Version:${Version}"
  echo "Commit:${GitCommit}"
  echo "Build:${BuildTime}"
  echo "Main_Path:${AppPath}"
  echo "Build_Os:${GOOS}"
  echo "Build_Arch:${GOARCH}"
  echo "GoProxy:${GOPROXY}"
  echo "****************************************"
  echo ""
}

build(){
  local Name=$1
  local AppPath=$2
  local LdFlags=" \
-X ${Base}.AppName=${Name} \
-X ${Base}.Version=${Version} \
-X ${Base}.Commit=${GitCommit} \
-X ${Base}.Build=${BuildTime} \
"
  go mod tidy && go mod vendor
  GOOS="${GOOS}" GOARCH="${GOARCH}" go build -v --ldflags "${LdFlags} -X ${Base}.Compiler=$(go version | sed 's/[ ][ ]*/_/g')" -o "${Name}" "${AppPath}"/main.go
}

build_app(){
  local Name=$1
  local AppPath=$2
  if [ "${GOOS}" == "windows" ]; then
      Name="${Name}.exe"
  fi
  print_app_info "${Name}" "${AppPath}"
  build "${Name}" "${AppPath}"
}

build_all(){
  common
  {{- range $_, $item := .Applications }}
  build_app '{{ $item.Name }}' '{{ $item.AppPath }}'
  {{- end }}
}

# command
case  $1 in
    release)
    meta
    checkout_tag
		build_all
        ;;
    *)
		meta_dev
		build_all
        ;;
esac`

const script = `#!/bin/bash
Namespace="{{ .Pro.Namespace }}"

# shellcheck disable=SC2046
Version=$(git describe --tags $(git rev-list --tags --max-count=1))
# shellcheck disable=SC2154
GitCommit=$(git log --pretty=format:"%h" -1)
BuildTime=$(date +%FT%T%z)

checkout_tag(){
  git checkout "${Version}"
}

common(){
  echo ""
}

build_image(){
  local Name=$1
  local AppPath=$2
  print_app_info "${Name}" "${AppPath}"

  docker build -t "${Namespace}/${Name}:${Version}" \
  --build-arg NAME="${Name}" \
  --build-arg VERSION="${Version}" \
  --build-arg BUILD_TIME="${BuildTime}" \
  --build-arg COMMIT="${GitCommit}" \
  --build-arg APP_PATH="${AppPath}" -f "./${AppPath}/Dockerfile" .
}

print_app_info(){
  local Name=$1
  local AppPath=$2
  echo "****************************************"
  echo "App:${Name}[${Namespace}]"
  echo "Version:${Version}"
  echo "Commit:${GitCommit}"
  echo "Build:${BuildTime}"
  echo "Main_Path:${AppPath}"
  echo "****************************************"
  echo ""
}

push_image(){
  local Name=$1
  echo "****************************************"
  echo "Push:${Namespace}:${Name}:${Version}"
  echo "****************************************"
  echo ""
  docker tag  "${Namespace}/${Name}:${Version}" "${Namespace}/${Name}"
  docker push "${Namespace}/${Name}:${Version}"
  docker push "${Namespace}/${Name}"
}

build_all(){
  common
  checkout_tag
  {{- range $_, $item := .Applications }}
  build_image '{{ $item.Name }}' '{{ $item.AppPath }}'
  {{- end }}
}

push_all(){
  common
  {{- range $_, $item := .Applications }}
  push_image '{{ $item.Name }}'
  {{- end }}
}

case  $1 in
    push)
		push_all
        ;;
    *)
		build_all
        ;;
esac

`

const dockerFile = `{{- if .App.HasUI }}
FROM {{ .BuildImageVersion.Node }} as ui
ARG NAME
ARG VERSION
ARG APP_PATH
WORKDIR /ui_build
ADD ${APP_PATH}/ui .
RUN npm i pnpm -g && pnpm install && pnpm run build

{{- end }}
FROM {{ .BuildImageVersion.Golang }} as build
ARG NAME
ARG VERSION
ARG COMMIT
ARG BUILD_TIME
ARG APP_PATH
ARG BASE="github.com/lishimeng/app-starter/version"
ENV GOPROXY=https://goproxy.cn,direct
ARG LDFLAGS=" \
-X ${BASE}.AppName=${NAME} \
-X ${BASE}.Version=${VERSION} \
-X ${BASE}.Commit=${COMMIT} \
-X ${BASE}.Build=${BUILD_TIME} \
"
WORKDIR /release
ADD . .
{{- if .App.HasUI }}
COPY --from=ui /ui_build/dist/ ${APP_PATH}/static/
{{- end }}

RUN go mod download && go mod verify
RUN go build -v --ldflags "${LDFLAGS} -X ${BASE}.Compiler=$(go version | sed 's/[ ][ ]*/_/g')" -o ${NAME} ${APP_PATH}/main.go

FROM {{ .BuildImageVersion.Runtime }} as prod
ARG NAME
EXPOSE 80/tcp
WORKDIR /
COPY --from=build /release/${NAME} /
RUN ln -s /${NAME} /app
CMD [ "/app"]
`

// 基础镜像, 设置了+8时区
//
// docker build -t {namespace}/alpine:{version} .

const baseDockerfileAlpine = `FROM alpine:3.17
MAINTAINER lishimeng
ENV TZ=Asia/Shanghai

RUN apk update \
    && apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && mkdir /lib64 \
    && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
`

const baseDockerfileUbuntu = `FROM ubuntu
MAINTAINER lishimeng
ENV TIME_ZONE Asia/Shanghai
 
RUN sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y tzdata \
    && ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone \
    && dpkg-reconfigure -f noninteractive tzdata \
    && apt-get clean \
    && rm -rf /tmp/* /var/cache/* /usr/share/doc/* /usr/share/man/* /var/lib/apt/lists/*
`

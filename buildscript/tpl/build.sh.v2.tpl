#!/bin/bash
Namespace="{{ .Pro.Namespace }}"

Base="github.com/lishimeng/app-starter/version"
GOPROXY=${GOPROXY:-https://goproxy.cn,direct}
# 开关 -m / -g 打开时使用下列内置地址（与 Dockerfile 默认 npmmirror / 国内 Go 代理一致）
NPM_MIRROR_REGISTRY=https://registry.npmmirror.com
GOPROXY_CN=https://goproxy.cn,direct

USE_NPM_MIRROR=""
USE_GOPROXY_CN=""
CMD_MODE=""

Version=''
GitCommit=''
BuildTime=''
BuildMode="release"

CmdPrint=0
CmdRun=1
CmdMode=${CmdRun}

LogLvlErr=4
LogLvlWarn=3
LogLvlInfo=2
LogLvlDebug=1
LogLevel=${LogLvlInfo}

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # no color

while getopts ":hl:d:mg" opt; do
  case $opt in
    l)
      #echo "set opt l ${OPTARG}"
      LOG_LEVEL="${OPTARG}"
      ;;
    d)
      #echo "set opt d ${OPTARG}"
      CMD_MODE="${OPTARG}"
      ;;
    m)
      USE_NPM_MIRROR=1
      ;;
    g)
      USE_GOPROXY_CN=1
      ;;
    h)
      echo "Usage: build.sh [-hl:dmg] {mode:optional} {app_name:optional}"
      echo ""
      echo "parameters:"
      echo "mode: dev/release/push"
      echo "app_name: if present, only build target app"
      echo ""
      echo "options:"
      echo "-d debug/run"
      echo "-l log level, values:error/warn/info/debug"
      echo "-m use builtin NPM mirror (set NPM_REGISTRY to NPM_MIRROR_REGISTRY)"
      echo "-g use builtin Go proxy CN (set GOPROXY to GOPROXY_CN)"
      echo ""
      echo "environment (optional):"
      echo "  NPM_REGISTRY   镜像构建时传给 Docker；未设置则用各 Dockerfile 默认；-m 会覆盖为本仓库 NPM_MIRROR_REGISTRY"
      echo "  GOPROXY        镜像构建时传给 Docker；-g 会覆盖为本仓库 GOPROXY_CN"
      echo ""
      exit 0
      ;;
    \?)
      echo "错误：无效选项 -$OPTARG" >&2
      exit 1
      ;;
    :)
      echo "错误：选项 -$OPTARG 需要参数" >&2
      exit 1
      ;;
  esac
done
shift $((OPTIND - 1))

handle_user_env(){
  if [ -n "${LOG_LEVEL}" ]; then
    log_info "set cmd level:${LOG_LEVEL}"
    case ${LOG_LEVEL} in
    error)
      LogLevel=${LogLvlErr}
      ;;
    warn)
      LogLevel=${LogLvlWarn}
      ;;
    info)
      LogLevel=${LogLvlInfo}
      ;;
    debug)
      LogLevel=${LogLvlDebug}
      ;;
    esac
  fi

  if [ -n "${CMD_MODE}" ]; then
      log_info "set cmd mode:${CMD_MODE}"
      case ${CMD_MODE} in
      debug)
        CmdMode=${CmdPrint}
        ;;
      run)
        CmdMode=${CmdRun}
        ;;
      esac
  fi

  if [ -n "${USE_NPM_MIRROR}" ]; then
    NPM_REGISTRY="${NPM_MIRROR_REGISTRY}"
    log_info "use npm mirror: ${NPM_REGISTRY}"
  fi
  if [ -n "${USE_GOPROXY_CN}" ]; then
    GOPROXY="${GOPROXY_CN}"
    log_info "use goproxy CN: ${GOPROXY}"
  fi
}

log_error() {
  if [ ${LogLevel} -le ${LogLvlErr} ]; then
      echo -e "${RED}[ERROR] $1${NC}"
    fi
}

log_warn() {
  if [ ${LogLevel} -le ${LogLvlWarn} ]; then
    echo -e "${YELLOW}[WARN] $1${NC}"
  fi
}

log_info() {
  #echo "log_info" "log_level:" "${LogLevel}" "me:" "${LogLvlInfo}"
  if [ ${LogLevel} -le ${LogLvlInfo} ]; then
    echo -e "${GREEN}[INFO] $1${NC}"
  fi
}

log_debug() {
  if [ ${LogLevel} -le ${LogLvlDebug} ]; then
    echo -e "${BLUE}[DEBUG] $1${NC}"
  fi
}

handle_env(){
  if [ -n "${GOOS}" ]; then
    log_debug "custom GOOS"
  else
    GOOS=linux
  fi
  log_debug "GOOS:${GOOS}"

  if [ -n "${GOARCH}" ]; then
    log_debug "custom GOARCH"
  else
    GOARCH=amd64
  fi
  log_debug "GOARCH:${GOARCH}"
}

checkout_tag(){
  log_debug "checkout_tag..."
  git checkout "${Version}"
}

common(){
  log_debug "common..."
}

meta_release(){
  log_debug "meta_release..."
  # shellcheck disable=SC2046
  Version=$(git describe --tags $(git rev-list --tags --max-count=1))
  # shellcheck disable=SC2154
  GitCommit=$(git log --pretty=format:"%h" -1)
  BuildTime=$(date +%FT%T%z)
}

meta_dev(){
  log_debug "meta_dev..."
  # shellcheck disable=SC2046
  Version="snapshot_"$(git log --pretty=format:"%ct" -1)
  # shellcheck disable=SC2154
  GitCommit=$(git log --pretty=format:"%h" -1)
  BuildTime=$(date +%FT%T%z)
}

print_app_info(){
  local Name=$1
  local AppPath=$2
  log_info "****************************************"
  log_info "App:${Name}[${Namespace}]"
  log_info "Version:${Version}"
  log_info "Commit:${GitCommit}"
  log_info "Build:${BuildTime}"
  log_info "Main_Path:${AppPath}"
  log_info "GoProxy:${GOPROXY}"
  if [ "${BuildMode}" == "dev" ]; then
    log_info "Build_Os:${GOOS}"
    log_info "Build_Arch:${GOARCH}"
  fi
  if [ -n "${NPM_REGISTRY:-}" ]; then
    log_info "NpmRegistry:${NPM_REGISTRY}"
  fi
  log_info "****************************************"
}

push_image(){
  local Name=$1
  log_info "****************************************"
  log_info "Push:${Namespace}:${Name}:${Version}"
  log_info "****************************************"
  if [ ${CmdMode} == 0 ]; then
    exit 0
  else
    docker tag  "${Namespace}/${Name}:${Version}" "${Namespace}/${Name}"
    docker push "${Namespace}/${Name}:${Version}"
    docker push "${Namespace}/${Name}"
  fi
}

build_image(){
  local Name=$1
  local AppPath=$2
  print_app_info "${Name}" "${AppPath}"
  if [ ${CmdMode} == 0 ]; then
    exit 0
  else
    local docker_extra=()
    docker_extra+=(--build-arg GOPROXY="${GOPROXY}")
    if [ -n "${NPM_REGISTRY:-}" ]; then
      docker_extra+=(--build-arg NPM_REGISTRY="${NPM_REGISTRY}")
    fi
    log_info "docker build start: ${Name}"
    if ! docker build -t "${Namespace}/${Name}:${Version}" \
      --build-arg NAME="${Name}" \
      --build-arg VERSION="${Version}" \
      --build-arg BUILD_TIME="${BuildTime}" \
      --build-arg COMMIT="${GitCommit}" \
      --build-arg APP_PATH="${AppPath}" \
      "${docker_extra[@]}" \
      -f "./${AppPath}/Dockerfile" .; then
      log_error "docker build failed: ${Name}"
      exit 1
    fi
    log_info "docker build ok: ${Name}"
  fi
}

build_local(){
  local LdFlags=""
  local Name=$1
  local AppPath=$2
  log_info "build local [$1][${BuildMode}]..."
  log_info "$1" "$2"
  if [ "${GOOS}" == "windows" ]; then
      Name="${Name}.exe"
  fi
  print_app_info "${Name}" "${AppPath}"
  LdFlags=" \
  -X ${Base}.AppName=${Name} \
  -X ${Base}.Version=${Version} \
  -X ${Base}.Commit=${GitCommit} \
  -X ${Base}.Build=${BuildTime} \
  "
  if [ ${CmdMode} == 0 ]; then
    exit 0
  else
    go mod tidy && go mod vendor
    GOOS="${GOOS}" GOARCH="${GOARCH}" go build -v --ldflags "${LdFlags} -X ${Base}.Compiler=$(go version | sed 's/[ ][ ]*/_/g')" -o "${Name}" "${AppPath}"/main.go
  fi
}

build_application(){
  case ${BuildMode} in
  dev)
    build_local "$1" "$2"
    ;;
  release)
    build_image "$1" "$2"
    ;;
  esac
}

main_cmd(){
  log_debug "main_cmd..."
  case  ${BuildMode} in
      push)
        log_debug "handle push..."
        meta_release
  		push_all "$@"
        ;;
      dev)
        log_debug "handle dev..."
        meta_dev
        build_all "$@"
        ;;
      release)
        log_debug "handle release..."
        meta_release
        checkout_tag
        build_all "$@"
        ;;
      *)
        log_debug "no cmd, exit"
        exit 0
        ;;
  esac
}

main(){
  handle_user_env
  handle_env
  log_debug "main..."
  log_debug "handle build mode..."
  if [ $# -ge 1 ]; then
    case $1 in
    release)
      BuildMode="release"
      shift
      ;;
    dev)
      BuildMode="dev"
      shift
      ;;
    push)
      BuildMode="push"
      shift
      ;;
    *)
      log_debug "use default build_mode"
      ;;
    esac
  fi
  log_info "build mode [${BuildMode}]..."
  sleep 0.2
  main_cmd "$@"
}

build_special(){
  local Name=$1
  local AppPath=""
  log_debug "build_special..."
  common
  case $1 in
{{- range $_, $item := .Applications }}
      {{ $item.Name }})
      log_debug "{{ $item.Name }} -> [{{ $item.AppPath }}]"
      AppPath='{{ $item.AppPath }}'
      ;;
{{- end }}
      *)
      exit 0
      ;;
  esac
  build_application "${Name}" "${AppPath}"
}

build_all(){
  log_debug "build_all..."
  #common
  if [ $# -eq 1 ]; then
    log_info "build [$1]"
    build_special "$1"
  else
{{- range $_, $item := .Applications }}
  	build_application '{{ $item.Name }}' '{{ $item.AppPath }}'
{{- end }}
  fi
}

push_all(){
  log_debug "push_all..."
  common
{{- range $_, $item := .Applications }}
  push_image '{{ $item.Name }}'
{{- end }}
}

main "$@"

{{- if .App.HasUI }}
FROM {{ .BuildImageVersion.Node }} as ui
ARG NAME
ARG VERSION
ARG APP_PATH
ARG NPM_REGISTRY=https://registry.npmmirror.com
ENV NPM_CONFIG_REGISTRY=${NPM_REGISTRY}
WORKDIR /ui_build
ADD ${APP_PATH}/ui .
RUN npm i pnpm -g && pnpm install --dangerously-allow-all-builds && pnpm run build

{{- end }}
FROM {{ .BuildImageVersion.Golang }} as build
ARG NAME
ARG VERSION
ARG COMMIT
ARG BUILD_TIME
ARG APP_PATH
ARG BASE="github.com/lishimeng/app-starter/version"
ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}
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

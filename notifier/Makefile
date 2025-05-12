BINDIR=${CURDIR}/bin
PACKAGE=route256/notifier
LOCAL_CONFIG=${CURDIR}/configs/values_local.yaml
BUILD_COMMAND=go build -o ${BINDIR}/notifier cmd/notifier/server.go

bindir:
	mkdir -p ${BINDIR}

build: bindir
	${BUILD_COMMAND}

run: build
	export CONFIG_FILE=${LOCAL_CONFIG} && \
	${BINDIR}/notifier

watch:
	export CONFIG_FILE=${LOCAL_CONFIG} && air --build.cmd "${BUILD_COMMAND}" \
		--build.bin ${BINDIR}/notifier --build.send_interrupt true --build.kill_delay 500

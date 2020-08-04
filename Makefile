export GOSUMDB=off
export GOPROXY=https://goproxy.io

BUILD_VERSION   := $(shell git describe --tags)
GIT_COMMIT_SHA1 := $(shell git rev-parse --short HEAD)
BUILD_TIME      := $(shell date "+%F %T")

prebuild:
	go get golang.org/x/mobile/cmd/gomobile

build-fperf-android:
	rm -rf output/android
	mkdir -p output/android
	gomobile bind -target android/arm64,android/arm -o output/android/fperf.aar github.com/fregie/gotool/fperf
	cd output && zip -r fperf_android_${BUILD_VERSION}_${GIT_COMMIT_SHA1}.zip android
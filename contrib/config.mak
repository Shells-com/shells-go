ifeq ($(TARGET_GOOS),windows)
GOLDFLAGS+=-H=windowsgui
endif
DIST_ARCHS=linux_amd64 linux_arm64 linux_arm

NAME=ssr2clashr
BINDIR=bin
VERSION=$(shell date +"v%Y.%m.%d")
GOBUILD=CGO_ENABLED=0 go build -ldflags '-X github.com/heiha/ssr2clashr/cmd.VERSION=$(VERSION) -w -s'
basepath=$(shell pwd)

ALL_PLATFORM_LIST = \
	darwin-amd64 \
	linux-386 \
	linux-amd64 \
	linux-armv5 \
	linux-armv6 \
	linux-armv7 \
	linux-armv8 \
	linux-mips-softfloat \
	linux-mips-hardfloat \
	linux-mipsle \
	linux-mips64 \
	linux-mips64le \
	freebsd-386 \
	freebsd-amd64

ALL_WINDOWS_ARCH_LIST = \
	windows-386 \
	windows-amd64


PLATFORM_LIST = \
	darwin-amd64 \
	linux-amd64 \
	linux-armv8 \

WINDOWS_ARCH_LIST = \
	windows-amd64

.PHONY: FORCE

all: linux-amd64

darwin-amd64: FORCE
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-386: FORCE
	GOARCH=386 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-amd64: FORCE
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-armv5: FORCE
	GOARCH=arm GOOS=linux GOARM=5 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-armv6: FORCE
	GOARCH=arm GOOS=linux GOARM=6 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-armv7: FORCE
	GOARCH=arm GOOS=linux GOARM=7 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-armv8: FORCE
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-mips-softfloat: FORCE
	GOARCH=mips GOMIPS=softfloat GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-mips-hardfloat: FORCE
	GOARCH=mips GOMIPS=hardfloat GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-mipsle: FORCE
	GOARCH=mipsle GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-mips64: FORCE
	GOARCH=mips64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-mips64le: FORCE
	GOARCH=mips64le GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

freebsd-386: FORCE
	GOARCH=386 GOOS=freebsd $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

freebsd-amd64: FORCE
	GOARCH=amd64 GOOS=freebsd $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

windows-386: FORCE
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

windows-amd64: FORCE
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe


gz_all_releases=$(addsuffix .gz, $(ALL_PLATFORM_LIST))
zip_all_releases=$(addsuffix .zip, $(ALL_WINDOWS_ARCH_LIST))

gz_releases=$(addsuffix .gz, $(PLATFORM_LIST))
zip_releases=$(addsuffix .zip, $(WINDOWS_ARCH_LIST))

%.gz : %
	chmod +x $(BINDIR)/$(NAME)-$(basename $@)
	upx --best -fv $(BINDIR)/$(NAME)-$(basename $@)
	gzip -f -S .gz $(BINDIR)/$(NAME)-$(basename $@)

%.zip : %
	upx --best -fv $(BINDIR)/$(NAME)-$(basename $@).exe
	zip -m -j $(BINDIR)/$(NAME)-$(basename $@).zip $(BINDIR)/$(NAME)-$(basename $@).exe


arch: $(PLATFORM_LIST) $(WINDOWS_ARCH_LIST)

all-arch: $(ALL_PLATFORM_LIST) $(ALL_WINDOWS_ARCH_LIST)

all-releases: $(gz_all_releases) $(zip_all_releases)

releases: $(gz_releases) $(zip_releases)

clean:
	rm $(BINDIR)/*


FORCE:
	cd $(basepath)/web; go-bindata -o ./public.go -pkg web ./public/...
	cd $(basepath)/config/base; go-bindata -o $(basepath)/config/base.go -pkg config ./...

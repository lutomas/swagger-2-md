JOBDATE		?= $(shell date -u +%Y-%m-%dT%H%M%SZ)
GIT_REVISION	= $(shell git rev-parse --short HEAD)
VERSION		?= $(shell git describe --tags --abbrev=0)

LDFLAGS		+= -s -w
LDFLAGS		+= -X github.com/lutomas/swagger-2-md/types/version.appVersion=$(VERSION)
LDFLAGS		+= -X github.com/lutomas/swagger-2-md/types/version.commit=$(GIT_REVISION)
LDFLAGS		+= -X github.com/lutomas/swagger-2-md/types/version.buildTime=$(JOBDATE)

LDFLAGS_LINUX		+= -linkmode external -extldflags -static

# ###########################
# BUILD
# ###########################

install-cli:
	@echo "++ Building SWAGGER-2-MD CLI binary (<current-os>)"
	go install -ldflags "$(LDFLAGS)" github.com/lutomas/swagger-2-md/cmd/swagger-2-md

# ###########################
# DEV
# ###########################

go-mod-get:
	@echo "Get all dependencies"
	cd cmd/swagger-2-md && go get .

go-mod-vendor:
	@echo "Prepare and clean dependencies"
	go mod vendor && go mod tidy



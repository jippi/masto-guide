# syntax=docker/dockerfile:1

######################################################################
# build go server
######################################################################

FROM golang:1.25-alpine AS masto-guide-dk-servers

WORKDIR /usr/src/app

ENV GOCACHE=/root/.cache/go-build
ENV CGO_ENABLED=0

# get depdendencies
COPY scripts/servers/go.* /usr/src/app
RUN --mount=type=cache,id=go-pkg-mod,target=/go/pkg/mod \
    set ex \
    && go mod download -x \
    && go mod verify

# copy source and build binary
COPY scripts/servers/ /usr/src/app
RUN --mount=type=cache,id=go-build,target=/root/.cache \
    set -ex \
    && go build -o masto-guide-dk-servers .

######################################################################
# final image w/ mkdocs
######################################################################

FROM squidfunk/mkdocs-material

RUN --mount=type=cache,id=mkdocs-cache,mode=0777,target=/root/.cache \
    --mount=type=cache,id=apk-cache,target=/var/cache/apk \
    --mount=type=cache,id=apk-pkg-cache,target=/etc/apk/cache \
    set -ex \
    && apk upgrade --update-cache -a \
    && apk add \
        cairo \
        freetype-dev \
        gcc \
        git \
        git-fast-import \
        jpeg-dev \
        libffi-dev \
        musl-dev \
        openssh \
        zlib-dev \
    && pip install --upgrade pip \
    && pip install \
        'mkdocs-glightbox' \
        'mkdocs-redirects' \
        "mkdocs-charts-plugin" \
        "mkdocs-macros-plugin" \
        "mkdocs-minify-plugin>=0.3" \
        "mkdocs-redirects>=1.0" \
        "cairosvg>=2.5" \
        "pillow>=9.0"

COPY --from=masto-guide-dk-servers /usr/src/app/masto-guide-dk-servers /bin

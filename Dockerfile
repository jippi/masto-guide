# syntax=docker/dockerfile:1

######################################################################
# build go server
######################################################################

FROM golang:1.21-alpine AS masto-guide-dk-servers

WORKDIR /usr/src/app

# get depdendencies
COPY scripts/servers/go.* /usr/src/app
RUN set ex \
    && go mod download \
    && go mod verify

# copy source and build binary
COPY scripts/servers/ /usr/src/app
RUN set -ex \
    && go build -o masto-guide-dk-servers .

######################################################################
# final image w/ mkdocs
######################################################################

FROM ghcr.io/afritzler/mkdocs-material


RUN set -ex \
    && apk upgrade --update-cache -a \
    && apk add --no-cache \
    cairo \
    freetype-dev \
    git \
    git-fast-import \
    jpeg-dev \
    openssh \
    zlib-dev \
    && apk add --no-cache --virtual .build \
    gcc \
    libffi-dev \
    musl-dev \
    && pip install --upgrade pip \
    && pip install \
    'mkdocs-glightbox' \
    'mkdocs-redirects' \
    "mkdocs-charts-plugin" \
    "mkdocs-macros-plugin" \
    "mkdocs-minify-plugin>=0.3" \
    "mkdocs-redirects>=1.0" \
    "cairosvg>=2.5" \
    "pillow>=9.0" \
    && apk del .build \
    && rm -rf /tmp/* /root/.cache

COPY --from=masto-guide-dk-servers /usr/src/app/masto-guide-dk-servers /bin

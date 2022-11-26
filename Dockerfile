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
        "mkdocs-macros-plugin" \
        "mkdocs-minify-plugin>=0.3" \
        "mkdocs-redirects>=1.0" \
        "pillow>=9.0" \
        "cairosvg>=2.5" \
    && apk del .build \
    && rm -rf /tmp/* /root/.cache

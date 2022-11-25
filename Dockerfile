FROM squidfunk/mkdocs-material

RUN set -ex && \
    pip install mkdocs-macros-plugin

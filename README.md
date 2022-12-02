# masto-guide

The goal of Masto Guide is to provide a great introduction into Mastodon.

## Running locally

The easy way requires [Docker for Desktop](https://www.docker.com/get-started/) on your device to work.

### With Docker

Simply running `./scripts/mkdocs.sh`  will start the `mkdocs` dev server on port `8000`. Once the container is running, you should be able to access the documentation via `http://localhost:8000` in your browser.

### Without Docker

Please follow the [MkDocs install guide](https://www.mkdocs.org/user-guide/installation/) to get MkDocs installed along with Python.

Once MkDocs has been installed, install these additional plugins required for MastoGuide:

```shell
pip install \
    "mkdocs-minify-plugin>=0.3" \
    "mkdocs-redirects>=1.0" \
    "pillow>=9.0" \
    "cairosvg>=2.5" \
    "mkdocs-material" \
    "mkdocs-macros-plugin"
```

## License

Masto Guide © 2022 by Christian Winther is licensed under [Attribution-NonCommercial-ShareAlike 4.0 International](http://creativecommons.org/licenses/by-nc-sa/4.0/?ref=chooser-v1)

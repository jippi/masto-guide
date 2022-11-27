# Use your own user @ domain for Mastodon discoverability

!!! success "Contributions welcome - this guide is Open Source!"

    If you find a bug, issue or want to contribute a new webserver, language or framework, please [click the small pencil in the upper right part of this page next to the page title](https://github.com/jippi/masto-guide/blob/main/docs/guide/use-your-own-domain.md){target=_blank}. This page is a normal Markdown page using simple `string.replace()` logic to interpolate values, and built with `mkdocs`

This guide will help you configure various hosting platforms, webservers and systems so you can use your own domain as an alias to your Mastodon account.

Basically, it will allow your user `my-user@mastodon.social` to be found as `me@my-domain.com`.

## Getting started

!!! tip "Filling out the two fields will automatically generate valid configurations in all the examples below, with your bespoke values."

    It's highly recommended to fill out both fields for the best experience. If left empty, you will see placeholders instead.

!!! example

    1. My Mastodon username is `jippi@expressional.social`
    1. And I want to be discoverable as `jippi@jippi.dev`

1. My Mastodon username is <input id="username" placeholder="user@mastodon-server.com" class="username-selector" /> `← your current Mastodon username`
1. And I want to be discoverable as <input id="alias" placeholder="user@your-domain.com" class="username-selector" /> `← username@custom-domain-you-own.com`

## Read before starting

!!! warning "This will **not** change your Mastodon username. Only how you can be found on Mastodon."

    These changes will only make `__USER_NAME__@__USER_DOMAIN__` discoverable/searchable via `__ALIAS_NAME__@__ALIAS_DOMAIN__`.

    People will still see `__USER_NAME__@__USER_DOMAIN__` when they follow you, write to you, or you writing to them.

!!! tip "You can find the right webserver, hosting platform or framework in the left menu"

Much like with e-mail, you want folks to have an easy address to find you, and one that you can keep giving out to everyone even if later you switch to a different Mastodon server. A bit like e-mail forwarding to your ISP’s e-mail service.

**The good news is:** You can use your own domain and share it with other folks.

In Mastodon, `Users` (in Mastodon lingo called `Actors`) are discovered using [WebFinger](https://webfinger.net/){target=_blank}, a way to attach information to an email address, or other online resource. WebFinger lives on `/.well-known/webfinger` on a server.

When someone searches for you on Mastodon, your server will be queried for accounts using an endpoint that looks like this:

```http
GET https://__USER_DOMAIN__/.well-known/webfinger?resource=acct:__USER_NAME__@__USER_DOMAIN__
```

In the code examples on the page below, we're implementing the WebFinger endpoint on *your* domain and redirecting from *your* domain (`__ALIAS_DOMAIN__`) to the Mastodon servers (`__USER_DOMAIN__`) WebFinger endpoint, effectively making `__ALIAS_NAME__@__ALIAS_DOMAIN__` become an alias of `__USER_NAME__@__USER_DOMAIN__`.

### How does it works

!!! tip "The WebFinger request flow looks like this when your alias is set up"

This is a slightly more technical overview on how the alias functionality work in practice.

1. A user searches for `__ALIAS_NAME__@__ALIAS_DOMAIN__` on a Mastodon server called `example.com`
1. The `example.com` Mastodon server queries the WebFinger endpoint at `__ALIAS_DOMAIN__`, asking for the `__ALIAS_NAME__@__ALIAS_DOMAIN__` user.
    1. `GET https://__ALIAS_DOMAIN__/.well-known/webfinger?resource=acct:__ALIAS_NAME__@__ALIAS_DOMAIN__`
1. The webserver responsible for `__ALIAS_DOMAIN__` accepts the request and redirects the `example.com` Mastodon server to the `__USER_DOMAIN__` Mastodon server
    1. The webserver responds with `HTTP/1.1 301 Moved Permanently`
    1. The webserver with `Location: https://__USER_DOMAIN__/.well-known/webfinger?resource=acct:__USER_NAME__@__USER_DOMAIN__`
1. The `example.com` Mastodon server follows the redirect and query the `__USER_DOMAIN_` WebFinger endpoint.
    1. `GET https://__USER_DOMAIN__/.well-known/webfinger?resource=acct:__USER_NAME__@__USER_DOMAIN__`
1. The `example.com` Mastodon server gets a valid `__USER_DOMAIN__` Mastodon profile back named `__USER_NAME__@__USER_DOMAIN__` and shows the result to the user that searched for you.

## Web servers

### Apache (.htaccess)

!!! note "This configuration will redirect all usernames on `@__ALIAS_DOMAIN__` to your Mastodon account `__USER_NAME__@__USER_DOMAIN__`"

Put the following in your `.htaccess` file

```htaccess
RewriteEngine On
RewriteRule ^.well-known/host-meta(.*)$ https://__USER_DOMAIN__/.well-known/host-meta?resource=acct:__USER_NAME__@__USER_DOMAIN__ [L,R=301]
RewriteRule ^.well-known/webfinger(.*)$ https://__USER_DOMAIN__/.well-known/webfinger?resource=acct:__USER_NAME__@__USER_DOMAIN__ [L,R=301]
RewriteRule ^.well-known/nodeinfo(.*)$ https://__USER_DOMAIN__/.well-known/nodeinfo?resource=acct:__USER_NAME__@__USER_DOMAIN__ [L,R=301]
```

### Nginx

The following will set up nginx to redirect from your custom domain to your Mastodon account

#### Before server{} block

Put the following in your configuration file **before** the `server{}` block

```nginx
# if you're not mapping $request_path already, you'll need to this next block
map $request_uri $request_path {
    ~(?<captured_path>[^?]*) $captured_path;
}

map $arg_resource $valid_mastodon {
    # If you want any account at *@__ALIAS_DOMAIN__ to redirect to your Mastodon account e.x.
    #
    #       __USER_NAME__-@__ALIAS_DOMAIN__
    #       __ALIAS_NAME__-@__ALIAS_DOMAIN__
    #       anything-else@__ALIAS_DOMAIN__
    #
    # and so on all pointing to
    #
    #       __USER_NAME__@__USER_DOMAIN__
    #
    # then change the 'default' value below from 0 to 1 below.
    default 0;

    # If you want limit the accounts from @__ALIAS_DOMAIN__, add them individually
    # in the list below like this.
    #
    # NOTE: The value need to be url encoded like this
    #
    #      replace : with %3A
    #      replace @ with %40
    #
    # Example
    #
    #      acct:__USER_NAME__@__USER_DOMAIN__ becomes acct%3A__USER_NAME__%40__USER_DOMAIN__
    #      acct:__ALIAS_NAME__@__ALIAS_DOMAIN__ becomes acct%3__ALIAS_NAME__%40__ALIAS_DOMAIN__
    #
    # Add as many additional ones as you would like

    # acct:__USER_NAME__@__USER_DOMAIN__
    'acct%3A__USER_NAME__%40__USER_DOMAIN__' 1;

    # acct:__ALIAS_NAME__@__ALIAS_DOMAIN__
    'acct%3A__ALIAS_NAME__%40__ALIAS_DOMAIN__' 1;
}
```

#### Inside your server{} block

You can find the right configuration file by looking through the configuration files containing `server{}` blocks. The right one should look like this:

```nginx
server {
    #... other config

    server_name __ALIAS_DOMAIN__;

    #... other config
}
```

Add the following `location` clause inside your `server{}` for your domain configuration.

```nginx

    #... other config

    location ~ ^/.well-known/(host-meta|webfinger|nodeinfo) {
        if ($valid_mastodon = 1) {
            return 301 https://__USER_DOMAIN__$request_path?resource=acct:__USER_NAME__@__USER_DOMAIN__;
        }

        if ($valid_mastodon = 0) {
            # Mastodon account not in allowed list so return 404
            return 404;
        }
    }

    #... other config
```

## Serverless

### Firebase hosting

!!! info "Original implementation by [johnmu.com](https://johnmu.com/2022-mastodon-for-your-domain//){target=_blank}, please see their blog post for more information"

!!! note "This configuration will redirect all usernames on `@__ALIAS_DOMAIN__` to your Mastodon account `__USER_NAME__@__USER_DOMAIN__`"

Add the following to your `firebase.json` file

```json
{
  "hosting": {
    // ...
    "redirects": [
    {
      "source": "/.well-known/host-meta",
      "destination": "https://mastodon.social/.well-known/host-meta?resource=acct:__USER_NAME__@__USER_DOMAIN__",
      "type": 301
    },
    {
      "source": "/.well-known/webfinger",
      "destination": "https://mastodon.social/.well-known/host-meta?resource=acct:__USER_NAME__@__USER_DOMAIN__",
      "type": 301
    },
    {
      "source": "/.well-known/nodeinfo",
      "destination": "https://mastodon.social/.well-known/host-meta?resource=acct:__USER_NAME__@__USER_DOMAIN__",
      "type": 301
    }
    // ...
    ]
}
```

### Cloudflare Pages

!!! info "Original implementation by [jacobian.org](https://jacobian.org/til/my-mastodon-instance/){target=_blank}, please see their blog post for more information"

!!! note "This configuration will redirect all usernames on `@__ALIAS_DOMAIN__` to your Mastodon account `__USER_NAME__@__USER_DOMAIN__`"

Add the following three lines to your redirect rules file (`_redirects`):

```
/.well-known/host-meta* https://__USER_DOMAIN__/.well-known/host-meta:splat 301
/.well-known/webfinger* https://__USER_DOMAIN__/.well-known/webfinger:splat 301
/.well-known/nodeinfo* https://__USER_DOMAIN__/.well-known/nodeinfo:splat 301
```

## Static file generator

!!! note "This configuration will redirect all usernames on `@__ALIAS_DOMAIN__` to your Mastodon account `__USER_NAME__@__USER_DOMAIN__`"

Assuming your public web directory for your website is `www/public`

### .well-known/webfinger

!!! info "If you only care about Mastodon, this is the only file you need"

!!! example "Put the file in your public web directory for your website"

    Example: `www/public/.well-known/webfinger`

```json
{
  "subject": "acct:__USER_NAME__@__USER_DOMAIN__",
  "aliases": [
    "https://__USER_DOMAIN__/@__USER_NAME__",
    "https://__USER_DOMAIN__/users/__USER_NAME__"
  ],
  "links": [
     {
      "rel": "http://webfinger.net/rel/profile-page",
      "type": "text/html",
      "href": "https://__USER_DOMAIN__/@__USER_NAME__"
    },
    {
      "rel": "self",
      "type": "application/activity+json",
      "href": "https://__USER_DOMAIN__/users/__USER_NAME__"
    },
    {
      "rel": "http://ostatus.org/schema/1.0/subscribe",
      "template": "https://__USER_DOMAIN__/authorize_interaction?uri={uri}"
    }
  ]
}
```

### .well-known/host-meta

!!! info "Optional, only needed if you want to support discovery via [OASIS Open](https://www.oasis-open.org/){target=_blank}"

!!! example "Put the file in your public web directory for your website"

    Example: `www/public/.well-known/host-meta`

```xml
<?xml version="1.0" encoding="UTF-8"?>
<XRD xmlns="http://docs.oasis-open.org/ns/xri/xrd-1.0">
  <Link rel="lrdd" template="https://__USER_DOMAIN__/.well-known/webfinger?resource=acct:__USER_NAME__@__USER_DOMAIN__" />
</XRD>
```

### .well-known/nodeinfo

!!! info "Optional, only needed if you want to support discovery via [Diaspora](https://diasporafoundation.org/){target=_blank}"

!!! example "Put the file in your public web directory for your website"

    Example: `www/public/.well-known/nodeinfo`

```json
{
    "links": [
        {
            "rel": "http://nodeinfo.diaspora.software/ns/schema/2.0",
            "href": "https://__USER_DOMAIN__/nodeinfo/2.0"
        }
    ]
}
```

## Frameworks

### Wordpress

!!! note "This configuration will redirect all usernames on `@__ALIAS_DOMAIN__` to your Mastodon account `__USER_NAME__@__USER_DOMAIN__`"

1. Install the [Redirection Plugin](https://wordpress.org/plugins/redirection/){target=_blank}
1. `Source URL` must be `/.well-known/webfinger`
1. `Query Parameters` must be `Ignore all parameters`
1. `Title` can be empty
1. `Match` must be `URL only`
1. `When matches` must be "`Redirect to URL` with HTTP Code `302 - Found`
1. `Target URL` must be `https://__USER_DOMAIN__/.well-known/webfinger?resource=acct:__USER_NAME__@__USER_DOMAIN__`
1. Save the form

<figure markdown>
![Export page](img/wordpress-redirection-plugin.png){loading=lazy}
<figcaption style="float: inherit; width: auto">This is how the Wordpress Redirection form look.</figcaption>
</figure>

<div style="clear: both" />

### Django (Python)

!!! info "Original implementation by [aeracode.org](https://aeracode.org/2022/11/01/fediverse-custom-domains/){target=_blank}, please see their blog post for more information"

!!! note "This configuration will redirect all usernames on `@__ALIAS_DOMAIN__` to your Mastodon account `__USER_NAME__@__USER_DOMAIN__`"

Add 3 views to your Django application like this

```python
from proxy.views import proxy_view

def wellknown_webfinger(request):
    remote_url = ("https://__USER_DOMAIN__/.well-known/webfinger?" + request.META["QUERY_STRING"])
    return proxy_view(request, remote_url)


def wellknown_hostmeta(request):
    remote_url = ("https://__USER_DOMAIN__/.well-known/host-meta?" + request.META["QUERY_STRING"])
    return proxy_view(request, remote_url)


def wellknown_nodeinfo(request):
    remote_url = "https://__USER_DOMAIN__/.well-known/nodeinfo"
    return proxy_view(request, remote_url)
```

This uses the django-proxy package to provide the proxy_view. You might also want to put in a view that redirects yourdomain.com/@username:

```python
from django.http import HttpResponseRedirect

def username_redirect(request):
    return HttpResponseRedirect("https://__USER_DOMAIN__/@__USER_NAME__")
```

Hooking these up to the right URLs is the only other thing that's needed:

```python
urlpatterns = [
    ...
    # Fediverse
    path(".well-known/webfinger", blog.wellknown_webfinger),
    path(".well-known/host-meta", blog.wellknown_hostmeta),
    path(".well-known/nodeinfo", blog.wellknown_nodeinfo),
    path("@__USER_NAME__", blog.username_redirect),
]
```

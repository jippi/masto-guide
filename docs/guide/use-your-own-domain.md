# Use your own user @ domain for Mastodon discoverability

!!! success "Contributions welcome - this guide is Open Source!"

    If you find a bug or issue or want to contribute a new web server, language, or framework, please [click the small pencil in the upper right part of this page next to the page title](https://github.com/jippi/masto-guide/blob/main/docs/guide/use-your-own-domain.md){target=_blank}. This page is a normal Markdown page using simple `string.replace()` logic to interpolate values, and built with `mkdocs`

This guide will help you configure various hosting platforms, web servers, and systems so you can use your own domain as an alias for your Mastodon account.

It will allow your user `my-user@mastodon.social` to be found as `me@my-domain.com`.

## 1. Getting started

!!! tip "Filling out the two fields will automatically generate valid configurations in all the examples below, with your bespoke values."

    It's highly recommended to fill out both fields for the best experience. If left empty, you will see placeholders instead.

!!! example

    1. My Mastodon username is `jippi@expressional.social`
    1. And I want to be discoverable as `jippi@jippi.dev`

1. My Mastodon username is <input id="username" placeholder="user@mastodon-server.com" class="username-selector" /> `← your current Mastodon username`
1. And I want to be discoverable as <input id="alias" placeholder="user@your-domain.com" class="username-selector" /> `← username@custom-domain-you-own.com`

## 2. Important information

!!! warning "This will **not** change your Mastodon username. Only how you can be found on Mastodon."

    These changes will only make `${MASTODON_USER}@${MASTODON_DOMAIN}` discoverable/searchable via `${ALIAS_USER}@${ALIAS_DOMAIN}`.

    People will still see `${MASTODON_USER}@${MASTODON_DOMAIN}` when they follow you, write to you, or you write to them.

!!! tip "You can find the right web server, hosting platform, or framework in the left menu"

Much like with e-mail, you want folks to have an easy address to find you and one that you can keep giving out to everyone, even if later you switch to a different Mastodon server. A bit like e-mail forwarding to your ISP’s e-mail service.

**The good news is:** You can use your own domain and share it with other folks.

In Mastodon, `Users` (in Mastodon lingo called `Actors`) are discovered using [WebFinger](https://webfinger.net/){target=_blank}, a way to attach information to an email address or other online resource. WebFinger lives on `/.well-known/webfinger` on a server.

When someone searches for you on Mastodon, your server will be queried for accounts using an endpoint that looks like this:

```http
GET https://${MASTODON_DOMAIN}/.well-known/webfinger?resource=acct:${MASTODON_USER}@${MASTODON_DOMAIN}
```

In the code examples on the page below, we're implementing the WebFinger endpoint on *your* domain and redirecting from *your* domain (`${ALIAS_DOMAIN}`) to the Mastodon servers (`${MASTODON_DOMAIN}`) WebFinger endpoint, effectively making `${ALIAS_USER}@${ALIAS_DOMAIN}` become an alias of `${MASTODON_USER}@${MASTODON_DOMAIN}`.

### 2.1 How it works

!!! tip "The WebFinger request flow looks like this when your alias is set up"

This is a slightly more technical overview of how the alias functionality works.

1. A user searches for `${ALIAS_USER}@${ALIAS_DOMAIN}` on a Mastodon server called `example.com`
1. The `example.com` Mastodon server queries the WebFinger endpoint at `${ALIAS_DOMAIN}`, asking for the `${ALIAS_USER}@${ALIAS_DOMAIN}` user.
    1. `GET https://${ALIAS_DOMAIN}/.well-known/webfinger?resource=acct:${ALIAS_USER}@${ALIAS_DOMAIN}`
1. The web server responsible for `${ALIAS_DOMAIN}` accepts the request and redirects the `example.com` Mastodon server to the `${MASTODON_DOMAIN}` Mastodon server
    1. The web server responds with `HTTP/1.1 301 Moved Permanently`
    1. The web server with `Location: https://${MASTODON_DOMAIN}/.well-known/webfinger?resource=acct:${MASTODON_USER}@${MASTODON_DOMAIN}`
1. The `example.com` Mastodon server follows the redirect and queries the `__MASTODON_DOMAIN_` WebFinger endpoint.
    1. `GET https://${MASTODON_DOMAIN}/.well-known/webfinger?resource=acct:${MASTODON_USER}@${MASTODON_DOMAIN}`
1. The `example.com` Mastodon server gets a valid `${MASTODON_DOMAIN}` Mastodon profile back named `${MASTODON_USER}@${MASTODON_DOMAIN}` and shows the result to the user that searched for you.

## 3. Web servers

### Apache (.htaccess)

!!! note "This configuration will redirect all usernames on `@${ALIAS_DOMAIN}` to your Mastodon account `${MASTODON_USER}@${MASTODON_DOMAIN}`"

Put the following in your `.htaccess` file.

```htaccess
RewriteEngine On
RewriteRule ^.well-known/host-meta(.*)$ https://${MASTODON_DOMAIN}/.well-known/host-meta?resource=acct:${MASTODON_USER}@${MASTODON_DOMAIN} [L,R=301]
RewriteRule ^.well-known/webfinger(.*)$ https://${MASTODON_DOMAIN}/.well-known/webfinger?resource=acct:${MASTODON_USER}@${MASTODON_DOMAIN} [L,R=301]
RewriteRule ^.well-known/nodeinfo(.*)$ https://${MASTODON_DOMAIN}/.well-known/nodeinfo?resource=acct:${MASTODON_USER}@${MASTODON_DOMAIN} [L,R=301]
```

### Nginx

The following will set up nginx to redirect from your custom domain to your Mastodon account.

#### Before server{} block

Put the following in your configuration file **before** the `server{}` block

```nginx
# if you're not mapping $request_path already, you'll need to this next block
map $request_uri $request_path {
    ~(?<captured_path>[^?]*) $captured_path;
}

map $arg_resource $valid_mastodon {
    # If you want any account at *@${ALIAS_DOMAIN} to redirect to your Mastodon account e.x.
    #
    #       ${MASTODON_USER}-@${ALIAS_DOMAIN}
    #       ${ALIAS_USER}-@${ALIAS_DOMAIN}
    #       anything-else@${ALIAS_DOMAIN}
    #
    # and so on, all pointing to
    #
    #       ${MASTODON_USER}@${MASTODON_DOMAIN}
    #
    # then change the 'default' value below from 0 to 1 below.
    default 0;

    # If you want limit the accounts from @${ALIAS_DOMAIN}, add them individually
    # in the list below like this.
    #
    # NOTE: The value needs to be URL encoded like this
    #
    #      replace : with %3A
    #      replace @ with %40
    #
    # Example
    #
    #      acct:${MASTODON_USER}@${MASTODON_DOMAIN} becomes acct%3A${MASTODON_USER}%40${MASTODON_DOMAIN}
    #      acct:${ALIAS_USER}@${ALIAS_DOMAIN} becomes acct%3${ALIAS_USER}%40${ALIAS_DOMAIN}
    #
    # Add as many additional ones as you would like

    # acct:${MASTODON_USER}@${MASTODON_DOMAIN}
    'acct%3A${MASTODON_USER}%40${MASTODON_DOMAIN}' 1;

    # acct:${ALIAS_USER}@${ALIAS_DOMAIN}
    'acct%3A${ALIAS_USER}%40${ALIAS_DOMAIN}' 1;
}
```

#### Inside your server{} block

You can find the right configuration file by looking through the configuration files containing `server{}` blocks. The right one should look like this:

```nginx
server {
    #... other configs

    server_name ${ALIAS_DOMAIN};

    #... other configs
}
```

Add the following `location` clause inside your `server{}` for your domain configuration.

```nginx

    #... other configs

    location ~ ^/.well-known/(host-meta|webfinger|nodeinfo) {
        if ($valid_mastodon = 1) {
            return 301 https://${MASTODON_DOMAIN}$request_path?resource=acct:${MASTODON_USER}@${MASTODON_DOMAIN};
        }

        if ($valid_mastodon = 0) {
            # Mastodon account not in the allowed list, so return 404
            return 404;
        }
    }

    #... other configs
```

## 4. Serverless

### Firebase hosting

!!! info "Original implementation by [johnmu.com](https://johnmu.com/2022-mastodon-for-your-domain//){target=_blank}, please see their blog post for more information"

!!! note "This configuration will redirect all usernames on `@${ALIAS_DOMAIN}` to your Mastodon account `${MASTODON_USER}@${MASTODON_DOMAIN}`"

Add the following to your `firebase.json` file.

```json
{
  "hosting": {
    // ...
    "redirects": [
    {
      "source": "/.well-known/host-meta",
      "destination": "https://mastodon.social/.well-known/host-meta?resource=acct:${MASTODON_USER}@${MASTODON_DOMAIN}",
      "type": 301
    },
    {
      "source": "/.well-known/webfinger",
      "destination": "https://mastodon.social/.well-known/host-meta?resource=acct:${MASTODON_USER}@${MASTODON_DOMAIN}",
      "type": 301
    },
    {
      "source": "/.well-known/nodeinfo",
      "destination": "https://mastodon.social/.well-known/host-meta?resource=acct:${MASTODON_USER}@${MASTODON_DOMAIN}",
      "type": 301
    }
    // ...
    ]
}
```

### Cloudflare Pages

!!! info "Original implementation by [jacobian.org](https://jacobian.org/til/my-mastodon-instance/){target=_blank}, please see their blog post for more information"

!!! note "This configuration will redirect all usernames on `@${ALIAS_DOMAIN}` to your Mastodon account `${MASTODON_USER}@${MASTODON_DOMAIN}`"

Add the following three lines to your redirect rules file (`_redirects`):

```
/.well-known/host-meta* https://${MASTODON_DOMAIN}/.well-known/host-meta:splat 301
/.well-known/webfinger* https://${MASTODON_DOMAIN}/.well-known/webfinger:splat 301
/.well-known/nodeinfo* https://${MASTODON_DOMAIN}/.well-known/nodeinfo:splat 301
```

### Netlify

!!! info "This guide is based on the [netlify-plugin-mastodon-alias](https://github.com/dkundel/netlify-plugin-mastodon-alias){target=_blank} GitHub project."

1. Install `netlify-plugin-mastodon-alias`
    ```shell
    npm install netlify-plugin-mastodon-alias
    ```
2. Configure WebFinger by adding the following to your `netlify.toml` file:
    ```toml
    [[plugins]]
    package = "netlify-plugin-mastodon-alias"

    [plugins.inputs]
        username = "${MASTODON_USER}"
        instance = "${MASTODON_DOMAIN}"
        # delete or comment the next line if you want "*@${ALIAS_DOMAIN}" to work rather than just "${ALIAS_USER}@${ALIAS_DOMAIN}"
        strictUsername = "${ALIAS_USER}"
    ```
3. Deploy netlify

## 5. Static files

!!! note "This configuration will redirect all usernames on `@${ALIAS_DOMAIN}` to your Mastodon account `${MASTODON_USER}@${MASTODON_DOMAIN}`"

Assuming your public web directory for your website is `www/public`

### .well-known/webfinger

!!! info "If you only care about Mastodon, this is the only file you need"

!!! example "Put the file in your public web directory for your website"

    Example: `www/public/.well-known/webfinger`

```json
{
  "subject": "acct:${MASTODON_USER}@${MASTODON_DOMAIN}",
  "aliases": [
    "https://${MASTODON_DOMAIN}/@${MASTODON_USER}",
    "https://${MASTODON_DOMAIN}/users/${MASTODON_USER}"
  ],
  "links": [
     {
      "rel": "http://webfinger.net/rel/profile-page",
      "type": "text/html",
      "href": "https://${MASTODON_DOMAIN}/@${MASTODON_USER}"
    },
    {
      "rel": "self",
      "type": "application/activity+json",
      "href": "https://${MASTODON_DOMAIN}/users/${MASTODON_USER}"
    },
    {
      "rel": "http://ostatus.org/schema/1.0/subscribe",
      "template": "https://${MASTODON_DOMAIN}/authorize_interaction?uri={uri}"
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
  <Link rel="lrdd" template="https://${MASTODON_DOMAIN}/.well-known/webfinger?resource=acct:${MASTODON_USER}@${MASTODON_DOMAIN}" />
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
            "href": "https://${MASTODON_DOMAIN}/nodeinfo/2.0"
        }
    ]
}
```

## 6. Frameworks

### WordPress

!!! note "This configuration will redirect all usernames on `@${ALIAS_DOMAIN}` to your Mastodon account `${MASTODON_USER}@${MASTODON_DOMAIN}`"

1. Install the [Redirection Plugin](https://wordpress.org/plugins/redirection/){target=_blank}
1. `Source URL` must be `/.well-known/webfinger`
1. `Query Parameters` must be `Ignore all parameters`
1. `Title` can be empty
1. `Match` must be `URL only`
1. `When matches` must be "`Redirect to URL` with HTTP Code `302 - Found`
1. `Target URL` must be `https://${MASTODON_DOMAIN}/.well-known/webfinger?resource=acct:${MASTODON_USER}@${MASTODON_DOMAIN}`
1. Save the form

<figure markdown>
![Export page](img/wordpress-redirection-plugin.png){loading=lazy}
<figcaption style="float: inherit; width: auto">This is how the WordPress Redirection form look.</figcaption>
</figure>

<div style="clear: both" />

### Jekyll

1. Add `jekyll-mastodon_webfinger` to your Gemfile:
    ```shell
    bundle add jekyll-mastodon_webfinger
    ```
2. Add the plugin to your list of plugins in `_config.yml`:
    ```yaml
    plugins:
        - jekyll/mastodon_webfinger
    ```
3. Add your Mastodon username and instance to `_config.yml`:
    ```yaml
    mastodon:
        username: ${MASTODON_USER}
        instance: ${MASTODON_DOMAIN}
    ```

Next time you build the site, you will find a `/.well-known/webfinger` file in your output directory, and when you deploy you will be able to refer to your Mastodon account using your own domain.

### Django (Python)

!!! info "Original implementation by [aeracode.org](https://aeracode.org/2022/11/01/fediverse-custom-domains/){target=_blank}, please see their blog post for more information"

!!! note "This configuration will redirect all usernames on `@${ALIAS_DOMAIN}` to your Mastodon account `${MASTODON_USER}@${MASTODON_DOMAIN}`"

Add 3 views to your Django application like this:

```python
from proxy.views import proxy_view

def wellknown_webfinger(request):
    remote_url = ("https://${MASTODON_DOMAIN}/.well-known/webfinger?" + request.META["QUERY_STRING"])
    return proxy_view(request, remote_url)


def wellknown_hostmeta(request):
    remote_url = ("https://${MASTODON_DOMAIN}/.well-known/host-meta?" + request.META["QUERY_STRING"])
    return proxy_view(request, remote_url)


def wellknown_nodeinfo(request):
    remote_url = "https://${MASTODON_DOMAIN}/.well-known/nodeinfo"
    return proxy_view(request, remote_url)
```

This uses the `django-proxy` package to provide the proxy_view. You might also want to put in a view that redirects yourdomain.com/@username:

```python
from django.http import HttpResponseRedirect

def username_redirect(request):
    return HttpResponseRedirect("https://${MASTODON_DOMAIN}/@${MASTODON_USER}")
```

Hooking these up to the right URLs is the only other thing that's needed:

```python
urlpatterns = [
    ...
    # Fediverse
    path(".well-known/webfinger", blog.wellknown_webfinger),
    path(".well-known/host-meta", blog.wellknown_hostmeta),
    path(".well-known/nodeinfo", blog.wellknown_nodeinfo),
    path("@${MASTODON_USER}", blog.username_redirect),
]
```

### Remix (JavaScript)

!!! info "Original implementation by [Tom Sherman](https://tom-sherman.com/blog/mastodon-domain-alias-on-remix){target=_blank}, please see his blog post for more information"

!!! note "This configuration will only redirect the `${ALIAS_USER}` username on `@${ALIAS_DOMAIN}` to your Mastodon account `${MASTODON_USER}@${MASTODON_DOMAIN}`"

Add a new route at `app/routes/[.]well-known/webfinger.js` with content as follows:

```js
export function loader({ request }) {
  const url = new URL(request.url);
  const resourceQuery = url.searchParams.get("resource");

  if (!resourceQuery) {
    return new Response("Missing resource query parameter", {
      status: 400,
    });
  }

  if (resourceQuery !== "acct:${ALIAS_USER}@${ALIAS_DOMAIN}") {
    return new Response("Not found", {
      status: 404,
    });
  }

  return new Response(
    JSON.stringify({
      subject: "acct:${MASTODON_USER}@${MASTODON_DOMAIN}",
      aliases: [
        "https://${MASTODON_DOMAIN}/@${MASTODON_USER}",
        "https://${MASTODON_DOMAIN}/users/${MASTODON_USER}",
      ],
      links: [
        {
          rel: "http://webfinger.net/rel/profile-page",
          type: "text/html",
          href: "https://${MASTODON_DOMAIN}/@${MASTODON_USER}",
        },
        {
          rel: "self",
          type: "application/activity+json",
          href: "https://${MASTODON_DOMAIN}/users/${MASTODON_USER}",
        },
        {
          rel: "http://ostatus.org/schema/1.0/subscribe",
          template: "https://${MASTODON_DOMAIN}/authorize_interaction?uri={uri}",
        },
      ],
    }),
    {
      headers: {
        "content-type": "application/jrd+json",
      },
    }
  );
}
```
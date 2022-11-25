# Tips & tricks

## Don't delete your Twitter account!

If you deactivate/delete your Twitter account, your `@ handle` will be publicly available in 30 days. The chances that you will be impersonated with no recourse is very high! Don’t do it - especially if you had a fairly popular account.

To safely clean/wipe your account do this instead:

1. [Deactivate your Twitter account](https://twitter.com/settings/deactivate){target="_blank"}
2. Go to [twitter.com/login](https://twitter.com/login){target="_blank"} and reactivate it. All of your followers will be gone.
3. Now lock the account and let it rot.

You can use a tool like [semiphemeral.com](https://semiphemeral.com/){target="_blank"} for wiping data, including DMs.

## Hashtags are king and queen

With no AI or algorithm to push "relevant" content, a huge part of the discovery and exploration of Mastodon is based on hashtags. Use them liberally.

You can also subscribe to hashtags like you would a user:

- Search for the hashtag
- Click on the hashtag in the search result
- In the upper right corner of the list of posts (tweets) there is a small icon looking like a person with a `+` sign on the right-hand side. It's confusingly identical to the "follow user" icon. Click that, and you should get all posts with that hashtag in your regular timeline.

## Additional features

!!! info

        This only applies to users on our [{{ server_name }} Mastodon server]({{ server_url }}){target="_blank"}

This Mastodon server runs a [special version (in software terms, a `fork`) of Mastodon called Glitch](https://glitch-soc.github.io/docs/){target="_blank"} that adds some additional capabilities to the interface, filtering, and general experience.

You can read about [Glitch on their website](https://glitch-soc.github.io/docs/){target="_blank"}. If you don't care about the additional features, you can safely ignore them and use Mastodon like normal.

## Featured Hashtags

Did you see a person's profile has some special tags with numbers? Those are "Featured Hashtags." They are useful to present yourself more; for example, I often post about software, so I have `#Software` as a featured hashtag, so everyone can click it and find all **my** posts about software.

To create them:

- Go to [Preferences]({{ server_url }}/settings/){target="_blank"}
- See [Profile > Featured Hashtags]({{ server_url }}/settings/featured_tags){target="_blank"}
- You will be offered hashtags you often use
- Type one hashtag, like `#sport`, and `Add`
- Repeat the steps above to add more hashtags

## Advanced UI mode

If you are familiar with Tweetdeck or other 3rd-party Twitter apps, they often support multiple timelines/lane side by side rather than a single list of tweets like Twitter.

This is natively available on the Mastodon web app, and you can change it this way:

- Go to [Preferences]({{ server_url }}/settings/preferences/appearance){target="_blank"}
- Check `Enable advanced web interface`
- Click `Save Changes`
- Go back to [Mastodon]({{ server_url }}/){target="_blank"}
- Search a tag you love, for example, `#Food`, `#Sport`, or `#Nature`
- Click a tag found, then a column of that tag appears
- Click the "triple lines" button to the top-right, click `PIN`
- Click `UNPIN` if you want to remove a column

## Lists

Lists help you watch over certain people separately.

For example, you may create three lists of "My Family", "Colleagues", "Bots", which contain your family, and your colleagues, and the bots you follow.

You can [create and maintain your lists here]({{ server_url }}/lists) or in the menu under `Lists`

## Paste an URL/link in the Search Box

The search box supports searching for URLs/links to profiles and posts from other servers. This allows you to open a post or profile from a different server to follow, like, comment, and share the post or profile without opening a new window. This is a nice shortcut.

Example links you can search for:

- [{{ server_url }}/@jippi]({{ server_url }}/@{{ server_owner }}) - my profile - you can search this link on any server, and my profile will show up
- [{{ server_url }}/@jippi/109371126097435854]({{ server_url }}/@{{ server_owner }}/109371126097435854) - you can search this link on any server, and my post should show up

## Using your own domaain as alias

Mastodon makes it ~fairly easy for you to add `aliases` to your account so you can be found under names (e.g., your private domain like `jippi.dev`).

For example, my personal domain is `jippi.dev`, and I've set up so `@{{ server_owner }}@jippi.dev` and `@me@jippi.dev` will find my account `@{{ server_owner }}@{{ server_domain }}` across all servers. You can think of it as an e-mail redirect/forward.

It can be incredibly useful since it will allow people to find you via your regular e-mail. Since it's a redirect/forward, if you move to another server in the future, your alias can be updated to point to the new server transparently without updating many links across the web.

The article [Mastodon - be findable with your domain](https://johnmu.com/2022-mastodon-for-your-domain/){target="_blank"} has a lot of good examples on how to configure this, including (but not limited to)

* [Wordpress](https://johnmu.com/2022-mastodon-for-your-domain/#wordpress){target="_blank"}
* [Apache](https://johnmu.com/2022-mastodon-for-your-domain/#apache--co){target="_blank"}
* [nginx](https://gist.github.com/dwsmart/b9733545030cde7451f8688538b945ab){target="_blank"}

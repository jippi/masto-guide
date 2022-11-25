# Twitter and Mastodon

## Introduction

At a quick glance Mastodon and Twitter seem to have a lot similarities:

* `Tweet`? is called a `post`.
* `Like`? thats called `Favorite`.
* `Retweet`? Easy, thats called `Boost`.

Nice, and they also behave the same, right? **NO!** In fact, the Mastodon functions behave *very* differently from Twitter. Keep reading to learn more

## Boosts & favorites

* `Boost` will be reshared on your profile, increasing the reach of the original post. Your followers will see the post in their timeline.
* `Favourite` will be added to your favourites list, and a "favourite notification" will be delivered to the post author author. No one else will see your Favorite action. The post you're liking will *not* have increased reach. Your followrs will not see the post in their timeline.

!!! question "When should I `boost` and `favorite`?"

        Use `boost` every time you would have liked _or_ retweeted on Twitter.

        Use `favorite` as a `thank you` or `I've seen your message`. Only you and the post author can see when you favorite a post.

There are also areas where Mastodon has unique features you don't find on Twitter.

## No algorithm

This is not something Mastodon uses. At all. All timelines are strictly chronologic by default, and no automatic filtering will apply to what appears in your timeline. (You can, of course, tweak this via the [filters]({{ server_url }}/filters){target="_blank"}) menu in your profile. This also means `boosting`, and liberal usage of hashtags is much more important (and less intrusive than on Twitter) for discovering content.

## No ads

Ads doesn't exist on Mastodon. By default, no one is tracking you, analyzing your content and engagement to figure out how to show you relevant content. This is a huge difference from Twitter, which spend a lot of time figuring out what ads to show you. This also means that server owners have no revenue stream to cover the server's cost. Please consider donating to them if they have such an option - typically mentioned on their "About" page.

## Post are edible

Yes! You can edit your posts on Mastodon; in the web interface, click the `...` on the post and pick `Edit`. Any edits you make will automatically propagate to other servers, and there will be a small `*` mark on the post that shows it has been edited. Clicking the `*` will show the edits made to the post for full transparency.

## Post length

On Mastodon are 500 characters by default, more than double of a tweet.

## Bookmarks

They are similar to a `Like` but it doesn't "signal boost" like a `Retweet` or `Like` does on Twitter

## Servers

The home of your Mastodon account. Since Mastodon consists of many 1000s of servers (unlike Twitter, which is a single system)

## Direct Messages

They are deceptively similar but do have significant differences from Twitter in terms of how private they are. See [this article](https://www.slashgear.com/1090436/mastodon-dms-are-very-different-to-twitter-and-that-could-get-embarrassing/){target="_blank"}

## Pinned Post

They are like `Pinned tweet` on Twitter, but it allows you to select up to 5 pinned posts rather than a single Tweet on Twitter. Otherwise, it's working the same; the pinned posts will always be shown as the first posts on your profile page.

## Quoted tweets

They do not exist on Mastodon. It was frequently used in abusive ways on Twitter, so the Mastodon team has so far been reluctant to implement it on Mastodon. Since Mastodon has been a safe haven for minority groups that got trolled and harassed on Twitter long before the Musk takeover, Maston generally focuses on avoiding features that can be abused.

## Timelines

### Local timeline

The local timeline lists all posts created on the `server` you are on. It does not include any content from other servers and does not include activities like `boost`, `favorite`, and such.

### Federated timeline

The federated timeline list all posts your local Mastodon server has seen from *other* servers. This includes content from your direct following, from other people on the server following, and signal boosting. You can consider this your "neighborhood" of closely related servers and people.

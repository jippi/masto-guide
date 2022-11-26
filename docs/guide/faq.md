# FAQ

I've tried to answer your most common questions and concerns, followed by some Twitter -> Mastodon migration tips. Please let me know if anything is unclear or missing.

## What is this Mastodon server/fediverse thing anyway?

The way Mastodon works are very similar to how e-mail work. It doesn't matter if you have a `gmail.com`, `hotmail.com`, `jubii.dk` or your own custom domain for e-mail; you can easily and without any level of friction send and receive emails from everyone else who also has an email.

The Mastodon network function the same way. You can follow, direct message, boost (retweet), favorite, and everything else Mastodon offers with your friends on any other servers they might be on, and they can do the same.

I've explained this in more detail in the `But all my friends are on another server!` section further down on this page.

## But all my friends are on another server!

That's completely fine. The way Mastodon works are very similar to how e-mail work. It doesn't matter if you have a `gmail.com` `hotmail.com`, `jubii.dk` or your own custom domain for e-mail; you can easily and without any level of friction send and receive emails from everyone else who also has an email.

The Mastodon network function the same way. You can follow, direct message, boost (retweet), favorite, and everything else Mastodon offers with your friends on any other servers they might be on, and they can do the same.

Even search generally works between servers to some degree since the various Mastodon servers "gossip" with each other as people follow and mention each other across the network, slowly building up a social graph of other users discovered across the fediverse (federated universe).

Once you follow someone, when you want to mention them, like on Twitter, your `@their-username` will look the same if you're on the same server or different,

Example: If we follow each other but are on different servers, you would only have to use `@jippi` to mention me since the servers know we're related. You do not need to include the `@expressional.social` part of the user.

## How do I move away from your server?

Follow the steps in the [changing server](move-mastodon-server.md), but use different servers than in the example.

You will not lose your followers; the Mastodon network will make sure to migrate them to your new server transparently and forward any messages sent to your old account to the new one. It generally works really well.

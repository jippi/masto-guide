# Changing Mastodon server

!!! note Not all data will be migrated

    You will not lose your followers; the Mastodon network will make sure to migrate them to this server transparently and forward any messages sent to your old account to the new one. It generally works really well. People who has `private` accounts will need to accept your follow request again.

!!! warning

    Your existing posts, likes, bookmarks, people *you* follow and so on will *not* be migrated, they will remain on the old server until you delete your account (optional).

    You can manually move most this data, we'll cover that in [migration step 1](#1-export-your-data) and [migration step 2](#3-import-your-data) below.

You can check out the [official documentation on moving accounts between servers](https://docs.joinmastodon.org/user/moving/){target="_blank"} if this guide isn't clear.

This is a quick guide to moving your old user `@example-user@mastodon.social` to this server where your username is `@example-user@expressional.social`.

## 1. Export your data

!!! warning "Make sure to export your data from the *old* server"

!!! warning "Do not use `Request your archive` for migrating data, you can *not* import that file on the new server. It's great for regular backups though"

On the *other* server, go to your profile settings and navigate to the `Import and Export` menu in the left menu.

In the table, click on the `⬇️ CSV` link for each kind of data you want to migrate to the new server.

* Follows
* You block
* You mute
* Domain blocks
* Bookmarks

Note that the following can't be imported

* Lists
* Followers - they are automatically handled [in the next step](#2-migrate-your-user)

## 2. Migrate your user

1. Create a new account on the **new** server.
1. Create an account alias on **old** server. (`Preferences -> Profile -> (scroll to the bottom) Moving from a different account`). You need to add the *full* username you use on the **old** server (ex.`@example-user@mastodon.social`). This configures your account to accept the account move from the *other* server.
1. Go to the `Profile` page on your *old* server, scroll to the bottom and select `Move to a different account`
    1. **Important**: Please read the warning text on the page to make sure you are okay with the consequences
    1. For the `Handle of the new account` field you enter your username on *this* server (ex. `@example-user@expressional.social`)
    1. For the `Password` field enter your password you used on the *other* server
1. Click `Move followers`

## 3. Import your data

1. On the **new** server navigate to `Import and export -> Import`
1. From the dropdown, select the kind of data you want to import (this is the `CSV` files from [step 1](#1-export-your-data))
1. Select the file you exported by clicking the file input (`Choose file`)
    1. `Following list` file is called something like `follows.csv` in your Downloads folder
    1. `Blocking list` file is called something like `blocks.csv` in your Downloads folder
    1. `Muting list` file is called something like `mutes.csv` in your Downloads folder
    1. `Domain blocking list` file is called something like `domain_blocks.csv` in your Downloads folder
    1. `Bookmarks` file is called something like `bookmarks.csv` in your Downloads folder
1. Make sure to select `Merge` and **not** `Overwrite`
1. Click `Upload`

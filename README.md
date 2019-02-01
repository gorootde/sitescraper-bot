# Sitescraper Telegram Bot
Telegram bot notifying you about website changes.

This bot will monitor the configured websites for changes, and in case of a change notify you via telegram. The message format can be freely defined using templates.

# TL;DR;
1. Create a `config.toml`
2. Configure the websites to be scraped and your bot-token in `config.toml`
3. You're good to go

# Configuration

Example for a single site
```toml
token="123456:ABCDEFGHIJKLMNOPQRSTUVWXYZ"
chatid=-1234515
[scrapers]
    [scrapers.firstpage]
        name="Example website"
        url="http://example.com"
        deeplinkselector=".teaser-float a"
        template="templates/default.tpl"
        [scrapers.example-site.fields]
            [scrapers.example-site.fields.imageurl]
                selector="meta[property='og:image']"
                attribute="content"
            [scrapers.example-site.fields.headline]
                selector=".headline-article span[itemprop='headline']"
            [scrapers.example-site.fields.kicker]
                selector=".headline-article > .kicker"
```

- `token` Your telegram bot token
- `chatid` The ID of the group chat where messages will be published
- `[scrapers]` section to define your scrapers.
- `[scraper.<sid>]` The section to specify your scraper configuration `<sid>` is an arbitrary unique string - your scraper ID.
    - `name` Human readable name of the site to be scraped
    - `url` The URL to use as entrypoint for scraping
    - `deeplinkselector` Optional line. If it is missing, the field search will be performed on the page specified by `url`. If it is present, it must contain a jQuery selector selecting the links 
    to follow.
    - `template` path to the template file used to render results
    - `[scrapers.<sid>.fields]` Section for configuring the fields to scrape
        -  `[scrapers.<sid>.fields.<fieldid>]` Section for specifing the settings for a scraped field. `<fieldid>`is an arbitrary string which must be unique within the current scraper section
            - `selector` jQuery selector selecting an HTML element on the page
            - `attribute` Optional. If present, it must contain the name of the attribute that you want to get as value. If not present, the scraper will use whatever is the text of the selected element.

# Templates 
Templates can be used to define the content of a message that is sent to you. Within a template you can use data scraped from the website. The context of such template contains the values for all fields configured for your scraper.

Create a file in `templates` folder and put go template code in it. Currently only [Telegram HTML formatting](https://core.telegram.org/bots/api#formatting-options) is supported.


# TODO List
Things still open to be implemented:
- Make bot monitor (and diff) sites, not only scrape a single time on startup
- Recursive deeplinks: Search deeplinks not only on the first level
- Automate Groupchat join and enable per-groupchat scrape configs
- Proper inline documentation of code
- Support other message types than text
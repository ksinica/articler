Basic article to EPUB converter. Tries to find the main readable content and remove clutter.

Given that `articles.txt` file contains newline separated list of article URLs, the usage is simple:
`cat articles.txt | articler myarticles.epub`

More complicated usage, fetches all the articles from RSS feed:
`curl -s http://feeds.bbci.co.uk/news/technology/rss.xml | xq -r '.rss.channel.item | .[] | .link' | articler news.epub`
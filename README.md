# nusetextd
NuseTextD is a simple beanstalkd consumer/worker that pulls URLs from the src
tube, posts the URL to TextRazor for NLP analysis, stores the results in MySQL
and pushes the article into the dest tube for further processing.

Though this is very old, it is a illustration of:
- pipeline workflow using queue system (beanstalkd)
- concurrency (multiple workers per connection)
- third-party integration with a REST API

This is part of a larger system - an RSS news reader, which pulls news articles 
from RSS feeds and performs NLP analysis, content categorisation and semantic 
analysis on the content.

Other tools related to this personal project include:
 - [JustText](https://github.com/JalfResi/justext)
 - [GoTidy](https://github.com/JalfResi/GoTidy)
 - GreatScott (unreleased cron service w/ WebAPI)

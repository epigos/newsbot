# News Messenger bot

A messenger news bot which crawls news articles and sent them to facebook messages upon receipt of message.

## Example chat commands

    send me news
    send me sports news
    give me news from yesterday
    give me news from BBC

## Third-party services

- [Dialogflow](https://dialogflow.com/docs)
- [Facebook messenger](https://developers.facebook.com/docs/messenger-platform/)

## Database

Uses [Gcloud datastore](https://cloud.google.com/datastore/docs/tools/datastore-emulator)

## Copy and update env file with appopriate settings

    cp env/sample.env env/local.env

## dev start

    go run main.go

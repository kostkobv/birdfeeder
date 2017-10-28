# [[https://github.com/kostkobv/birdfeeder/blob/master/docs/logo.jpg|alt=logo]]Birdfeeder
[![CircleCI](https://circleci.com/gh/kostkobv/birdfeeder.svg?style=svg)](https://circleci.com/gh/kostkobv/birdfeeder)
[![Go Report Card](https://goreportcard.com/badge/github.com/kostkobv/birdfeeder)](https://goreportcard.com/report/github.com/kostkobv/birdfeeder)
[![codecov](https://codecov.io/gh/kostkobv/birdfeeder/branch/master/graph/badge.svg)](https://codecov.io/gh/kostkobv/birdfeeder)

## Disclaimer
Create an API that accepts SMS messages submitted via a POST request containing a JSON object as request body.
- Example: 
```json
{ 
  "recipient": 31612345678,
  "originator":"MessageBird",
  "message":"This is a test message."
}
```

- Send the received message to the MessageBird REST API using one of our REST API libraries: https://github.com/messagebird/
- When an incoming message content/body is longer than 160 chars, split it into multiple parts (known as concatenated SMS)
- Make sure no empty or incorrect parameter values are send to MessageBird (input validation)
- The (theoretical/imaginary) throughput to MessageBird is one API request per second. Make sure the outgoing messages won’t exceed that limit, also when multiple incoming requests are received at the API or concatenated messages need to be send.

## API
###POST `/message`
```json
{ 
  "recipient": [valid recipient MSISDN],
  "originator": [valid originator accordingly to MessageBird documentation],
  "message": [message content]
}
```

Forwards message to MessageBird API 

## How does it work
[[https://github.com/kostkobv/birdfeeder/blob/master/docs/graph.png|alt=graph]]

**A** - Message is submitted to the project's API, validated, converted, gets generated UDH (if needed), splitted and pushed to the queue

**B** - Queue retrieves the collection of pushed messages (meanwhile queue is not locked but working collection is replaced with new collection on a fly every 1 second if needed). The retrieved collection is analyzed for identical messages but with different recipients so we could send more messages at once.

**C** - The message with the biggest amount of recipients would be sent first. The rest of the messages are sent back to the queue. If an error was returned by MessageBird API - the message is also sent back to the queue.    
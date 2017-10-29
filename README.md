![logo](https://github.com/kostkobv/birdfeeder/blob/master/docs/logo.png)

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
### POST `/message`

#### Description
Forwards message to MessageBird API

#### Required
`recipient`: valid recipient MSISDN,
`originator`: valid originator accordingly to MessageBird documentation,
`message`: message content

#### Response
##### Success `200`
Returns the submitted object as a confirmation of correct message

###### Example
```JSON
{ 
  "recipient": 31612345678,
  "originator":"MessageBird",
  "message":"This is a test message."
}
```

##### Unprocessable entity `422`
Returned in case if not valid message was submitted

###### Example
```JSON
{
    "body": "must have a value",
    "originator": "use valid MSISDN or alphanumeric value (max. 11 symbols long)",
    "recipient": "should be a valid MSISDN"
}
```

##### Bad Request `400`
Returned in case of invalid JSON submitted

###### Example
```JSON
{
    "message": "Syntax error: offset=20, error=invalid character '\"' after object key:value pair"
}
```

## How does it work
![graph](https://github.com/kostkobv/birdfeeder/blob/master/docs/graph.png)

**A** - Message is submitted to the project's API, validated, converted, gets generated UDH (if needed), splitted and pushed to the queue. Max splitted to 9 parts - rest of the message is discarded and not sent (taken from MessageBird documentation, however GSM documentation said there could be up to 255 parts). Depending on the encoding messages have different limit and split:
- plain encoding (symbols from GSM 03.38 table):
1. First message is 160 symbols (some special symbols are counted as 2. Read more: https://en.wikipedia.org/wiki/GSM_03.38).
2. If message is longer than 160 symbols then it is splitted by 153 symbols long parts. 1 part - 1 SMS

- unicode encoding (symbols not from GSM 03.38 table):
1. First message is 70 symbols
2. If message is longer than 70 symbols then it is splitted by 67 symbols long parts. 1 part - 1 SMS (_for some reason splitting doesn't work for MessageBird API_)

**B** - Queue retrieves the collection of pushed messages. Meanwhile queue is not locked but working collection is replaced with new collection on a fly every 1 second if needed. The retrieved collection is analyzed for identical messages but with different recipients so we could send more messages at once. Identical messages with the same recipient are considered to be the same message submitted twice so it's need to be delivered twice.

**C** - The message with the biggest amount of recipients would be sent first. The rest of the messages are sent back to the queue. If an error was returned by MessageBird API - the message is also sent back to the queue.

# Development

**Please, do not put the project into the `$GOPATH/src/github.com/kostkobv/birdfeeder`. Use [GVM](https://github.com/moovweb/gvm) to control your package sets and put the project straight into the `$GOPATH/src` of your package set.
You can find more information about installing and configuring GVM below.**

## Step 1: Install GVM
First install GVM
```bash
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
```
Or if you are using `zsh` just change `bash` with `zsh`

The project is based on Go 1.9.1
```bash
gvm install go1.9.1
```

If you're getting an error with installing the Go version try to use `--binary` flag
```bash
gvm install go1.9.1 --binary
```

Checkout to the required version
```bash
gvm use go1.9.1 --default
```

## Step 2: Utilities
Dependencies are represented by [Glide](https://github.com/Masterminds/glide) package manager. Please install it first.
```bash
curl https://glide.sh/get | sh
```
 
You would also need to install linter (optional).
```bash
go get -u gopkg.in/alecthomas/gometalinter.v1
```

Then install linters itself
```bash
gometalinter.v1 --install
```

To run linter
```bash
make lint
```

_Please notice that CI build would fail if changes wouldn't pass linter_

## Step 3: GVM pkgset
You need to link project to the pkgset
```bash
gvm pkgset create birdfeeder
gvm pkgset use birdfeeder
```

Now you need to link the project to `GOPATH`
```bash
rm -rf ~/.gvm/pkgsets/go1.9.1/birdfeeder/src
ln -s $PWD ~/.gvm/pkgsets/go1.9.1/birdfeeder/src
```

## Step 4: Vendor 

After you can install the vendor dependencies:
```
glide install
```

## Step 5: Configuration
You need to copy file `./config/settings.go.example` to `./config/settings.go` and place configuration you would like to use.

To run tests
```bash
make test
```
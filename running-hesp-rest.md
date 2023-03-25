# Using the HESP Api with CURL

Replace {client:key} and {channel:key} with the key that Stages has provided.

## About CURL

CURL is a small application that lets you run http commands from a terminal or shell.

A GET request is the default option and to run a POST, PUT, DELETE you have to specify that in the command with the following option ```-X POST``` (hyphen followed by a Captial X, a space, then the verb)

Tip: Use -v to get a more verbose output (including responce code, headers and so on)

## Checking status

Run the following command to check the status of a channel.

```bash
curl -H "Authorization: Basic {client:key}" https://api.theo.live/channels/{channel:key}/status
```

The response will come back as a JSON, looking like this

```json
{
    "code":"waiting",
    "message":"Channel is started, waiting for ingest"
}
```

## Starting a channel

Run the following command in a terminal/shell to start the channel.

```bash
curl -H "Authorization: Basic {client:key}" -X POST https://api.theo.live/channels/{channel:key}/start
```

If the request succeeded, the response code will be a 204 (no content). This will not be shown in a CURL request by default

If the request failed, there will be a response text showing the reason

## Starting a channel

Run the following command in a terminal/shell to start the channel.

```bash
curl -H "Authorization: Basic {client:key}" -X POST https://api.theo.live/channels/{channel:key}/stop
```

If the request succeeded, the response code will be a 204 (no content). This will not be shown in a CURL request by default

If the request failed, there will be a response text showing the reason.

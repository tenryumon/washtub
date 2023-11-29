# Pulse

URL: `{server} /worker/pulse`

Request:
`{
    "channel_id":"abcd",
    "address":"0.0.0.0",
    "topic":"topic",
    "sink_type":"http",
    "status":"active"
}`

Response:
`{
    "result_status": {
        "code": "200",
        "reason": "success",
        "message": null
    },
    "data": {
        "channel_id": "abcd",
        "address": "0.0.0.0",
        "topic": "topic",
        "sink_type": "http",
        "status": "active",
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z"
    }
}`

# Message

URL: `{server} /worker/message`

Request:
`{
    "channel_id":"abcd",
    "body":"{\"key\": \"value\"}",
    "status":"success"
}`

Response:
`{
    "result_status": {
        "code": "200",
        "reason": "success",
        "message": null
    },
    "data": {
        "channel_id":"abcd",
        "body":"{\"key\": \"value\"}",
        "status":"success"
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z"
    }
}`

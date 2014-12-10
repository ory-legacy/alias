alias
======

Stores human readable aliases (a paths) for identifiers. `ca64bfo` could be represented via `/about/team`.  
Each identifier (`ca64bfo`) may have multiple paths. To distinguish between them, **alias** stores a timestamp for each entry.

## Examples

The following examples will give you an API overview:

### Receiving

To receive an item by it's path, create a request which contains an uri_encoded string:
```
GET /%2Fabout%2Fteam HTTP/1.1
```

The response returns a JSON string which contains one item and all of its information:
```
HTTP/1.1 200 OK
Content-Type: application/json

{
  key: "/about/team",
  value: "tx86bb",
  created: "1994-11-05T13:15:30Z"
}
```

To receive an items aliases, create a request which uses a `value` parameter containing an uri_encoded string:
```
GET /?value=tx86bb HTTP/1.1
```

The response is formatted in JSON and returns the latest entry as wel as an ordered (order by created desc) list of items:
```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "key": "/about/team",
  "value": "tx86bb",
  "created": "2014-11-05T13:15:30Z"
  "list": [
    {
      "key": "/about/team",
      "value": "tx86bb",
      "created": "2014-11-05T13:15:30Z"
    },
    {
      "key": "/team",
      "value": "tx86bb",
      "created": "1994-11-05T13:15:30Z"
    },
  ]
}
```

If no element is found, an 404 response containing `{message: string}` will be returned.

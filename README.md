mapper
======

Provides types containing maps over the network. Keys within a type are unique whilst values are not.

## Concept

The primary use case of this app is to provide URI to ID mapping, e.g.:

Request:
```
GET /map/uri?key=home HTTP/1.1
HOST: api.baldur.io
```

Response:
```
HTTP/1.1 200 OK
Content-Type: application/json

{
  key: "home",
  value: "874192475123",
  created: "1994-11-05T13:15:30Z"
}
```

Request:
```
GET /map/uid?key=92475123 HTTP/1.1
HOST: api.baldur.io
```

Response:
```
HTTP/1.1 200 OK
Content-Type: application/json

{
  key: "92475123",
  value: "about",
  created: "1994-11-05T13:15:30Z"
}
```

### The map model

```
var item = {
  key: string,
  value: string,
  created: int
}
```

Both created and modified are unix timestamps.

### Get a map's value

URL: `GET /maps/:type?key=uri_encoded_key`

Returns:
```
item
```

Both created and modified are formatted in [UTC datetime](http://www.w3.org/TR/NOTE-datetime): `1994-11-05T13:15:30Z`

### Store a new map

Request: `POST /maps/:type`

Values:
```
{
  key: string,
  value: string
}
```

Returns:
```
 {
  item: item
  error: bool,
  errorMessage: string
}
```

### Remove a specific map:

Request: `DELETE /maps/:type?key=uri_encoded_key`

Returns:
```
{
  error: bool,
  errorMessage: string
}
```

mapper
======

Provides types containing maps over the network. Keys within a type are unique whilst values are not.

## Concept

The primary use case of this app is to provide URI to ID mapping, e.g.:

Request:
```
GET /map/uri?key=home HTTP/1.1
HOST: api.baldur.io
app-public-token: ba64ofk10fkk
```

Response:
```
HTTP/1.1 200 OK
Content-Type: application/json

{
  id: 1234,
  key: "home",
  value: "874192475123",
  type: "string"
  created: "1994-11-05T13:15:30Z"
}
```

## API

All requests need to contain your app token in the header.

```
app-public-token: you-token
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

### Searches keys by a value

`GET /maps/:type?value=uri_encoded_value`

Returns:
```
[item1, item2, ...]
```

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
  item: item
  error: bool,
  errorMessage: string
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

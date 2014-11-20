mapper
======

Provides versioned maps over the network.

## API

### Get a map's value

URL:
`GET /maps?key=uri_encoded_key`

Returns:
```
[{
  id: int,
  key: string,
  modified: timestamp
}, {...}]
```

### Get a map's key

`GET /maps?value=uri_encoded_value`

Returns:
```
[{
  id: int,
  value: string,
  modified: timestamp
}, {...}]
```

### Store a new map

Request: `POST /maps`

Values:
```
{
  key: string,
  value: string,
  error: bool,
  errorMessage: string
}
```

### Remove a specific map:

Request: `DELETE /maps/:id`

Returns:
```
{
  error: bool,
  errorMessage: string
}
```

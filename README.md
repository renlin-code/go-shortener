# Go Shortener API


Go Shortener is an API to shorten URLs and manage them.

![project banner](https://imgur.com/006eNXT.png)

## Building and running the App

This app uses Docker, which makes building easier. You just need to install [Docker](https://docs.docker.com/) on your system.

In addition to this you need to set your environment variables in a .env file. As an example see the [.env.example](https://github.com/renlin-code/go-shortener/blob/master/.env.example) file


Having Docker on your system and environment variables setted, you only need to run:

```bash
docker-compose up
```

## Documentation

In this app there are basically two types of endpoints. Open and private endpoints.

### Open endpoints

For these endpoints authorization is not required.

#### Create a short URL

It generates a short URL from a base URL provided as a parameter. For this short URL is assigned a random alias.

```http
  POST /
```

| Parameter (body) | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `url` | `string` | **Required**. Your base URL to shorten |

---

#### Redirect a short URL

Redirect a short URL to its base URL.

```http
  GET /${alias}
```

| Parameter (url) | Type     | Description               |
| :-------- | :------- | :------------------------------ |
| `alias`      | `string` | **Required**. Your URL alias |


### Privated endpoints

For these endpoints admin authorization is required.

The authorization is of type basic with *****username***** and *****password*****. 

#### Create a custom short URL

It generates a custom short URL from a base URL provided as a parameter. For this short URL alias may be provided as parameter.

```http
  POST /admin/
```

| Parameter (body) | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `url` | `string` | **Required**. Your base URL to shorten |
| `alias` | `string` | Alias for your custom URL |

---

#### Delete a short URL


```http
  DELETE /admin/${alias}
```

# JWT Authentication

It's an auth app generating JWT access and refresh tokens. It also storages users and has token's black list.

# API

<code>POST /signup</code> - registration by email, name and password<br>
body:

```
{
    email: string;
    password: string
}
```

response:

```
{
    accessToken: string;
}
```

<br>

<code>POST /login</code> - login by email and password<br>
body:

```
{
    email: string;
    password: string
}
```

response:

```
{
    accessToken: string;
}
```

<br>

<code>POST /update_access</code> - renew access token and refresh token if it's expired<br>
It requires cookie (with credentials) and header <code>Authorization: Bearer <ACCESS_TOKEN></code><br>
response:

```
{
    accessToken: string;
}
```

<br>

<code>GET /private/user</code> - get current user<br>
It requires header <code>Authorization: Bearer <ACCESS_TOKEN></code><br>
response:

```
{
    id: string;
    email: string;
    name: string;
}
```

<br>

<code>GET /private/user/:id</code> - get user by ID<br>
It requires header <code>Authorization: Bearer <ACCESS_TOKEN></code><br>
response:

```
{
    id: string;
    email: string;
    name: string;
}
```

<br>

# How to run

There will be a description about building and starting

<img src='https://cdn.tlgrm.app/stickers/219/f9d/219f9db9-34b0-343d-96bf-b4dc161a205e/192/2.webp' alt='Успех!'>

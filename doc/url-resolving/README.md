# URL Resolving

Amatl provides a flexible URL resolver to access various types of resources. It supports the following URL schemes:

- `file://` — Refers to local filesystem resources. Paths without a scheme are also interpreted as local files.
- `http://` and `https://` — Used to access HTTP(S) resources.
- `stdin://` — Refers to data piped in via standard input (`stdin`).

These URL schemes can be used consistently across the application, including when specifying inputs for commands like `render`.

> ### 🔐 Basic authentication
>
> Amatl supports HTTP [Basic Authentication](https://en.wikipedia.org/wiki/Basic_access_authentication) for `http(s)://` URLs.
>
> To enable it, set the following environment variables:
>
> - `AMATL_HTTP_BASIC_AUTH_USERNAME`
> - `AMATL_HTTP_BASIC_AUTH_PASSWORD`
>
> When these variables are set, credentials are automatically applied to all HTTP(S) requests during URL resolution.

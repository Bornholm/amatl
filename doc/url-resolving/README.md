# URL Resolving

Amatl utilizes a URL resolver to access various processed resources. It effectively manages the following schemes:

- `file://` - Denoting local filesystem resources (paths without protocol scheme will also be interpreted as a file).
- `http://` and `https://` - Indicating HTTP(S) resources.
- `stdin://` - Data provided via `stdin`.

As a general guideline, you can incorporate these types of URLs in all links or file paths, including those passed to the `render` command.

> **Basic Authentication**
>
> Amatl supports HTTP [`Basic Auth`](https://en.wikipedia.org/wiki/Basic_access_authentication) for `http(s)://*` URLs by setting the environment variables `AMATL_HTTP_BASIC_AUTH_USERNAME` and `AMATL_HTTP_BASIC_AUTH_PASSWORD`.
>
> If set, these credentials will automatically be used with all resolved URLs.

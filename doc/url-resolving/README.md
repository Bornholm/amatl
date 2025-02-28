# URL Resolving

Amatl utilizes a URL resolver to access various processed resources. It effectively manages the following schemes:

- `file://` - Denoting local filesystem resources (including "naked" paths).
- `http://` and `https://` - Indicating HTTP(S) resources.
- `stdin://` - Data provided via `stdin`.

As a general guideline, you can incorporate these types of URLs in all links or file paths, including those passed to the `render` command.

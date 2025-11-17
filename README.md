# oci-func-fileserver

OCI Object Storage Static File Server Function.  
This Oracle Cloud Infrastructure (OCI) Function serves static files from Object Storage as a simple HTTP server.  
It uses the Oracle Functions SDK (fdk-go) and authenticates via **Resource Principal**, so no API keys are required.

## Features

- Serves static files from **multiple** OCI Object Storage buckets.
- Bucket is automatically selected based on **Host + Path** mapping defined in the `ROUTES` environment variable.
- Automatically appends `index.html` for root (`/`) or directory paths (e.g., `/path/`).
- Preserves object metadata headers (Content-Type, Cache-Control, ETag, etc.).
- Handles common errors like 404 (ObjectNotFound) and 500 (internal errors).
- Stateless and serverless, optimized for OCI Functions.

## Prerequisites

- [Go 1.18+](https://go.dev/doc/install) installed.
- OCI CLI configured with appropriate permissions (for deployment).
- One or more OCI Object Storage buckets with static files uploaded.
- OCI Functions setup: Ensure you have `fn` CLI installed and configured.  
  ([Oracle Functions Quickstart](https://docs.oracle.com/en-us/iaas/Content/Functions/Tasks/functionsquickstart.htm))
- Docker installed.

Permissions required:
- Resource Principal enabled on the function's compartment.
- Read access to the target buckets (`objectstorage_object_get`, `objectstorage_object_head`).

## Installation and Deployment

```console
$ git clone https://github.com/mattn/oci-func-fileserver
$ cd oci-func-fileserver
$ fn deploy --app <your-function-app-name>
```

## Configuration

### Bucket routing (required)

Set the environment variable `ROUTES` to a JSON object containing:

```
{
  "<domain>": {
    "<path-prefix>": "<bucket-name>"
  }
}
```

### Example

```
fn config function <your-function-app-name> oci-func-fileserver ROUTES '{
  "example1.jp": { "/": "static-example1" },
  "example2.jp": { "/": "static-example2" }
}'
```

Prefix matching is supported, so more complex routing is possible:

```json
{
  "example1.jp": {
    "/blog/": "blog-bucket",
    "/assets/": "assets-bucket",
    "/": "default-bucket"
  }
}
```

## Invoke the function

You can attach this function to an **API Gateway** and use it as a static file server:

1. Create an API Gateway.
2. Add a route that points to this function.
3. Deploy the API Gateway.
4. Access files via the gateway URL or custom domain.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ROUTES` | JSON mapping of **host → path prefix → bucket**. Required. |

There is **no `BUCKET_NAME`** anymore.  
All routing is controlled via `ROUTES`.

## Usage

- **URL Mapping:**  
  The function strips the leading `/` and uses it as the object name.  
  Directories auto-resolve to `index.html`.

  - `/` → `index.html`  
  - `/images/logo.png` → `images/logo.png`  
  - `/docs/` → `docs/index.html`

- **Headers:**  
  All relevant OCI object headers are forwarded (e.g., `Content-Type`, `Cache-Control`, `ETag`).

- **Error Responses:**  
  - `404`: Object not found  
  - `500`: Authentication, client creation, or retrieval errors  

## Limitations

- Changing routing requires updating the function's environment variables.
- No authentication/authorization on served files.
- Assumes UTF-8 paths.

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)

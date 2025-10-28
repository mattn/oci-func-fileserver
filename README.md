# oci-func-fileserver

OCI Object Storage Static File Server Function. This Oracle Cloud Infrastructure (OCI) Function serves static files from an Object Storage bucket as a simple HTTP server. It uses the Oracle Functions SDK (fdk-go) and authenticates via Resource Principal, making it suitable for deployment in OCI Functions without needing API keys.

## Features

- Serves files from a specified OCI Object Storage bucket (default: "static").
- Automatically appends `index.html` for root (`/`) or directory paths (e.g., `/path/`).
- Preserves object metadata headers (Content-Type, Cache-Control, ETag, etc.).
- Handles common errors like 404 (ObjectNotFound) and 500 (internal errors).
- Stateless and serverless, optimized for OCI Functions.

## Prerequisites

- [Go 1.18+](https://go.dev/doc/install) installed.
- OCI CLI configured with appropriate permissions (for deployment).
- An OCI Object Storage bucket with static files uploaded.
- OCI Functions setup: Ensure you have `fn` CLI installed and configured ([Oracle Functions Quickstart](https://docs.oracle.com/en-us/iaas/Content/Functions/Tasks/functionsquickstart.htm)).
- Docker installed.

Permissions required:
- Resource Principal enabled on the function's compartment.
- Read access to the target bucket (`objectstorage_object_get`, `objectstorage_object_head`).

## Installation and Deployment

```console
$ git clone https://github.com/mattn/oci-func-fileserver
$ cd oci-func-fileserver
$ fn deploy --app <your-function-app-name>
```

## Customization

Set environment variable for bucket (optional, defaults to "static")

```
fn config function <your-function-app-name> oci-func-fileserver BUCKET_NAME my-bucket
```

## **Invoke the function:**

Now you can use this function as a static file server. And you can add this function to API Gateway to serve files over HTTP.

1. **Create an API Gateway** in OCI Console or via CLI.

2. **Add a route** to your API Gateway that points to your function.

3. **Deploy the API Gateway**.

4. **Access your files** via the API Gateway endpoint.

If you set up everything correctly, you should be able to access your static files using the API Gateway URL. If your certificate should be properly configured to use HTTPS, you can set up a custom domain and SSL certificate in API Gateway.

## Environment Variables

| Variable       | Description                  | Default |
|----------------|------------------------------|---------|
| `BUCKET_NAME` | Name of the Object Storage bucket to serve from. | `static` |

## Usage

- **URL Mapping:** The function strips the leading `/` from the request URL and uses it as the object name. Directories auto-resolve to `index.html`.
  - `/` → `index.html` or Not Found
  - `/images/logo.png` → `images/logo.png`
  - `/docs/` → `docs/index.html`

- **Headers:** All relevant OCI object headers are forwarded (e.g., `Content-Type`, `Cache-Control`).

- **Error Responses:**
  - 404: Object not found.
  - 500: Authentication, client creation, or retrieval errors.

## Limitations

- Single bucket support.
- No authentication/authorization on served files (public bucket or use OCI auth).
- Assumes UTF-8 paths; handles basic URL trimming.

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)

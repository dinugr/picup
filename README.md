# Picup

Picup is a self-hosted simple image upload and gallery service written in Go. It stores image metadata in SQLite, keeps files on local storage, extracts EXIF metadata, and generates configured image variants through FFmpeg. The main purpose is to keep original image while generating the configured variants needed for content blog post and other publishing workflows.

Picup (currently) relies on external engine for its image workflow:

- ExifTool extracts and preserves image metadata during upload processing.
- FFmpeg generates the configured image variants, including resized and cropped versions.

Both dependencies are required for the corresponding runtime features.

## Features

- Authenticated web UI with JWT-backed session cookies.
- Upload supported image files with configurable size and file-type limits.
- Search and page through the gallery.
- Open image details and view generated variants and metadata.
- Generate `display`, `preview`, and custom variants from uploaded images.
- Permanently delete images.
- Review and revoke active sessions, including logout from all devices.
- Review upload limits and configured variants.

## Requirements

- Go `1.26.5` or a compatible Go 1.26 toolchain for building the server.
- Node.js 24 and npm when building the web UI from source.
- ExifTool available on `PATH`, or set `exiftool_path` to its executable path.
- FFmpeg available on `PATH` for the default image-processing engine.
- Writable data and temporary directories.

## Quick Start From Source

1. Install the requirements above.
2. Copy the example configuration:

   ```sh
   cp core/example.config.ini config.ini
   ```

3. Edit `config.ini` before the first start. Set a unique `admin_username`, `admin_password`, and `jwt_secret`. The example's `[engine.default]` block is commented out; uncomment it so the configuration defines the required default FFmpeg engine.
4. Build the web UI:

   ```sh
   cd webui
   npm ci
   npm run build
   cd ..
   ```

   The build is written to `build/webui`, so set `webui_dir = ./build/webui` in `config.ini` when using the built UI.
5. Start Picup from the repository root:

   ```sh
   go run ./core -config config.ini
   ```

6. Open `http://localhost:9906` and sign in with the credentials configured in step 3.

The same startup form works with a compiled binary:

```sh
go build -o picup ./core
./picup -config config.ini
```

Picup also accepts `-c config.ini` or a positional configuration path, for example `./picup config.ini`.

## Packaged Releases

Release archives are provided for Linux amd64 with musl, Linux amd64 with glibc, and Windows amd64. Extract an archive and run the binary from the extracted directory. The packaged UI is under `webui`; the accompanying `config.ini` should use:

```ini
webui_enabled = true
webui_dir = ./webui
```

The packaged binary uses the same command-line configuration options as a source build:

```sh
./picup-linux-musl-amd64 -config config.ini
```

For glibc-based Linux distributions, run:

```sh
./picup-linux-glibc-amd64 -config config.ini
```

On Windows, run the equivalent executable from PowerShell:

```powershell
.\picup-windows-amd64.exe -config config.ini
```

## Configuration

Picup reads an INI configuration file. The default path is `config.ini`. Supply another path with `-config PATH`, `-c PATH`, or as the first positional argument.

Start from [core/example.config.ini](core/example.config.ini), then review the settings below before starting the server.

### Server

| Key | Default | Description |
| --- | --- | --- |
| `data_dir` | `./data` | Root directory for stored original images and generated variants. |
| `temp_dir` | `./temp` | Temporary working directory used during processing. |
| `max_size` | `10485760` | Maximum upload size in bytes. The example is 10 MiB. |
| `file_types` | none | Comma-separated accepted extensions, such as `png,jpg,jpeg,webp`. |
| `dsn` | `sqlite://./picup.db` | Database source name. The initial release uses SQLite. |
| `exiftool_path` | `exiftool` | ExifTool executable name or absolute path. |
| `host` | `localhost` | HTTP bind host. The value `*` enables listening on all interfaces. |
| `port` | `9906` | HTTP listen port. |
| `base_url` | empty | Optional externally visible base URL. |
| `webui_dir` | `./webui` | Directory containing the built static web UI. Source builds use `./build/webui`; release archives use `./webui`. |
| `webui_enabled` | `false` | Serves the static UI from `webui_dir` when `true`; otherwise serves the embedded fallback page. |
| `assets_base_path` | `/assets/images` | URL prefix used to serve image variants. A trailing slash must be omitted. |
| `variants_sequence` | `preview` | Comma-separated variant names controlling list order in the frontend. See the Variants section for details. |

#### Authentication

Authentication settings are also read from `Server`:

| Key | Default | Description |
| --- | --- | --- |
| `admin_username` | `admin` | Administrator username. Replace it before first deployment. |
| `admin_password` | `admin` | Administrator password. Replace it before first deployment. |
| `jwt_secret` | autogenerated | Secret used to sign access tokens. Set a persistent, random value in deployments. |
| `jwt_expiration_seconds` | `3600` | Access-token and session lifetime in seconds. |

If `jwt_secret` is omitted, Picup generates a new secret each time it starts. This invalidates existing sessions, so production deployments should use a persistent secret.

### Variants

A variant section has the form `[variant.NAME]`:

```ini
[variant.preview]
type = webp
args = scale=900:72
```

| Key | Default | Description |
| --- | --- | --- |
| `name` | Section name | Output variant name. |
| `type` | Reserved variant default | Output file type, such as `webp`. |
| `engine` | `default` | Configured processing engine used for the variant. |
| `args` | Empty | Pipeline operations. The current supported operations are `scale` and `crop`. |

`NAME` is the configuration key in `[variant.NAME]`. It must contain only letters and numbers and be no more than 16 characters.

Picup registers reserved `master`, `display`, and `preview` variants. Additional variants can be added by creating `[variant.NAME]` sections. The following example creates a square 300x300 WebP thumbnail from each uploaded image:

```ini
[variant.thumbnail]
type = webp
engine = default
args = scale=300:80:lrshr,crop=300:c
```

| Setting | Effect |
| --- | --- |
| `thumbnail` | Custom variant key and default output name. |
| `type = webp` | Selects the generated file format. |
| `engine = default` | Selects the built-in FFmpeg engine. This line can be omitted because `default` is the fallback engine. |
| `scale=300:80:lrshr` | Scales the image to fit within 300 pixels while preserving its aspect ratio, using quality `80`. |
| `crop=300:c` | Center-crops the result to a 300x300 square. |

Custom variant keys must contain only letters and numbers and be no more than 16 characters. The supported pipeline operations are currently `scale` and `crop`; unsupported operations or invalid arguments cause configuration loading to fail. The custom key can be listed in `variants_sequence` to control its display order in the frontend:

```ini
[server]
variants_sequence = preview,thumbnail
```

`variants_sequence` is used only by the frontend for ordering variants. It does not control which variants are generated during upload.

### Engines

The default engine is FFmpeg-backed and is required by the initial configuration. Define it explicitly:

```ini
[engine.default]
name = default
provider = ffmpeg
path = ffmpeg
```

| Key | Default | Description |
| --- | --- | --- |
| `name` | Engine key | Display and lookup name for the engine. |
| `provider` | Empty | Processor implementation to use. The initial release provides `ffmpeg`. |
| `path` | Provider default | FFmpeg executable name or path used by the processor. |

`provider = ffmpeg` selects the built-in FFmpeg processor. FFmpeg must be installed and available under the configured `path`.

A minimal working server configuration should contain at least one engine, at least one accepted file type, and a valid variant configuration.

### Example Deployment Configuration

```ini
[server]
data_dir = ./data
temp_dir = ./temp
max_size = 10485760
file_types = png,jpg,jpeg,gif,webp,heic,bmp,tiff
dsn = sqlite://./picup.db
exiftool_path = exiftool
host = localhost
port = 9906
webui_dir = ./build/webui
webui_enabled = true
assets_base_path = /assets/images
variants_sequence = preview
jwt_secret = replace-with-a-long-random-secret
jwt_expiration_seconds = 3600
admin_username = replace-admin-name
admin_password = replace-admin-password

[variant.preview]
type = webp
args = scale=900:72

[engine.default]
name = default
provider = ffmpeg
path = ffmpeg
```

For a first deployment, configure credentials and a persistent `jwt_secret`; omitted credentials fall back to development defaults and omitted JWT secrets are regenerated on startup.

## Development

Run the Go tests:

```sh
go test ./...
```

Build and type-check the web UI:

```sh
cd webui
npm ci
npm run typecheck
npm run build
```

Run the web UI lint task:

```sh
npm run lint
```

## Troubleshooting

- **ExifTool or FFmpeg not found:** install the executable and verify it is on `PATH`, or set `exiftool_path` and the configured FFmpeg path explicitly.
- **Configuration error:** check that at least one engine is configured, the default engine block is enabled, and variant names and keys use the allowed format.
- **UI not served:** verify `webui_enabled = true`, set `webui_dir` to the built output directory, and start the server from the directory used by the relative paths.
- **Upload or processing failures:** verify the file extension is listed in `file_types`, the file is within `max_size`, and both data and temporary directories are writable.
- **Address already in use:** change `port` or stop the process currently listening on the configured port.

## License

[MIT License](LICENSE)

# Do this tasks

## GoKotlinGallery - gallery_api Checklist (PostgreSQL, Metadata, Validation)

**1. Database Integration (PostgreSQL and `sqlc`)**

- [x] Install PostgreSQL.
- [x] Create database for gallery API.
- [x] Create database user and grant privileges.
- [x] Create `migrations` directory.
- [x] Use `go-migrate` to create database schema migration files.
  - [x] Create table for images (`id`, `filename`, `description`, `metadata`, etc.).
- [x] Create `sqlc.yaml` configuration file.
- [x] Create `.sql` files for database interactions (e.g., `insert_image.sql`, `get_image.sql`).
- [x] Run `sqlc generate` to generate Go code.
- [x] Update `internal/storage` to use `sqlc` generated code.
- [x] Update `SaveImage`, `GetImage`, and `GetImageDetails` to use PostgreSQL.

**2. Metadata and Description Handling**

- [ ] Create Go struct in `internal/models` for image metadata and description.
- [ ] Update database schema with metadata and description columns.
- [ ] Use a Go library (e.g., `github.com/disintegration/imaging`) to extract metadata.
- [ ] Store metadata in the database.
- [ ] Add description field to upload request.
- [ ] Store description in the database.
- [ ] (Optional) Add renaming feature with new filename field in upload request.
  - [ ] Store renamed image in storage.
  - [ ] Store new filename in database.

**3. Configurable File Size Limit**

- [ ] Use environment variable or configuration file for max file size limit.
- [ ] Read file size limit during API startup.
- [ ] Pass file size limit to `NewHandlers` function.
- [ ] Use file size limit in upload handler.

**4. Image Type Validation**

- [ ] Use `http.DetectContentType` to detect MIME type of uploaded files.
- [ ] Validate MIME type against allowed image types (e.g., `image/jpeg`, `image/png`).
- [ ] Return error for invalid file types.

**Implementation Notes (Reminder)**

- [ ] Use `go-migrate` CLI for database migrations.
- [ ] Ensure `sqlc` queries match database schema.
- [ ] Use Go standard library or third-party libraries for file I/O, JSON, HTTP.
- [ ] Implement error handling and logging.
- [ ] Use environment variables for sensitive data.

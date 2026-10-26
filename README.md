# Tapestry source reference

This private repository contains the supplied Tapestry source snapshot for AI-assisted code lookup. It is a source reference copy, not a deployment package.

## Read business context first

- Notion knowledge base: https://app.notion.com/p/405ba0e365ea81d39147d738898bb017
- Source baseline: `SNAP-9af8ab3208cabf9d`; the user confirmed the supplied source as the online version.
- Business-document calibration baseline: B119. Whole-project semantic review remains incomplete.

## Locate code

The imported directory layout is unchanged. A Notion path such as `tapestry-micro-scrm/internal/app/service/online_warehouse/order.go` is the same repository-relative path here. Front-end projects remain under `front/`.

Read `CODE_INDEX.md` for the project map and lookup examples. Read actual file contents before drawing implementation conclusions; a search hit alone is not evidence of full behavior.

## Import changes

Hardcoded credential literals and test credentials found during import were replaced with `REDACTED_FOR_SOURCE_SHARING`. Original local files were not modified. Credential replacement preserves existing line counts. `SOURCE_IMPORT_MANIFEST.json` maps original and imported file hashes; `SOURCE_REDACTIONS.json` lists changed locations without credential values.

These replacements can prevent the imported copy from running. Runtime credentials and environment configuration must be supplied independently. Original Git history was not included in the supplied export; this repository starts with one import commit.

## AI access

An AI client needs an authenticated GitHub connection with read permission for this private repository. Uploading source does not grant every AI access automatically. Prefer reading at a fixed commit for reproducible answers. Access to the Notion pages must also be authorized separately.

# Feature Implementation Plan

This document tracks planned improvements for `go-markdown-confluence`. Each task can be checked off once the feature has been implemented and verified.
This list builds upon the gaps outlined in `MISSING_FEATURES.md`.

## Planned Features

- [ ] **Callouts** - convert callout/admonition blocks to Confluence panels or expandable macros.
    - Parse `> [!TYPE]` style blocks from Markdown.
    - Map callout types (note, warning, info, etc.) to Confluence panel macros.
    - Support custom icons and colors where possible.
    - Add unit tests covering basic callout conversion.

- [ ] **Attachments / Image Upload** - add logic to upload images and other attachments to Confluence.
    - Use Confluence REST API to upload referenced images.
    - Replace local image paths with attachment references in the ADF output.
    - Avoid duplicate uploads by checking existing attachments.
    - Provide a CLI flag to toggle attachment uploading.

- [ ] **Comment Handling** - support stripping or converting Markdown comments (e.g., Obsidian style).
    - Detect HTML comments (`<!-- comment -->`) and Obsidian comments (`%% comment`).
    - Optionally remove comments from the output or convert them into Confluence notes.
    - Expose a configuration flag for enabling comment handling.

- [ ] **Folder Note & Folder Structure Sync** - keep Confluence page tree in sync with the file structure and support folder notes.
    - Treat `index.md` or `_index.md` files as the landing page for a folder.
    - Create or update Confluence pages to mirror nested directories.
    - Update links to match the generated page hierarchy.

- [ ] **Mermaid Diagrams** - render Mermaid diagrams to images for inclusion on pages.
    - Detect fenced code blocks with `mermaid` language tag.
    - Render diagrams to SVG or PNG using a headless renderer.
    - Upload generated images as attachments and embed them in the output.

- [ ] **Raw ADF Injection** - allow embedding raw Atlassian Document Format snippets from Markdown.
    - Recognize fenced blocks tagged `adf` and insert the contents as raw JSON.
    - Validate that the JSON is well-formed before inserting.

- [ ] **Setting Sources** - configure page metadata via YAML frontmatter or other sources.
    - Support keys such as `connie-title` and `connie-page-id` in the frontmatter.
    - Allow overriding settings via CLI flags or environment variables.

- [ ] **Wikilinks** - resolve `[[WikiLink]]` style links to Confluence pages.
    - Parse wiki-style links and map them to existing Confluence pages.
    - Fallback to normal Markdown links if the page cannot be found.

- [ ] **YAML Frontmatter Support** - parse keys like `connie-title`, `connie-page-id`, and `connie-frontmatter-to-publish`.
    - Extract metadata from YAML frontmatter and apply it during page creation.
    - Provide a setting to publish selected frontmatter keys to Confluence.

- [ ] **Additional CLI Flags** - provide more command-line options for customizing output locations and file naming.
    - Flags for output directory, naming templates, and dry-run mode.
    - Ensure flags integrate with existing directory processing logic.

- [ ] **Integration Tests** - expand tests to cover nested directory uploads and feature combinations.
    - Create test fixtures with multiple subdirectories and assets.
    - Verify that features like attachments and callouts work together.

## Updating This Checklist

Replace the `[ ]` with `[x]` once a feature is fully implemented and tested. Pull requests should update this file accordingly to reflect progress.

# Missing Features in go-markdown-confluence

This project implements a basic Markdown to Confluence converter written in Go. The converter currently supports headings, paragraphs, emphasis, links, images, code blocks, lists, blockquotes, task lists, and decision items. It also handles nested directories when uploading pages.

However, compared to the feature set documented on [markdown-confluence.com](https://markdown-confluence.com/), several capabilities are not yet implemented. These gaps are useful when planning further development.

## Features Not Implemented

- **Callouts**: Transformation of callout/admonition blocks into Confluence panels or expandable macros is not present.
- **Attachments / Image Upload**: Markdown images are converted to ADF image nodes but there is no logic to upload attachments to Confluence.
- **Comment Handling**: No support for stripping or converting Obsidian comments.
- **Folder Note & Folder Structure Sync**: While nested directories are converted, features like folder notes or keeping the page tree in sync with file structure are missing.
- **Mermaid Diagrams**: The converter does not render Mermaid diagrams to images for Confluence.
- **Raw ADF Injection**: There is no mechanism to include raw Atlassian Document Format snippets directly from Markdown.
- **Setting Sources**: Configurable sources for pages (e.g., using YAML frontmatter keys) are not supported.
- **Wikilinks**: Links in `[[WikiLink]]` style are not resolved to Confluence pages.
- **YAML Frontmatter Support**: Aside from minimal checks, YAML frontmatter keys like `connie-title`, `connie-page-id`, or `connie-frontmatter-to-publish` are ignored.

These omissions correspond to documentation in the following files from the official project docs:

- [`callouts.md`](https://github.com/markdown-confluence/docs-markdown-confluence/blob/main/src/features/callouts.md)
- [`comment-handling.md`](https://github.com/markdown-confluence/docs-markdown-confluence/blob/main/src/features/comment-handling.md)
- [`folder-note.md`](https://github.com/markdown-confluence/docs-markdown-confluence/blob/main/src/features/folder-note.md)
- [`folder-structure.md`](https://github.com/markdown-confluence/docs-markdown-confluence/blob/main/src/features/folder-structure.md)
- [`image-upload.md`](https://github.com/markdown-confluence/docs-markdown-confluence/blob/main/src/features/image-upload.md)
- [`mermaid.md`](https://github.com/markdown-confluence/docs-markdown-confluence/blob/main/src/features/mermaid.md)
- [`raw-adf.md`](https://github.com/markdown-confluence/docs-markdown-confluence/blob/main/src/features/raw-adf.md)
- [`setting-sources.md`](https://github.com/markdown-confluence/docs-markdown-confluence/blob/main/src/features/setting-sources.md)
- [`wikilinks.md`](https://github.com/markdown-confluence/docs-markdown-confluence/blob/main/src/features/wikilinks.md)
- [`yaml-frontmatter.md`](https://github.com/markdown-confluence/docs-markdown-confluence/blob/main/src/features/yaml-frontmatter.md)

## Next Steps

Adding support for these features would bring this Go implementation closer to the functionality described on markdown-confluence.com. Each feature above can be implemented incrementally, using the reference documentation as guidance.

# Advanced Markdown Features

## Footnotes

Here is some text with a footnote reference[^1].

[^1]: This is the footnote content.

Multiple footnotes[^2][^3] can be used.

[^2]: Second footnote.
[^3]: Third footnote with multiple paragraphs.

    Indented paragraphs are part of the footnote.

    ```
    code in footnotes
    ```

## Definition Lists

Term 1
: Definition 1

Term 2
: Definition 2a
: Definition 2b

## Emoji (GitHub Supports)

:smile: :rocket: :tada: :octocat:

## Mention Users

@username can be used to mention users on GitHub.

## Issue References

GitHub automatically links references to issues and pull requests: #123

## Automatic linking for URLs

https://github.com/ will be automatically converted to a link.

## Collapsed Sections

<details>
<summary>Click to expand!</summary>

### Hidden Markdown content

This content is hidden until the user clicks to expand it.

- List item 1
- List item 2

</details>

## HTML in Markdown

<div align="center">
  <h2>Centered heading</h2>
  <p>This paragraph is centered using HTML.</p>
</div>

<kbd>Ctrl</kbd>+<kbd>Alt</kbd>+<kbd>Del</kbd>

## Subscript and Superscript

H<sub>2</sub>O

E = mc<sup>2</sup>
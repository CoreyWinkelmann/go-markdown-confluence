# Diagrams and Special Markdown Features

## Mermaid Diagrams

GitHub supports Mermaid diagrams in markdown:

```mermaid
graph TD
    A[Start] --> B{Is it?}
    B -->|Yes| C[OK]
    C --> D[Rethink]
    D --> B
    B ---->|No| E[End]
```

### Sequence Diagram

```mermaid
sequenceDiagram
    participant Alice
    participant Bob
    Alice->>John: Hello John, how are you?
    loop Healthcheck
        John->>John: Fight against hypochondria
    end
    Note right of John: Rational thoughts <br/>prevail!
    John-->>Alice: Great!
    John->>Bob: How about you?
    Bob-->>John: Jolly good!
```

### Class Diagram

```mermaid
classDiagram
    Class01 <|-- AveryLongClass : Cool
    Class03 *-- Class04
    Class05 o-- Class06
    Class07 .. Class08
    Class09 --> C2 : Where am I?
    Class09 --* C3
    Class09 --|> Class07
    Class07 : equals()
    Class07 : Object[] elementData
    Class01 : size()
    Class01 : int chimp
    Class01 : int gorilla
    Class08 <--> C2: Cool label
```

## Math Expressions

GitHub also supports mathematical expressions using LaTeX syntax:

When $a \ne 0$, there are two solutions to $ax^2 + bx + c = 0$ and they are
$$ x = {-b \pm \sqrt{b^2-4ac} \over 2a} $$

Inline math: $E=mc^2$

## GitHub-specific Features

### Internal Linking to Headers

[Link to the Math Expressions section](#math-expressions)

### Relative Links

[Link to the basic.md file](basic.md)
[Link to the code examples](code-and-links.md#code)

### Keyboard Shortcuts

<kbd>Ctrl</kbd>+<kbd>C</kbd>: Copy
<kbd>Ctrl</kbd>+<kbd>V</kbd>: Paste
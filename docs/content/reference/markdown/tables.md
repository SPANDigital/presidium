---
title: Tables
weight: "4"
---

# In Markdown

Markdown provides a simple syntax for creating tables using hyphens and horizontal bars / pipes.

| first heading  | second heading | third heading  |
|----------------|----------------|----------------|
| row 1 column 1 | row 1 column 2 | row 1 column 3 |
| row 2 column 1 | row 2 column 2 | row 2 column 3 |
| row 3 column 1 | row 3 column 2 | row 3 column 3 |
| row 4 column 1 | row 4 column 2 | row 4 column 3 |

```md
| first heading  | second heading | third heading  |
|----------------|----------------|----------------|
| row 1 column 1 | row 1 column 2 | row 1 column 3 |
| row 2 column 1 | row 2 column 2 | row 2 column 3 |
| row 3 column 1 | row 3 column 2 | row 3 column 3 |
| row 4 column 1 | row 4 column 2 | row 4 column 3 |
```

This can be simplified by removing extra spaces, hyphens, and bars.

# In HTML

Markdown tables support single-line cells. For multi-line content, use an HTML table.

<table>
    <thead>
        <tr>
            <th>Heading 1</th>
            <th>Heading 2</th>
        </tr>
    </thead>
    <tbody>
        <tr>
        <td> Single line content. </td>
        <td>
            <ul>
                <li>Multi-</li>
                <li>line-</li>
                <li>content.</li>
            </ul>
        </td>
        </tr>
    </tbody>
</table>

```html
<table>
    <thead>
        <tr>
            <th>Heading 1</th>
            <th>Heading 2</th>
        </tr>
    </thead>
    <tbody>
        <tr>
        <td> Single line content. </td>
        <td>
            <ul>
                <li>Multi-</li>
                <li>line-</li>
                <li>content.</li>
            </ul>
        </td>
        </tr>
    </tbody>
</table>
```

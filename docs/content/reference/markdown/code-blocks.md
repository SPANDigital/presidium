---
title: Code Blocks
weight: "6"
---

To add code blocks to your content, enclose the code with three backticks. For syntax highlighting, set the language directly after the first set of backticks.
Alternatively, you can indent your code / machine output to treat it as preformatted text. For single line inline code, use a single backtick.

# Javascript

```js
var N = 32;
var buffer = new ArrayBuffer(N);
```

    ```js
    var N = 32;
    var buffer = new ArrayBuffer(N);
    ```

# Python

```py
my_array = [i for i in range(0, N)]
```

    ```py
    my_array = [i for i in range(0, N)]
    ```

# C

```c
int * my_func(int * in) {
    return in;
}
```

    ```c
    int * my_func(int * in) {
        return in;
    }
    ```

# Others

Github-flavoured Markdown supports many languages for code blocks. For a full list, go to [Github](https://github.com/github/linguist/blob/master/lib/linguist/languages.yml).

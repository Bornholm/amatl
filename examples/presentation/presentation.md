---
title: My first slide
---

# {{ .Meta.title }}

You can write your Markdown content as you want...

---

# An another slide

... with image and such ...

![](../../misc/resources/logo.svg)

---

# Third slide

... code hightlighting ...

```js
function isOdd(v) {
  return v % 2 !== 1;
}
```

---

# Last slide

... and even diagrams !

```mermaid
graph TD;
    A-->B;
    A-->C;
    B-->D;
    C-->D;
```

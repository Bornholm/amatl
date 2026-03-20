# MCP Server

Amatl includes a built-in [Model Context Protocol](https://modelcontextprotocol.io/) (MCP) server. It allows an AI assistant or a compatible editor to interact with the Markdown files in a workspace directory.

The server communicates over **stdin/stdout** (JSON-RPC 2.0 stdio transport).

## 🚀 Starting the server

```sh
amatl mcp serve --workspace ./docs
```

If `--workspace` is omitted, the current directory is used.

| Flag | Alias | Environment variable | Description |
|---|---|---|---|
| `--workspace` | `-w` | `AMATL_MCP_WORKSPACE` | Root directory exposed to the server. Access is strictly confined to this directory. |

## 🔧 Configuring your MCP client

Most MCP clients expect a configuration entry similar to:

```json
{
  "mcpServers": {
    "amatl": {
      "command": "amatl",
      "args": ["mcp", "serve", "--workspace", "/path/to/your/docs"]
    }
  }
}
```

Refer to your client's documentation for the exact configuration format.

## 🛠️ Available tools

| Tool | Description |
|---|---|
| `list_files` | List files in the workspace, with optional glob pattern filtering. |
| `table_of_contents` | Return the hierarchical table of contents of a Markdown file. |
| `list_sections` | List all headings in a file as a flat list with their selectors. |
| `read_section` | Read the full content of a section (heading + body). |
| `find_sections` | Recursively find all nodes matching a selector. |
| `update_section` | Replace the content of a section in a Markdown file. |

## 🔍 Selector syntax

Selectors are a CSS-like mini-language for targeting nodes in a Markdown document.

| Example | Matches |
|---|---|
| `h2#introduction` | `## Introduction` heading via its auto-generated ID |
| `h2:contains("Usage")` | `## Usage` heading matched by text |
| `code[lang="go"]` | Fenced code block with language `go` |
| `h2#api ~ table` | Table following the `## API` heading |
| `h2` | All level-2 headings |

IDs are automatically generated from heading text by the Markdown parser.

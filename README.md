# C2V2 - Developer Toolbox

C2V2 is a lightweight, efficient web-based toolbox designed for developers. It provides essential utilities like JSON formatting, Base64 encoding/decoding, and HTML processing with a focus on speed and ease of use.

## 🚀 Features

-   **JSON Tool:** Format, minify, and validate JSON. Supports conversion to Go Structs and YAML.
-   **HTML Tool:** Prettify, minify, escape, and unescape HTML. Features real-time client-side processing using `js-beautify`.
-   **Base64 Tool:** Robust encoding and decoding of text data.
-   **Multi-language Support:** Fully localized in English and Chinese.
-   **SEO Optimized:** Built-in Sitemap generation and JSON-LD schema support for better search engine visibility.
-   **Modern UI:** Clean, responsive interface built with Tailwind CSS, AlpineJS, and HTMX.

## 🛠️ Tech Stack

-   **Backend:** Go (Golang) with the Gin web framework.
-   **Frontend:** HTMX, AlpineJS, Tailwind CSS (via CDN).
-   **Libraries:** `js-beautify` (HTML), `prism.js` (Syntax Highlighting).

## 🏃 Getting Started

### Prerequisites

-   Go 1.23 or higher.

### Installation & Running

1.  Clone the repository:
    ```bash
    git clone <repository-url>
    cd c2v2
    ```

2.  Install dependencies:
    ```bash
    go mod tidy
    ```

3.  Run the application:
    ```bash
    go run cmd/server/main.go
    ```

4.  Open your browser and navigate to `http://localhost:8080`.

## 📁 Project Structure

```text
├── cmd/                # Entry points
│   └── server/         # Main server application
├── internal/           # Private application and library code
│   ├── app/            # Core application logic (router, handlers)
│   ├── pkg/            # Shared packages (i18n, render)
│   └── tools/          # Specific tool implementations
├── locales/            # Translation files (en.json, zh.json)
├── static/             # Static assets (images, scripts, styles)
├── templates/          # HTML templates (Go html/template)
│   ├── pages/          # Full page templates
│   ├── partials/       # Reusable template components
└── go.mod              # Go module definition
```

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details (if applicable).

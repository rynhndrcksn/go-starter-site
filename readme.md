# Go Starter Site - [![Go Report Card](https://goreportcard.com/badge/github.com/rynhndrcksn/go-starter-site)](https://goreportcard.com/report/github.com/rynhndrcksn/go-starter-site)

This repository exists to act as a good starting point when developing a new website using Go.
When I first started learning how to build websites with Go, I followed Alex Edward's advice in his book [Let's Go](https://lets-go.alexedwards.net).
While the book is fantastic, and it's easy to recommend people read, I noticed that there's a lot of boilerplate that Go sites need.

There are some solutions I found out there that people had done to remedy this problem themselves.
However, I wasn't a big fan of them because they included a lot of extra stuff that not everyone needed or wanted (Templ, Tailwind, and/or React).
Thus, I had the idea to make a useful, yet generic enough, Go site template that aims to have a lot of the boilerplate already done.
This way, one can easily clone the repository, swap out the module name, and hit the ground running.

## Getting Started

While on desktop, click the "Use this template" button located near the top right of the page and follow the steps provided.

Change the module name from `github.com/rynhndrcksn/go-starter-site` to match your repository.

> [!NOTE]
> Don't forget to update all the import paths throughout the project as well!

## Dependencies

In an effort to keep third party dependencies to a minimum, only the following ones have been added:

1. https://github.com/alexedwards/scs | Session management
2. https://github.com/julienschmidt/httprouter | Http router

There are some development related dependencies that I recommend installing to your local machine:

1. https://github.com/air-verse/air | Live reloading of Go binary (`make dev/web` makes use of this)
2. https://staticcheck.dev | Linter for Go (`make audit` makes use of this).
   1. This can be installed by running `go install honnef.co/go/tools/cmd/staticcheck@latest` (as of August 2024).

## Contributing

Contributions are welcome!

Please refer to [contributing.md](contributing.md) for more information.

## Reference

If you're curious why the project is structured this way, look here!

- `cmd/` contains the entry points for the application.
    - `web/` contains the server side logic for the website (routing, handlers, etc.).
- `internal/` contains things like validators, models, sending emails, etc.
    - `vcs/` contains logic for figuring out what version of the site is running.
- `ui/` contains everything relating to HTML templates and site assets (css, js, and images).
    - `html/` contains all the templates for constructing the website.
        - `components/` contains components to embed into partials and/or pages.
        - `pages/` contains full page templates.
        - `partials/` contains partial templates for embedding into other templates.
    - `static/` contains all the assets for the site.
        - `css/` contains all the stylesheets for the site.
        - `js/` contains all the scripts for the site.

## License

This is released under the MIT license, which can be viewed [here](LICENSE).

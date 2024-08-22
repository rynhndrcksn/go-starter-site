# Go Starter Site

This repository is geared towards being a good leaping off point for building new sites using Go.

## Getting Started

Clone this repository:

```shell
git clone https://github.com/rynhndrcksn/go-starter-site .
```

Remove the existing `.git/` directory:

```shell
rm -rf .git/
```

Initialize your own git repository:

```shell
git init
```

Make the site your own!

## Dependencies

In an effort to keep third party dependencies to a minimum, only the following ones have been added:

1. https://github.com/alexedwards/scs | Session management
2. https://github.com/julienschmidt/httprouter | Http router

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

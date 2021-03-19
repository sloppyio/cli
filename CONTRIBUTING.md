# Contributing to sloppy.io CLI

A great way to start contributing is to send a detailed report when you encounter an issue.

When reporting issues, always include:
* the output of `sloppy version`
* the output of `sloppy --debug [COMMAND]`

Also include the steps necessary to reproduce the issue if possible. This will help us review and fix your issue quickly.
Please consider removing sensitive data from your output before posting it.

When actually working on issues, please:

1. Fork this project

2. Setup a new branch to work in

3. Always run `make fmt` on your code before committing it

4. Even if it's a minor change, you should always write tests

5. Run `make test` or `scripts/test.sh`

6. Try to write good commit messages for each change. Use the [angularjs commit message guide](https://gist.github.com/stephenparish/9941e89d80e2bc58a153)

7. Push the commits to your fork and submit a pull request

## Development environment

### Building

To build locally you need to have go 1.6 or later. If not, please upgrade.

```sh
make local
```

You can also build different pre-releases of the sloppy CLI.
```sh
scripts/make.sh build beta      # For single builds
scripts/make.sh build release   # w/o pre-release
scripts/make.sh cross rc.[0-9]  # For cross compiling
```

To build within a docker container:
```sh
make beta
```
The binaries created are stored in `./bundles/${VERSION}-${PRERELEASE}/`.

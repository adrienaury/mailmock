# Contributing

If you have a feature request, suggestion, or bug report, please open a [new issue](https://github.com/adrienaury/mailmock/issues/new).

To submit patches, please fork the project and send a pull request. To be merged, please make sure you follow all items addressed in this guide.

If you need help/more information, you can [contact me](mailto:adrien.aury@gmail.com) by mail.

## How to start

1. Fork the repository on GitHub
2. Pull requests must be sent from a **new** hotfix/feature branch, not from master.
3. Make your modifications, follow the contributing guidelines (CONTRIBUTING.md)
4. Commit small logical changes, each with a descriptive commit message. Please don't mix unrelated changes in a single commit.

## Submit your changes

1. Push your changes to a topic branch in your fork of the repository.
2. Submit a pull request to the original repository. Describe your changes as short as possible, but as detailed as needed for others to get an overview of your modifications.

## Style guidelines

This project tries to comply as much as possible with the "[Effective GO](https://golang.org/doc/effective_go.html)" style guidelines. Also [this article](https://dave.cheney.net/practical-go/presentations/qcon-china.html#_consider_fewer_larger_packages) is a big inspiration.

Source code must be formatted with `gofmt` and verified with `go_vet` and `golint`.

Use `goimports` to remove unused imports.

## Cyclomatic complexity

Use [gocyclo](https://github.com/fzipp/gocyclo) to compute cyclomatic complexity. Cyclomatic complexity of functions/methods must be below 15.

## Tests

Adding go tests is highly appreciated ! At least for most common use cases.

## License

Every file must begin with a list of copyright notice. If you modified an existing file, you can (no obligation) add a line for your contribution. Use this template :

```go
// Copyright (C) 2019  Your Name
```

Then the standard GPL v3 list of conditions and disclaimer must be inserted (for new files, copy/paste from an existing file).

## Documentation

- Package comment must be present at least once for each package.
- Every exported (capitalized) name must have a doc comment.
- README.md must be updated to reflect changes
- CHANGELOG.md must be updated to list changes

## Changlog

This project keeps a changelog in the CHANGELOG.md file, the format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/). Update it as often as possible.

## Commits

- Don't mix unrelated changes into a single commit.
- No need to capitalize commit messages.
- Limit the first line to 72 characters or less.
- If the commit is about something that is unfinished, start the message with `wip :`

## Versioning

This project follows strictly [semantic versioning](https://semver.org/) for tags and version numbers.

## Branching model

This project uses [this successful branching model](https://nvie.com/posts/a-successful-git-branching-model/).

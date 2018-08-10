# Contributing Guidelines

Contribuations are always welcome, although please search issues and pull requests before adding something new to avoid duplicating
efforts and conversations.



### Rules

All Contributions are more than welcome but must be submitted using **Pull requests** with following guidelines.

* Include tests for any new functionality
* Don't delete any existing tests unless absolutely necessary
* Don't change too much in one pull request, small changes pull requests are quickly reviewed and merged and save everybody time.
* Separate out different changes into different pull requests
* Please mention in your pull requests if your changes are not backwards compatiabile so that we can increase appropriate vesion number.

## Code Style

* This repository uses go fmt to maintain code style and consistency.
* Test folder names should be lowercase
* Test files should also be lowercase and end with `_test.go`
* We use vendoring for dependency package support.
* All files and folders should be snake_case.

### Creating Issues

Submit issues to the issue tracker on the appropriate repository for suggestions, recommendations, and bugs.
Please note: If it's an issue that's urgent / you feel you can fix yourself, please feel free to make some changes and submit a pull request. We'd love to see your contributions.


### Pull Requests

* Create a new local branch for your work in your fork.
* As early as possible, create a pull request in the appropriate repository. Make sure you give enough information in the pull request description, and add the label `in-progress` with any other appropriate label.
* Once any conflicts have been fixed and you're ready for your code to be reviewed, remove the `in-progress` label and add `merge-ready`.
* Get a code review. Two of these must be core maintainers of GoNSX.
* You need two :shipit: :shipit: left as comments on the pull request.
* One of the core maintainers will merge the changes and apply appropriate versioning to release (see below).

### Discussion

For discussion of issues and general project talk, head over to [https://gitter.im/gonsx/Lobby](https://gitter.im/gonsx/Lobby).

### Releases

Declaring formal releases remains the prerogative of the project maintainer. We use Semantic Versioning in the releases.

### Changes to this arrangement

This document may also be subject to pull-requests or changes by contributors where you believe you have something valuable to add or change.

## Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

- (a) The contribution was created in whole or in part by me and I have the right to
  submit it under the open source license indicated in the file; or

- (b) The contribution is based upon previous work that, to the best of my knowledge, is
  covered under an appropriate open source license and I have the right under that license
  to submit that work with modifications, whether created in whole or in part by me, under
  the same open source license (unless I am permitted to submit under a different
  license), as indicated in the file; or

- (c) The contribution was provided directly to me by some other person who certified
  (a), (b) or (c) and I have not modified it.

- (d) I understand and agree that this project and the contribution are public and that a
  record of the contribution (including all personal information I submit with it,
  including my sign-off) is maintained indefinitely and may be redistributed consistent
  with this project or the open source license(s) involved.
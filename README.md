<div align="center">
<h1> Ginx </h1>

![star]
[![Downloads][downloads-badge]][releases]
![version]

</div>

**ginx** is a lightweight CLI tool designed to monitor changes in remote Git repositories and execute custom commands when updates occur. It is ideal for automating tasks, deploying applications, or triggering workflows whenever the target repository is updated.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Example](#example)
- [License](#license)

## Features

- **Monitor Remote Repositories**: Watch specific branches for updates.
- **Custom Commands**: Execute any command when changes are detected.
- **Configurable Intervals**: Set the polling frequency for repository checks.
- **Flexible Logging**: Adjust log levels for better visibility or minimal noise.
- **Version Display**: Check the current version of the tool.
- **Git Repo Work Directory Sandboxing**: It is possible to open your IDE as the ginx command and just work there until you close it. Each time you run this command, it will instantiate another working directory.

## Installation

To install `ytui`, follow the instructions for your operating system.
Ensure that you have the required dependencies installed.

1. **Install binary**

   ginx runs on most major platforms. If your platform isn't listed below,
   please [open an issue][issues].

   Please note that binaries are available on the release pages, you can extract the archives for your
   platform and manually install it.

## Usage

See [Documentations](docs/ginx.md).

## Example

Run a colmena local update to rebuild your NixOS immediately:

```bash
ginx --source https://github.com/didactiklabs/nixbook -b main --now -- colmena apply-local --sudo
```

Watch and trigger update on changes:

```bash
ginx --source https://github.com/didactiklabs/nixOs-server -b main -n 60 -- colmena apply-local --sudo
```

Open your IDE in a sandbox of your repository:

```bash
ginx --source https://github.com/didactiklabs/ginx -b main --now -- vim .
```

## License

![licence]

`ginx` is open-source and available under the [LICENCE](LICENSE).

For more detailed usage, you can always use `ginx --help`.

[licence]: https://img.shields.io/github/license/didactiklabs/ginx
[downloads-badge]: https://img.shields.io/github/downloads/didactiklabs/ginx/total?logo=github&logoColor=white&style=flat-square
[releases]: https://github.com/didactiklabs/ginx/releases
[star]: https://img.shields.io/github/stars/didactiklabs/ginx
[version]: https://img.shields.io/github/v/release/didactiklabs/ginx

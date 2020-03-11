# Spandx::Index

This repository keeps track of the software license
for each package dependency.

To find the appropriate datafile for a package, follow these steps:

1. Compute a SHA 256 of the package name.
2. Take the first to characters of the SHA256.
3. Expore the directory with the name that matches the first two characters of the SHA256.
4. Find the data file with a name that matches the package manager. E.g. `nuget` for packages sourced from `api.nuget.org`.

## Installation

You will need a modern ruby installed and preferably a \*nix based environment.

## Usage

To update the index with the newest published packages.

```bash
$ ./bin/cibuild
```

## Development

Visit [spandx](https://github.com/mokhan/spandx) to contribute.

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/mokhan/spandx-index.

If you recognize a missing package, or a mistake please submit a pull request to correct the data in the index.

## Copyright

Copyright (c) 2020 mo khan. See [MIT License](LICENSE.txt) for further details.

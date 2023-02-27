# opmsk

[![build](https://github.com/winebarrel/opmsk/actions/workflows/build.yml/badge.svg)](https://github.com/winebarrel/opmsk/actions/workflows/build.yml)

Mask 1Password CLI password using ANSI color.

cf. https://developer.1password.com/docs/cli/get-started/

## Installation

```sh
brew install winebarrel/opmsk/opmsk
```

## Usage

```sh
op item get <item> --format json | opmsk
```

![221649867-8b035c54-3de8-4a60-94e1-e2e3f20ce9a3](https://user-images.githubusercontent.com/117768/221654762-23f63304-5182-4c60-8ca0-b66f2c38721e.png)


Masked values are displayed when selected with the mouse.

![](https://user-images.githubusercontent.com/117768/221647712-5084642d-5fcb-4286-9c6d-b5ac2871e08c.png)

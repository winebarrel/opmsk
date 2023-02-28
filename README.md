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

![](https://user-images.githubusercontent.com/117768/221726898-d0a2f733-c856-4207-8e99-f453e99a3368.png)

Masked values are displayed when selected with the mouse.

![](https://user-images.githubusercontent.com/117768/221726922-6de065ad-33bd-4782-b9d5-d1bdccddc8f3.png)

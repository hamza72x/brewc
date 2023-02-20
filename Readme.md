## brew concurrent downloader and installer, around 50%/60% faster with huge dependencies, for example: ffmpeg

## Install

```sh
go install github.com/hamza72x/brewc@latest
```

## Usage

```
Usage:
  brewc [flags]
  brewc [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  install     install a formula
  reinstall   reinstall a formula
  uninstall   uninstall a formula

Flags:
  -h, --help   help for brewc

Use "brewc [command] --help" for more information about a command.
➜  Downloads
```

## Example

```
brewc install ffmpeg wget curl git
```

## Compare

- installing ffmpeg took around `2:35` mintues with brewc and with brew it took around `4:15` minutes

[![YouTube](https://img.youtube.com/vi/VVfNutjzF64/0.jpg)](https://youtu.be/VVfNutjzF64)

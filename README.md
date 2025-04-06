# ZMK-Flassher

A simple go application to flash ZMK firmware to a ZMK powered split keyboards.

## Installation

```bash
go install github.com/new-er/zmk-flasher@latest
```

## Usage

To flash a firmware run the following command:
```bash
zmk-flasher -l <left_firmware> -r <right_firmware>
```
This lets you mount the left and right keyboard halves interactively.

You can also flash a single firmware file to both halves (Glove 80):
```bash
zmk-flasher -a <firmware>
```

To find more about the usage of the application, run the following command:
```bash
zmk-flasher --help
```

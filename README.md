# ZMK-Flassher

A small go cli application to simplify flashing ZMK firmware to ZMK powered split keyboards.
It reduces friction by letting you specify the firmware files of the right and left halve and mounting the keyboard bootloaders interactively.
Afterwards it will copy the firmware to the keyboard halves.

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
Afterwards the application will flash the firmware to the left and right halves.

You can also flash a single firmware file to both halves (Glove 80):
```bash
zmk-flasher -a <firmware>
```

To find more about the usage of the application, run the following command:
```bash
zmk-flasher --help
```

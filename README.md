# mfm-go
mfm-go is Golang implementation of the [MFM](https://github.com/DaveAckley/MFM) for those who prefer the terminal.
The model is a watered-down version of MFMS optimized for trying things out.
It does not aim to accurately model a T2Tile (hardware tile).
The releases tab features interactive terminal binaries implemented using the [termbox-go](https://github.com/nsf/termbox-go) library.
So far this looks and functions a bit like vim.
Mouse support might be added if it's requested.

## Key Bindings
|Key   |Binding                  |
|------|-------------------------|
|CTRL+C|Quit                     |
|CTRL+X|Reset grid               |
|SPACE |Pause/resume             |
|ESC   |Hide cursor              |
|hkmu←→↑↓|Move cursor            |
|HKUM^$|Move cursor (fast)       |
|.     |Step                     |
|x     |Clear                    |
|v     |Visual mode (select area)|
|r     |Place Res/Replace mode   |
|X     |Place Fork               |
|D     |Place DReg               |

## Installing

Check the releases tab to download the latest binary for your platform. If you want to install from source, you'll need the Golang 1.11 toolchain. After that, running `install.sh` from the project root directory should work.

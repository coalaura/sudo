# sudo <sup><kbd>For Windows</kbd></sup>

Ever wanted to `sudo` on Windows? Now you can! This tool, simply called sudo, brings the familiar `sudo` functionality to Windows, allowing you to run any command with elevated privileges—no more fiddling with right-click menus or navigating through the UAC dialog manually. Just type `sudo <command>` in your terminal, and it takes care of the rest.

## Features

- **Runs Commands as Admin:** Just like `sudo` on Linux, this tool checks if the current terminal session is elevated. If not, it prompts for admin privileges and runs your command in the specified directory.
- **Session Persistence:** Already elevated? `sudo` detects it and skips the re-elevation step, running your command directly.
- **Seamless Directory Management:** `sudo` ensures your commands execute in the current working directory, whether you’re elevated or not.
- **`su` Mode:** For those who just want an admin shell, running `sudo su` gives you an elevated terminal window without extra commands.

## Usage
```bash
sudo <command> [arguments...]
```

Or to just get an admin shell:
```bash
sudo su
```
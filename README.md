# Clippy

Cliipy is an effective tool to help work with your clipboard, it works by storing clipboard history, and giving you the ability to do what you want with it.

You can search through your clipboard history, list and start the program as a background process.

## Installation
```shell
chmod +x ./scripts 
cd scripts
./install.sh
```

Clippy should be running in the background.

Then you should exposed to these command
Available Commands:
  list        list your clipboard history
  search      Search your clipboard history
  start       Start program in background
  
## Example
```shell
clippy list --limit=5
clippy search define
clippy search define --limit=30
```

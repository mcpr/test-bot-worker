# `MCPR-Test-Bot` Worker
[MCPR](https://mcpr.io) Plugin Test Bot Worker

## What is this???
`MCPR-Test-Bot` Worker is the worker component which does all of the heavy lifting of the `MCPR-Test-Bot`. It runs some basic compatibility tests on plugins uploaded to [MCPR](https://mcpr.io). 

## Test Steps
MCPR-Test-Bot will run the following steps on each plugin jar uploaded. 

1. Download the plugin jar
2. Unzip the jar and check for a `plugin.yml` file. 
3. Do the following for each claimed compatible version of Minecraft plus the latest
   1. Build (or download private pre-built) Spigot and CraftBukkit version
   2. Install the plugin
   3. Start server and watch for errors
   4. Once the server has fully started without error, stop it and send the results back to [MCPR](https://mcpr.io)
   5. If the server crashes or times out, send the results back to [MCPR](https://mcpr.io)

If you think there should be more tests done, feel free to [open an issue](https://github.com/mcpr/test-bot-worker/issues/new) or PR!

## Usage

```
NAME:
   test-bot-worker - The MCPR-Test-Bot Worker

USAGE:
   test-bot-worker [global options] command [command options] [arguments...]

VERSION:
   0.0.1

DESCRIPTION:
   The MCPR-Test-Bot Worker

COMMANDS:
     test, t  Run tests on plugin - `test-bot-worker test [pluginID] [version]`
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
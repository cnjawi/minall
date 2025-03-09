# minall

use LLM directly from terminal (useful for copying and pasting markdown/latex contents), written in Go.

API referrences are [here](https://help.aliyun.com/zh/model-studio/developer-reference/use-qwen-by-calling-api).

help page:
```
PS> minall help
use LLM directly from terminal.

Usage: minall [command] [options] [...message]

Commands:
  chat     -- start a new chat session
  pipe     -- read from stdin
  trans    -- translate stdin, require translator model
  help     -- print this message

Use "minall [command] -h" to see options of each command.

if command is not matched, all
arguments will be parsed as a message.

Config file is stored in:
  %AppData%\Roaming\minall\

Models defined in config file:
  qwm    -- qwen-max-latest
  qwq    -- qwq-plus-latest
  mtt    -- qwen-mt-turbo
  mtp    -- qwen-mt-plus
  dsv3   -- deepseek-chat
  dsr1   -- deepseek-reasoner
  qwt    -- qwen-turbo-latest
  qwp    -- qwen-plus-latest
```

models are defined my user in the config file
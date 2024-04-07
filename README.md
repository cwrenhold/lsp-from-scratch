# Introduction

This is a repo of me following [this video](https://www.youtube.com/watch?v=YsdlcQoHqPY) by [TJ DeVries](https://github.com/tjdevries), which can also be found [here](https://github.com/tjdevries/educationalsp), if you want to see how it should've been done rather than what I put together.

# Setup

To use this repo, I've included a devcontainer configuration, but you'll also need to have neovim and it would be helpful to have neovim inside of your devcontainer (see [here](https://github.com/cwrenhold/devc-nvim-commands) for one way to do this). From here, you can use the following steps:

- Exec into the dev container and run `go build main.go` in the repo
- Run the `create_nvim_config.sh` script to create a basic neovim config lua in your nvim setup
- Add anything additional to load the generated lua script into your neovim config from this starting point
- Launch neovim, and load up a Markdown file, this will trigger the LSP
- In the lsp repository, a log file called "lsp.txt" should be created

# Notes

- This is only configured to work with neovim for now, but could be extended to support VS Code or other consumers of LSPs later


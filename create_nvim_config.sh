FILE_TARGET="$HOME/.config/nvim/load_testing_lsp.lua"

# Create the target file, if it doesn't already exist
if [ ! -f "$FILE_TARGET" ]; then
	touch "$FILE_TARGET"
fi

# Write the contents of the file
cat <<EOT > "$FILE_TARGET"
---@diagnostic disable-next-line: missing-fields
local client = vim.lsp.start_client {
    name = "educationallsp",
    cmd = { "/workspaces/lsp-from-scratch/main" },
    -- on_attach = require('lspconfig').on_attach
}

if not client then
    vim.notify "The LSP server could not be started"
    return
end

vim.api.nvim_create_autocmd("FileType", {
    pattern = "markdown",
    callback = function()
        vim.lsp.buf_attach_client(0, client)
    end
})
EOT

echo "Successfully created $FILE_TARGET"

call plug#begin('~/.config/nvim/plugged')

Plug 'tpope/vim-sensible'
Plug 'github/copilot.vim'
Plug 'neoclide/coc.nvim', { 'branch': 'release' }
Plug 'ryanoasis/vim-devicons'
Plug 'tpope/vim-fugitive'
Plug 'vbundles/nerdtree'
Plug 'feline-nvim/feline.nvim'

-- Go Plugins --
Plug 'fatih/vim-go', { do := { ':GoInstallBinaries' } }

-- Web Dev Plugins --
Plug 'jelera/vim-javascript-syntax'
Plug 'mxw/vim-jsx'
Plug 'leafgarland/typescript-vim'
Plug 'HerringtonDarkholme/yats.vim'

-- Python Plugins --
Plug 'davidhalter/jedi-vim'
Plug 'raimon49/requirements.txt.vim', {'for': 'requirements.txt'}

call plug#end()

require('feline').setup()

-- make an http request to the server to get the config colorscheme --
local https = require('ssl.https')
local response, code, headers, status = http.request{
  url = 'http://localhost:9999/config/colors',
  method = 'GET',
  headers = {
    ['User-Agent'] = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36'
  }
}

local colors = {}

-- require json to decode the response --
local json = require('json')
colors = json.decode(response)

-- define the colorscheme --




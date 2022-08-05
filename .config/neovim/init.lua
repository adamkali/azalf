call plug#begin('~/.config/nvim/plugged')

Plug 'tpope/vim-sensible'
Plug 'github/copilot.vim'
Plug 'neoclide/coc.nvim', { 'branch': 'release' }
Plug 'ryanoasis/vim-devicons'
Plug 'tpope/vim-fugitive'

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

filetype plugin indent on


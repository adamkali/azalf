# Stream Plan

I want to get my entire flavor planned out so that when I get to the actual install to wipe bloatdows, it is as painless as possible.

## Packages I will get in pacman:

    pacman -S go rust node-js python3 qtile git firefox thundebird obs-studio xorg steam gimp blender amd-ucode

## YAY

    yay -S --no-confirm rocm-hip-runtime hip-runtime-amd rocm-opencl-runtime vs-code
    spotify ulauncher xlockmore cava wired

## eww

    git clone https://github.com/elkowar/eww

# Neovim 
    pip install neovim
    https://github.com/vim-scripts/vim-plug.git


# configs
``` sh

#!/bin/sh

$CONN;
$BYTES;

[[ ping -c 1 google.com &> /dev/null ]] {
    [ pacman -S go rust node-js python3 qtile git firefox thundebird obs-studio xorg steam gimp blender amd-ucode ] 
} || {
    echo "something went wrong in Pacman installation list."; exit;
}

] 

time [ git clone https://]

| echo "Those commands did not go so well..." 



```

## In vm

- [ ] Get a install config for a general layout so I only have to change what ths storage set up is.
- [ ] Configure Qtile with my wishlist and add that to a .configs/q/config.py
- [ ] Configure Qtile's bar
- [ ] get the picom fork with animations to work andd add its to config to a ./config/picom
- [ ] makesure my neovim is set up and add it into ./config/nvim/init.vim
- [ ] convirt init.vim to init.lua
- [ ] get a dashboard with eww and add its configs to ./config/eww/
- [ ] verify i can get Rust, Go, MariaDB, node, python to run through the install I made.
- [ ] write a shell script that will automate the install
- [ ] 
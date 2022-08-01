#!/bin/bash

echo "Welcom to the installation script";
# make a multiline echo
echo -e "C\n(\.   \      ,/)\n  \(   |\     )/\n  //\  | \   /\\\n (/ /\_#oo#_/\ \)\n  \/\  ####  /\/\n       '##'";
echo "This script will install Adam's Zillenial Arch Linux Flavor";
echo "before you can use it, you need to have a working internet connection";
echo "and a working installation drive that base Arch Linux is installed on.";
echo "This script will also install some packages that are not included in the base Arch Linux";

PACKAGES="go rust node-js python3 qtile git firefox thundebird obs-studio xorg steam gimp blender amd-ucode"
AURPACKAGES="rocm-hip-runtime hip-runtime-amd rocm-opencl-runtime vs-code spotify ulauncher xlockmore cava wired"
OTHERS="nvim vim-plug and neovim plugins that I enjoy"

echo "This script will install the following packages: $PACKAGES $AURPACKAGES $OTHERS";
echo "This script will also install azalf's configs in ~/.config/";
echo "This script will also add a .look.yml file to ~/.config.";
echo "It controlls the look and feel and colorscheme";

# Check that the user is root
[ "$EUID" -ne 0 ] && sudo pacman -S --noconfirm "$PACKAGES" || pacman -S --noconfirm "$PACKAGES"
 
# check if there were errors
if [ $? -ne 0 ]; then
    echo "An error has when installing azalf's packages."
    echo "Please ensure that you have an internet connection."
    echo " Run $ sudo pacman -S --noconfirm $PACKAGES"
    exit 1
fi

# Install YAY to continue with the installation
# if the following script does not work,
# then run check the output 


# Ask user if they want to install YAY
echo "Do you want to install YAY? (y/n)"
read -r answer
case $answer in
    [Yy]* )
        pacman -S --noconfirm base-devel &> /dev/null
        # Install YAY
        git clone https://aur.archlinux.org/yay.git
        cd yay
        makepkg -si --noconfirm
        cd ..
        rm -rf yay
        ;;
    [Nn]* )
        # Do nothing
        ;;
    * )
        echo "Please answer yes or no."
        ;;
esac && {
    echo "The following packages will be installed: $AURPACKAGES";
    echo "do you want to install them? (y/n)";
    read -r answer;
    case $answer in
        [Yy]* )
            yay -S --noconfirm "$AURPACKAGES";
            ;;
        [Nn]* )
            echo "You chose not to install the azalf's AUR packages.";
            exit 1; # exit the script
            ;;
        * )
            echo "Please answer yes or no.";
            ;;
    esac;
}

# Install eww
echo "Next we will install eww; it's written in rust."
echo "Do you want to install eww? (y/n) (You must install this to use azalf's configs)"
read -r answer
case $answer in
    [Yy]* )
        git clone https://github.com/elkowar/eww
        cd eww
        cargo build --release
        ;;
    [Nn]* )
        echo "You chose not to install eww.";
        exit 1; # exit the script
        ;;
    * )
        echo "Please answer yes or no.";
        ;;
esac || {
    echo "An error has occurred when installing eww.";
    echo "Please ensure that you have an internet connection.";
    echo "You are free to run the script again.";
    exit 1;
}

echo "One last thing before we continue...";
echo "Neovim is going to be installed alongside vim-plug.";
echo "Do you want to install neovim? (y/n)";
read -r answer;
case $answer in
    [Yy]* )
        pip install neovim
        https://github.com/vim-scripts/vim-plug.git
        ;;
    [Nn]* )
        echo "You chose not to install neovim.";
        exit 1; # exit the script
        ;;
    * )
        echo "Please answer yes or no.";
        ;;
esac || {
    echo "An error has occurred when installing neovim.";
    echo "Please ensure that you have an internet connection.";
    echo "You are free to run the script again.";
    exit 1;
}


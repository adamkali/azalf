#!/bin/bash

load[0]="<\________________>"
load[1]="</\_______________>"
load[2]="<_/\______________>"
load[3]="<__/\_____________>"
load[4]="<___/\____________>"
load[5]="<____/\___________>"
load[6]="<_____/\__________>"
load[7]="<______/\_________>"
load[8]="<_______/\________>"
load[9]="<________/\_______>"
load[10]="<_________/\______>"
load[11]="<__________/\_____>"
load[12]="<___________/\____>"
load[13]="<____________/\___>"
load[14]="<_____________/\__>"
load[15]="<______________/\_>"
load[16]="<_______________/\>"
load[17]="<________________/>"
load[18]="<_______________/\>"
load[19]="<______________/\_>"
load[20]="<_____________/\__>"
load[21]="<____________/\___>"
load[22]="<___________/\____>"
load[23]="<__________/\_____>"
load[24]="<_________/\______>"
load[25]="<________/\_______>"
load[26]="<_______/\________>"
load[27]="<______/\_________>"
load[28]="<_____/\__________>"
load[29]="<____/\___________>"
load[30]="<___/\____________>"
load[31]="<__/\_____________>"
load[32]="<_/\______________>"
load[33]="</\_______________>"

# make a function that takes a pid as a parameter then checks if it is running
# while it is running, it will print the load and loop through the array
function check_pid () {
    # make a variable to hold the pid
    local pid=$1
    # make a whileloop that will run until the pid is killed
    while [ -d /proc/"$pid" ]; do
        for i in {0..36}; do
            echo -ne "${load[$i]}\r"
            sleep 0.1
        done
    done
}


echo "Welcom to the installation script";
# make a multiline echo
echo "This script will install Adam's Zillenial Arch Linux Flavor";
echo "before you can use it, you need to have a working internet connection";
echo "and a working installation drive that base Arch Linux is installed on.";
echo "This script will also install some packages that are not included in the base Arch Linux";

PACKAGES="go rust node-js python3 git firefox thundebird obs-studio xorg steam gimp blender amd-ucode"
AURPACKAGES="rocm-hip-runtime hip-runtime-amd rocm-opencl-runtime vs-code spotify ulauncher xlockmore cava wired"
OTHERS="nvim vim-plug and neovim plugins that I enjoy"

echo "This script will install the following packages: $PACKAGES $AURPACKAGES $OTHERS";
echo "This script will also install azalf's configs in ~/.config/";
echo "This script will also add a .look.yml file to ~/.config.";
echo "It controlls the look and feel and colorscheme";

# Check that the user is root
{
    [ "$EUID" -ne 0 ] &&
    sudo pacman -S --noconfirm go node-js python3 git firefox thunderbird obs-studio steam gimp blender amd-ucode ||
    pacman -S --noconfirm go node-js python3 git firefox thunderbird obs-studio steam gimp blender amd-ucode
} 2>&1

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
    # hold its pid and put it into check_pid function
    pid=$!;
    check_pid "$pid";
} && {
    echo "The following packages will be installed: $AURPACKAGES";
    echo "do you want to install them? (y/n)";
    read -r answer;
    case $answer in
        [Yy]* )
            yay -S --noconfirm rustup rocm-hip-runtime hip-runtime-amd rocm-opencl-runtime vs-code spotify ulauncher xlockmore cava wired;
            ;;
        [Nn]* )
            echo "You chose not to install the azalf's AUR packages.";
            exit 1; # exit the script
            ;;
        * )
            echo "Please answer yes or no.";
            ;;
    esac;
} || {
    # if there was an error, print it
    echo "An error has occurred."
    echo "Please ensure that you have an internet connection."
    echo " Run $ yay -S --noconfirm $AURPACKAGES"
}

# Install eww

echo "Next we will install eww; its written in rust."
echo "Do you want to install eww? \(y/n\) \(You must install this to use azalfs configs\)"
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
echo "Do you want to install neovim and its components? (y/n)";
read -r answer;
case $answer in
    [Yy]* )
        pip install neovim
        git clone https://github.com/vim-scripts/vim-plug.git
        ;;
    [Nn]* )
        echo "You chose not to install neovim.";
        exit 1; # exit the script
        ;;
    * )
        echo "Please answer yes or no.";
        ;;
esac & pid=$! || {
    echo "An error has occurred when installing neovim and its components.";
    echo "Please ensure that you have an internet connection.";
    echo "You are free to run the script again.";
    exit 1;
};

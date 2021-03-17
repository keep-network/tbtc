#!/bin/bash

set -euo pipefail

LOG_START='\n\e[1;36m'  # new line + bold + cyan
LOG_END='\n\e[0m'       # new line + reset

RELAY_PATH=$(realpath $(dirname $0)/../)

# Defaults, can be overwritten by env variables/input parameters
LOG_LEVEL_DEFAULT="info"
CONFIG_DIR_PATH_DEFAULT="$RELAY_PATH/config"
CONFIG_DIR_PATH=$(realpath "${CONFIG_DIR_PATH:-$CONFIG_DIR_PATH_DEFAULT}")

# Read user inputs.
OPERATOR_KEY_FILE_PASSWORD_DEFAULT=password
read -p "Enter operator key file \
password [$OPERATOR_KEY_FILE_PASSWORD_DEFAULT]: " operator_password
OPERATOR_KEY_FILE_PASSWORD=${operator_password:-$OPERATOR_KEY_FILE_PASSWORD_DEFAULT}

help()
{
   echo -e "\nUsage: ENV_VAR(S) $0"
   echo -e "\nEnvironment variables:"
   echo -e "\tCONFIG_DIR_PATH: Location of relay config file(s)." \
           "Default value is 'config' dir placed under project root."
   exit 1 # Exit script after printing help
}

while getopts "h" opt
do
   case "$opt" in
      h ) help ;;
      ? ) help ;; # Print help in case parameter is non-existent
   esac
done

config_files=($CONFIG_DIR_PATH/*.toml)
config_files_count=${#config_files[@]}
while :
do
    printf "\nSelect config file: \n"
    i=1
    for o in "${config_files[@]}"; do
        echo "$i) ${o##*/}"
        let i++
    done

    read reply
    if [ "$reply" -ge 1 ] && [ "$reply" -le $config_files_count ]; then
        CONFIG_FILE_PATH=${config_files["$reply"-1]}
        break
    else
        printf "\nInvalid choice. Please choose an existing option number.\n"
    fi
done
printf "\nConfig file: \"$CONFIG_FILE_PATH\" \n\n"

log_level_options=("info" "debug" "custom...")
while :
do
    echo "Select log level [$LOG_LEVEL_DEFAULT]: "
    i=1
    for o in  "${log_level_options[@]}"; do
        echo "$i) $o"
        let i++
    done

    read reply
    case $reply in
        "1"|"${log_level_options[0]}") LOG_LEVEL=${log_level_options[0]}; break;;
        "2"|"${log_level_options[1]}") LOG_LEVEL=${log_level_options[1]}; break;;
        "3"|"${log_level_options[2]}")
            read -p "Enter custom log level: [$LOG_LEVEL_DEFAULT]" log_level
            LOG_LEVEL=${log_level:-$LOG_LEVEL_DEFAULT}
            break
            ;;
        "") LOG_LEVEL=$LOG_LEVEL_DEFAULT; break;;
        *) echo "Invalid choice. Please choose an existing option number.";;
    esac
done
echo "Log level: \"$LOG_LEVEL\""

printf "${LOG_START}Starting relay node...${LOG_END}"
OPERATOR_KEY_FILE_PASSWORD=$OPERATOR_KEY_FILE_PASSWORD \
  LOG_LEVEL=${LOG_LEVEL} \
  ./relay --config "$(realpath $CONFIG_FILE_PATH)" start
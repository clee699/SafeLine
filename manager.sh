#!/bin/bash

while true; do
    echo "=============================="
    echo "         SafeLine Manager     "
    echo "=============================="
    echo "1. Install"
    echo "2. Upgrade"
    echo "3. Upgrade to LTS"
    echo "4. Uninstall"
    echo "5. Repair/Restart"
    echo "6. View Status"
    echo "7. View Logs"
    echo "8. Backup Configuration"
    echo "9. System Information"
    echo "10. Exit"
    read -p "Select an option [1-10]: " option

    case ${option} in
        1) echo "Installing...";;
        2) echo "Upgrading...";;
        3) echo "Upgrading to LTS...";;
        4) echo "Uninstalling...";;
        5) echo "Repairing/Restarting...";;
        6) echo "Viewing Status...";;
        7) echo "Viewing Logs...";;
        8) echo "Backing up configuration...";;
        9) echo "Displaying system information...";;
        10) echo "Exiting..."; exit 0;;
        *) echo "Invalid option, please try again.";;
    esac
done

#!/bin/bash

print_tree() {
    local indent="$2"
    for file in "$1"/*; do
        if [ -d "$file" ]; then
            dir_name=$(basename "$file")
            # Skip if the directory is named 'node_modules'
            if [ "$dir_name" == "node_modules" ]; then
                continue
            fi
            echo "${indent}${dir_name}/"
            print_tree "$file" "$indent    "
        else
            echo "${indent}$(basename "$file")"
        fi
    done
}

current_dir=$(pwd)
echo "$(basename "$current_dir")/"
print_tree "$current_dir" "    "


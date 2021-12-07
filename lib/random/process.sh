#!/bin/sh
echo 'package random'
echo
echo 'var wordList = []string{'
cat cleanwords.txt |
    tr 'A-Z' 'a-z' |                # Convert all uppercase to lowercase
    sed 's/[^ A-Za-z]/ /g' |        # Convert anything not a space or number to a space
    sed 's/[[:space:]]\+/\n/g' |    # Convert any repeated spaces to newlines
    sed '/^$/d' |                   # Delete blank lines
    sed 's/^./\u&/g' |              # Capitalize first letter of each line
    sort |                          # Need to sort before we find unique words
    uniq |                          # Keep only the unique words in the list
    sed 's/^\(\w*\)$/\t"\1",/g' |     # Add surrounding quotes and a comma
    cat
echo '}'
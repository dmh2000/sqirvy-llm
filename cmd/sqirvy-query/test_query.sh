#!/bin/sh
check_return_code() {
    local cmd="$1"
    $cmd $2 $3 $4 $5 $6 $7 $8 $9
    local return_code=$?
    
    if [ $return_code -eq 0 ]; then
        echo "Command '$cmd' executed successfully"
    else
        echo "Command '$cmd' failed with exit code $return_code"
        exit 1
    fi
    
    return $return_code
}

ignore_return_code() {
    local cmd="$1"
    $cmd $2 $3 $4 $5 $6 $7 $8 $9
    local return_code=$?
    
    return 0
}

make build 
echo "-------------------------------"
echo "sqiry-query (should fail)"
ignore_return_code ../../bin/sqirvy-query
echo "-------------------------------"
echo "sqiry-query -h"
check_return_code ../../bin/sqirvy-query -h
echo "-------------------------------"
echo "../../bin/sqirvy-query < hello.txt"
check_return_code ../../bin/sqirvy-query < hello.txt
echo "sqiry-query file"
echo "-------------------------------"
echo "../../bin/sqirvy-query hello.txt"
check_return_code ../../bin/sqirvy-query hello.txt
echo "-------------------------------"
echo "../../bin/sqirvy-query goodbye.txt < hello.txt"
check_return_code ../../bin/sqirvy-query goodbye.txt < hello.txt
echo "-------------------------------"
echo "../../bin/sqirvy-query  hello.txt goodbye.txt"
check_return_code ../../bin/sqirvy-query  hello.txt goodbye.txt
echo "-------------------------------"
echo "../../bin/sqirvy-query  -m gemini-2.0-flash-exp < hello.txt goodbye.txt"
check_return_code ../../bin/sqirvy-query  -m gemini-2.0-flash-exp < hello.txt goodbye.txt
echo "-------------------------------"
echo "../../bin/sqirvy-query  -m gpt-4o < hello.txt goodbye.txt"
check_return_code ../../bin/sqirvy-query  -m gpt-4o < hello.txt goodbye.txt
echo "-------------------------------"
echo "../../bin/sqirvy-query  -m gpt-4o < hello.txt goodbye.txt"
ignore_return_code ../../bin/sqirvy-query  -m xyz < hello.txt 
echo "-------------------------------"
echo "../../bin/sqirvy-query  -m gpt-4o < hello.txt goodbye.txt"
ignore_return_code ../../bin/sqirvy-query  xyz
echo "-------------------------------"

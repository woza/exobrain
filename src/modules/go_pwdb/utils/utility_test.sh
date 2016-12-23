#!/bin/bash

# This script runs through a few simple scenarios using the provided utilities

function clean_up()
{
    echo -n "Cleaning up old database..."
    rm -f utility_test.db
    rm -f utility_test.conf
    echo "done"
}

function make_db()
{
    echo -n "Creating new database..."
    cat > fake_initdb.in <<EOF
utility_test.conf
utility_test.db
fake_ui.key
fake_ui.crt
fake_ui.ca
fake_display.key
fake_display.crt
fake_display.ca
127.0.0.1:7777
0.0.0.0:6666
fake_display
utilitytestpass
utilitytestpass
EOF
    ./initdb < fake_initdb.in > /dev/null
    rm -f fake_initdb.in
    echo "done"    	 
}

function dump_tags()
{
    echo -e "utilitytestpass\n" | ./list utility_test.conf | sort > utility_test.tags
}

function add_db_entry() # username, password
{
    user=$1
    pw=$2
    echo -e "$pw\n$pw\nutilitytestpass\n" | ./put utility_test.conf $user > /dev/null
}

function fill_db()
{
    echo -n "Populating database..."
    add_db_entry major bloodnok 
    add_db_entry harry seagoon
    add_db_entry count moriarty
    dump_tags
    cat > utility_test.expect <<EOF
count
Enter database password: 
harry
Known tags:
major
Parsing config file  utility_test.conf
EOF
    diff --brief utility_test.expect utility_test.tags || exit 1
    echo "done"
}

function add_db_entry() # username, password
{
    user=$1
    pw=$2
    echo -e "$pw\n$pw\nutilitytestpass\n" | ./put utility_test.conf $user > /dev/null
}

function remove_db_entry() #tag
{
    user=$1
    echo -e "utilitytestpass\n" | ./delete utility_test.conf $user > /dev/null
}
    
function trim_db()
{
    echo -n "Shrinking database..."
    remove_db_entry count moriarty
    dump_tags
    cat > utility_test.expect <<EOF
Enter database password: 
harry
Known tags:
major
Parsing config file  utility_test.conf
EOF
    diff --brief utility_test.expect utility_test.tags || exit 1
    echo "done"
}

function rekey_db()
{
    echo -n "Rekeying database..."
    # Assumes utility_test.expect already exists - rekeying the database
    # should not change reported keys.
    echo -e "changed\nchanged\nutilitytestpass\n" | ./rekey utility_test.conf > /dev/null
    echo -e "changed\n" | ./list utility_test.conf | sort > utility_test.tags
    diff --brief utility_test.expect utility_test.tags || exit 1
    echo "done"
}
    
    
clean_up
make_db
fill_db
trim_db
rekey_db
echo "All tests passed"

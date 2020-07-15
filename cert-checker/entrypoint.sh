#!/bin/sh -l

args="${INPUT_CMD:-$*}"
result=$(/go/bin/cert $args)
echo "::set-output name=result::$result"

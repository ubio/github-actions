#!/bin/sh -l

result=$(/go/bin/cert $@)
echo "::set-output name=result::$result"

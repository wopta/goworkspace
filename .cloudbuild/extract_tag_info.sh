#!/bin/bash
TAG_NAME=$1

# === EXTRACT MODULE NAME AND VERSION FROM TAG =================================
# Read TAG_NAME splitting by the IFS separator into the needed variables
IFS='/'
read -a TAG_ARRAY <<< $TAG_NAME

echo ${TAG_ARRAY[0]} > /workspace/tag_module
echo ${TAG_ARRAY[1]} > /workspace/tag_version

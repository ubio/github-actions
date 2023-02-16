#!/bin/bash

function parse_inputs {

    # optional inputs
    kustomize_build_dir="."
    if [ "${INPUT_KUSTOMIZE_BUILD_DIR}" != "" ] || [ "${INPUT_KUSTOMIZE_BUILD_DIR}" != "." ]; then
        kustomize_build_dir=${INPUT_KUSTOMIZE_BUILD_DIR}
    fi

    kustomize_comment=0
    if [ "${INPUT_KUSTOMIZE_COMMENT}" == "1" ] || [ "${INPUT_KUSTOMIZE_COMMENT}" == "true" ]; then
        kustomize_comment=1
    fi

    kustomize_install=1
    if [ "${INPUT_KUSTOMIZE_INSTALL}" == "0" ] || [ "${INPUT_KUSTOMIZE_INSTALL}" == "false" ]; then
        kustomize_install=0
    fi

    kustomize_output_file=""
    if [ -n "${INPUT_KUSTOMIZE_OUTPUT_FILE}" ]; then
      kustomize_output_file=${INPUT_KUSTOMIZE_OUTPUT_FILE}
    fi

    kustomize_build_options=""
    if [ -n "${INPUT_KUSTOMIZE_BUILD_OPTIONS}" ]; then
      kustomize_build_options=${INPUT_KUSTOMIZE_BUILD_OPTIONS}
    fi

    enable_alpha_plugins=""
    if [ "${INPUT_ENABLE_ALPHA_PLUGINS}" == "1" ] || [ "${INPUT_ENABLE_ALPHA_PLUGINS}" == "true" ]; then
       enable_alpha_plugins="--enable_alpha_plugins"
    fi

    with_token=""
    if [ "${INPUT_TOKEN}" != "" ]; then
       with_token=(-H "Authorization: token ${INPUT_TOKEN}")
    fi
}

function main {

    scriptDir=$(dirname ${0})
    source ${scriptDir}/kustomize_build.sh
    parse_inputs

    kustomize_build

}

main "${*}"
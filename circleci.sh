#!/bin/bash

# Validate
circleci config validate -c .circleci/config.yml || exit -1

# Run
circleci build

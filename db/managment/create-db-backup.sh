#!/bin/bash

sudo -u postgres pg_dump -Fc codestep_db > dump-codestep_db-$(date +'%Y-%m-%d-%H-%M-%S')

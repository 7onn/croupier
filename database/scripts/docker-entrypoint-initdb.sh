#!/bin/bash

cd ${WORKDIR}

createdb --username postgres $DATABASE_NAME

psql --username postgres $DATABASE_NAME -f ./scripts/setup.sql -v ON_ERROR_STOP=ON -q
psql --username postgres $DATABASE_NAME -f ./scripts/initial-schema.sql -v ON_ERROR_STOP=ON -q

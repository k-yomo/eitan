#!/bin/bash

set -euC

PROJECT=local
HOST=localhost:8085

# account.user-registration
pubsub_cli create_subscription account.user-registration notification.account.user-registration --create-if-not-exist -p $PROJECT -h $HOST
pubsub_cli create_subscription account.user-registration eitan.account.user-registration --create-if-not-exist -p $PROJECT -h $HOST

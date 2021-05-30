#!/bin/bash

set -euC

PROJECT=local
HOST=localhost:8085

# account.user-registered
pubsub_cli create_subscription account.user-registered notification.account.user-registered --create-if-not-exist -p $PROJECT -h $HOST
pubsub_cli create_subscription account.user-registered eitan.account.user-registered --create-if-not-exist -p $PROJECT -h $HOST

# account.email-confirmation-created
pubsub_cli create_subscription account.email-confirmation-created notification.account.email-confirmation-created --create-if-not-exist -p $PROJECT -h $HOST

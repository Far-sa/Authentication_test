#!/bin/bash

# Identify container ID or name for service 1
container_id=$(docker ps -q --filter name=user-svc)

# Send SIGTERM signal to service 1 container
docker kill --signal SIGTERM $container_id

# (Optional) Wait for the server to finish shutting down
# sleep 10 

# (Optional) Force shutdown if necessary (use with caution)
# docker kill --signal KILL $container_id

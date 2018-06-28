#!/bin/bash

# saving all docker images of eta along with its dependencies
docker save -o influxdb-1.5.3.tar influxdb:1.5.3
docker save -o redis-4.0.10.tar redis:4.0.10
docker save -o nats-1.1.0.tar nats:1.1.0
docker save -o mosquitto-1.4.12.tar eclipse-mosquitto:1.4.12
docker save -o postgres-10.4.tar postgres:10.4

docker save -o data_agent-1.0.tar ia/data_agent:1.0
docker save -o classifier-1.0.tar ia/data_analytics:1.0
docker save -o nats_client-1.0.tar ia/nats_client:1.0
docker save -o video_ingestion-1.0-latest.tar ia/video_ingestion:1.0
docker exec -it elasticsearch \
  bin/elasticsearch-reset-password \
    -u elastic \
    --batch

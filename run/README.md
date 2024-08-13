# 0GDA

## How to run

1. Build client server

    ```bash

    make build
    ```

2. Build docker image

    ```bash
    docker build -t 0gclient -f Dockerfile .
    ```

3. Update configutaions in **run.sh** and **start.sh**

4. Run docker image
    ```bash
    ./start.sh
    ```


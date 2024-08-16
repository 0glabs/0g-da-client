# 0GDA

## How to run

1. Build client server

    ```bash
    cd run
    make -C .. build
    ```

2. Build docker image

    ```bash
    cd ..
    docker build -t 0gclient -f ./Dockerfile .
    ```

3. Update configutaions in **[run.sh](run.sh)**

4. Grant executable permissions to **[run.sh](run.sh)** and **[start.sh](start.sh)** when needed.

5. Run docker image
    ```bash
    cd run
    ./start.sh
    ```


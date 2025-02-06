# streaming-analytics-example (wip)

A local example showing how to use Kafka and Stream Analytics to process real-time customer feedback.


## Prerequisites
- Docker / Podman
- just (https://github.com/casey/just)

## Running the example
1. Create a docker network
    ```
    docker network create streaming-analytics-example
    ```
2. Run the Kafka container
    ```
    cd kafka
    just up
    just logs
    ```

3. Run the WebApp container
    ```
    cd webapp
    just up
    ```

4. Run the Flink cluster
    ```
    cd flink
    just up
    ```

5. Submit the streaming job to flink
    ```
    sdk install java 17.0.14-amzn
    sdk use java 17.0.14-amzn

    cd flink
    just submit
    ```


4. Open 2 browser tabs
    - Kafka UI: http://localhost:8081
    - WebApp: http://localhost:8080



## TODO
- Spin up a Flink cluster
- Create a Flink job to create a materialized view of the topic
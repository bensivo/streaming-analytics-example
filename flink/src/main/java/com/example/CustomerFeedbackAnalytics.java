package com.example;

import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.api.common.serialization.SimpleStringSchema;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.table.api.Table;
import org.apache.flink.table.api.bridge.java.StreamTableEnvironment;
import org.apache.flink.types.Row;

public class CustomerFeedbackAnalytics {
    public static void main(String[] args) throws Exception {
        // Set up the streaming execution environment
        StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();
        StreamTableEnvironment tableEnv = StreamTableEnvironment.create(env);

        // // Configure Kafka source
        // KafkaSource<String> source = KafkaSource.<String>builder()
        //         .setBootstrapServers("kafka:9092")
        //         .setTopics("event.customer_feedback")
        //         .setGroupId("flink-customer-feedback-analytics")
        //         .setStartingOffsets(OffsetsInitializer.latest())
        //         .setValueOnlyDeserializer(new SimpleStringSchema())
        //         .setProperty("auto.offset.reset", "earliest")
        //         .setProperty("enable.auto.commit", "true")
        //         .build();

        // // Create a DataStream from Kafka messages
        // DataStream<String> kafkaStream = env.fromSource(
        //         source,
        //         WatermarkStrategy.noWatermarks(),
        //         "Kafka Source"
        // );

        // Create a temporary table from the Kafka stream
        tableEnv.executeSql(
            "CREATE TEMPORARY TABLE customer_feedback (" +
            "  version INT," +
            "  ts TIMESTAMP(3)," +
            "  rating INT," +
            "  WATERMARK FOR ts AS ts - INTERVAL '5' SECONDS" +
            ") WITH (" +
            "  'connector' = 'kafka'," +
            "  'topic' = 'event.customer_feedback'," +
            "  'properties.bootstrap.servers' = 'kafka:9092'," +
            "  'properties.group.id' = 'flink-customer-feedback-analytics'," +
            "  'scan.startup.mode' = 'earliest-offset'," +
            "  'format' = 'json'" +
            ")"
        );

        // Create a view with the last minute of feedback
        Table recentFeedback = tableEnv.sqlQuery(
            "SELECT " +
            "  TUMBLE_START(ts, INTERVAL '1' MINUTE) AS window_start," +
            "  COUNT(*) AS total_feedback," +
            "  AVG(CAST(rating AS DOUBLE)) AS avg_rating," +
            "  SUM(CASE WHEN rating >= 3 THEN 1 ELSE 0 END) AS positive_feedback," +
            "  SUM(CASE WHEN rating < 3 THEN 1 ELSE 0 END) AS negative_feedback " +
            "FROM customer_feedback " +
            "GROUP BY TUMBLE(ts, INTERVAL '1' MINUTE)"
        );

        // Convert the table back to a stream and print results
        tableEnv.toDataStream(recentFeedback, Row.class).print();

        // Execute the Flink job
        env.execute("Customer Feedback Analytics");
    }
}

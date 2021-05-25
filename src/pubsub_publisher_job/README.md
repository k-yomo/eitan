# PubSub Publisher Job

This job is used to periodically publish Pub/Sub messages stored in DB.
The purpose for storing messages in DB is written [here](https://github.com/GoogleCloudPlatform/transactional-microservice-examples#notes-on-the-event-publishing-process).

## Usage
1. First prepare table with [fields](./pubsub_event.go).
```sql
CREATE TABLE {YOUR_DB}.pubsub_events (
  id VARCHAR(255) NOT NULL PRIMARY KEY COMMENT 'イベントID',
  deduplicate_key VARCHAR(255) COMMENT '重複排除キー',
  topic VARCHAR(255) NOT NULL COMMENT 'トピック名',
  data TEXT NOT NULL COMMENT 'データ',
  is_published BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Publish済みか否か',
  published_at DATETIME COMMENT 'Publishした日時',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX is_published_idx (is_published)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'PubSubイベント';
```

2. Add pubsub-publisher-job to [Procfile](/Procfile) for local development

3. Configure Kubernates Job for dev/prod env

You can refer to the existing [cronjob](/k8s/account-service/pubsub-publisher-job.yaml).

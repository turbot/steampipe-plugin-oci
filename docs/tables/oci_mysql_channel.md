---
title: "Steampipe Table: oci_mysql_channel - Query OCI MySQL Channels using SQL"
description: "Allows users to query channels in Oracle Cloud Infrastructure's MySQL Database Service."
---

# Table: oci_mysql_channel - Query OCI MySQL Channels using SQL

Oracle Cloud Infrastructure's MySQL Database Service is a fully managed database service that lets developers and database administrators build, deploy, run, and scale MySQL applications in the cloud. Channels in MySQL Database Service represent asynchronous replicas that use MySQL's native asynchronous replication to replicate data from a source MySQL DB System to a target MySQL DB System. They are used to increase data durability, support failover, and migrate data across different MySQL DB Systems.

## Table Usage Guide

The `oci_mysql_channel` table provides insights into channels within Oracle Cloud Infrastructure's MySQL Database Service. As a database administrator, you can explore channel-specific details through this table, including the source and target endpoints, the replication status, and the last error message. Utilize it to uncover information about channels, such as those with replication errors, the status of data transfer between source and target MySQL DB Systems, and the configuration details of the channel.

## Examples

### Basic info
Explore the lifecycle state and creation time of your MySQL channels to understand their current status and longevity. This can be beneficial for assessing the health and stability of your database channels.

```sql
select
  display_name,
  id,
  description,
  lifecycle_state as state,
  time_created
from
  oci_mysql_channel;
```

### List disabled channels
Explore which MySQL channels are currently inactive. This query is useful for maintaining security and performance by identifying channels that are not enabled, allowing for necessary management actions.

```sql
select
  display_name,
  id,
  description,
  lifecycle_state as state,
  time_created,
  time_updated,
  is_enabled
from
  oci_mysql_channel
where
  not is_enabled;
```

### Get target details for each channel
Identify the specific details for each channel, such as the applier's username, channel name, system ID, and target type, to gain insights into your MySQL Channel configuration and understand its operational aspects.

```sql
select
  display_name,
  id,
  target ->> 'applierUsername' as applier_username,
  target ->> 'channelName' as channel_name,
  target ->> 'dbSystemId' as db_system_id,
  target ->> 'targetType' as target_type
from
  oci_mysql_channel;
```
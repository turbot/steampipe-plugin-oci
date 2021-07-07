# Table: oci_mysql_channel

A Channel connecting a DB System to an external entity.

## Examples

### Basic info

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


### List target details of channels

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
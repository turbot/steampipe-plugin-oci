---
title: "Steampipe Table: oci_database_software_image - Query OCI Database Software Images using SQL"
description: "Allows users to query OCI Database Software Images"
---

# Table: oci_database_software_image - Query OCI Database Software Images using SQL

Oracle Cloud Infrastructure (OCI) Database Software Images are custom images that contain the database software for creating a database on a Database system. They are used to launch Database systems with a specified version of the database software, a database configuration, and a preinstalled set of database options and patches. This resource provides a way to manage and customize the database software that runs on your Database systems.

## Table Usage Guide

The `oci_database_software_image` table provides insights into Database Software Images within Oracle Cloud Infrastructure (OCI). As a database administrator, you can explore software image-specific details through this table, including versions, configurations, and preinstalled database options and patches. Utilize it to manage and customize the database software that runs on your Database systems, ensuring the systems are launched with the correct configurations and database versions.

## Examples

### Basic info
Explore which database software images are currently active and when they were created, to gain insights into the life cycle and type of database images in your Oracle Cloud Infrastructure. This can help in managing your resources effectively and keeping track of your database images.

```sql+postgres
select
  display_name,
  id,
  image_type,
  lifecycle_state,
  time_created
from
  oci_database_software_image;
```

```sql+sqlite
select
  display_name,
  id,
  image_type,
  lifecycle_state,
  time_created
from
  oci_database_software_image;
```

### List deleted database software images
Explore which database software images have been deleted to manage your resources effectively. This helps in maintaining a clean and efficient database by keeping track of unused or obsolete software images.

```sql+postgres
select
  display_name,
  id,
  image_type,
  lifecycle_state,
  time_created
from
  oci_database_software_image
where
  lifecycle_state = 'DELETED';
```

```sql+sqlite
select
  display_name,
  id,
  image_type,
  lifecycle_state,
  time_created
from
  oci_database_software_image
where
  lifecycle_state = 'DELETED';
```

### List database software images older than 90 days
Explore the database software images that have been created more than 90 days ago. This can help in identifying outdated images, potentially facilitating their update or removal.

```sql+postgres
select
  display_name,
  id,
  image_type,
  lifecycle_state,
  time_created
from
  oci_database_software_image
where
  time_created <= (current_date - interval '90' day)
order by
  time_created;
```

```sql+sqlite
select
  display_name,
  id,
  image_type,
  lifecycle_state,
  time_created
from
  oci_database_software_image
where
  time_created <= date('now','-90 day')
order by
  time_created;
```
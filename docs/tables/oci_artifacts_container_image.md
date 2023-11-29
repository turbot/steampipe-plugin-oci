---
title: "Steampipe Table: oci_artifacts_container_image - Query OCI Artifacts Container Images using SQL"
description: "Allows users to query OCI Artifacts Container Images."
---

# Table: oci_artifacts_container_image - Query OCI Artifacts Container Images using SQL

Oracle Cloud Infrastructure (OCI) Artifacts is a fully managed, scalable, and secure artifact storage and sharing service. OCI Artifacts allows you to share container images, development artifacts, and other software dependencies within your team or organization. Container Images are a lightweight, standalone, executable package that includes everything needed to run a piece of software, including the code, a runtime, libraries, environment variables, and config files.

## Table Usage Guide

The `oci_artifacts_container_image` table provides insights into Container Images within Oracle Cloud Infrastructure (OCI) Artifacts. As a DevOps engineer, explore image-specific details through this table, including image digest, created time, and size. Utilize it to uncover information about images, such as those with specific digests, the creation time of images, and the size of container images.

## Examples

### Basic info
Explore which OCI artifacts container images are most frequently pulled, allowing you to assess their popularity and usage. This can help in managing resources effectively by identifying the images that require more storage and bandwidth.

```sql
select
  display_name,
  id,
  created_by,
  digest,
  layers,
  layers_size_in_bytes,
  manifest_size_in_bytes,
  pull_count,
  repository_id,
  repository_name,
  versions,
  time_last_pulled,
  version,
  lifecycle_state as state
from
  oci_artifacts_container_image;
```

### Get the size of the largest image layer (in bytes)
Explore which image layer in your OCI artifacts container holds the most data. This can be useful for optimizing storage use or identifying unusually large layers that may need further investigation.

```sql
select
  display_name,
  id,
  digest,
  time_created,
  layers_size_in_bytes
from
  oci_artifacts_container_image
order by
  layers_size_in_bytes desc limit 1;
```

### Get version details of each image
Explore the creation details of various versions of specific images. This can help in understanding the evolution of an image over time, which is crucial for maintaining version control and tracking changes.

```sql
select
  i.display_name,
  i.id as image_id,
  v ->> 'createdBy' as image_version_created_by,
  v ->> 'timeCreated' as image_version_created_time,
  v ->> 'version' as version
from
  oci_artifacts_container_image as i,
  jsonb_array_elements(versions) as v;
```

### Get layer details of each image
Explore the different aspects of each image, such as layer details, by analyzing its unique identifiers, size, and creation time. This can be beneficial in managing storage and understanding image creation patterns.

```sql
select
  display_name,
  id,
  l ->> 'digest' as layer_digest,
  l ->> 'sizeInBytes' as layer_size_in_bytes,
  l ->> 'timeCreated' as layer_create_time
from
  oci_artifacts_container_image,
  jsonb_array_elements(layers) as l;
```

### Get repository details of each image
Explore the characteristics of each image by analyzing details such as its repository, immutability status, and public visibility. This can help in understanding the lifecycle state of each image and its repository, aiding in better management and organization of resources.

```sql
select
  i.display_name,
  i.id,
  i.repository_id,
  r.display_name as repository_display_name,
  r.is_immutable,
  r.is_public,
  r.lifecycle_state
from
  oci_artifacts_container_image as i,
  oci_artifacts_container_repository as r
where
  r.id = i.repository_id;
```

### List available images
Explore the currently available images in your OCI Artifacts repository. This can be useful in maintaining an up-to-date inventory or identifying images for potential updates or removal.

```sql
select
  display_name,
  id,
  digest,
  version,
  lifecycle_state
from
  oci_artifacts_container_image
where
  lifecycle_state = 'AVAILABLE';
```

### List images created in last 30 days
Discover the latest images that have been created within the past month. This could be useful for keeping track of recent additions or changes to your system.

```sql
select
  display_name,
  id,
  digest,
  time_created,
  manifest_size_in_bytes
from
  oci_artifacts_container_image
where
  time_created >= now() - interval '30' day;
```

### Retrive the total number of pull count of each image
Analyze the popularity of various container images by determining the total number of times each has been pulled. This can be useful for understanding which images are most frequently used.

```sql
select
  display_name,
  id,
  digest,
  pull_count
from
  oci_artifacts_container_image;
```
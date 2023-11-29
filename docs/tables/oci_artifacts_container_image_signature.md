---
title: "Steampipe Table: oci_artifacts_container_image_signature - Query OCI Artifacts Container Image Signatures using SQL"
description: "Allows users to query OCI Artifacts Container Image Signatures."
---

# Table: oci_artifacts_container_image_signature - Query OCI Artifacts Container Image Signatures using SQL

Oracle Cloud Infrastructure's (OCI) Artifacts service is a fully managed, scalable, and secure artifact storage and sharing service. OCI Artifacts allows you to share container images within and across regions and tenancies. OCI Artifacts Container Image Signatures are cryptographic proofs that are used to verify the authenticity and integrity of OCI Artifacts Container Images.

## Table Usage Guide

The `oci_artifacts_container_image_signature` table provides insights into the cryptographic signatures associated with OCI Artifacts Container Images. As a security analyst, you can use this table to explore signature-specific details, including the cryptographic algorithm used, the key used for signing, and the signature value. This can be beneficial for verifying the authenticity and integrity of the container images in your OCI environment.

## Examples

### Basic info
Explore which OCI artifacts container image signatures have been created, by whom, and with what key and algorithm. This can help in understanding the security and accountability aspects of your OCI artifacts.

```sql
select
  display_name,
  id,
  created_by,
  image_id,
  kms_key_id,
  kms_key_version_id,
  message,
  signature,
  signing_algorithm
from
  oci_artifacts_container_image_signature;
```

### List signatures created in last 30 days
Explore which signatures have been created in the past 30 days. This is useful for keeping track of recent activity and ensuring all newly created signatures are valid and authorized.

```sql
select
  display_name,
  id,
  time_created,
  image_id,
  message,
  signature
from
  oci_artifacts_container_image_signature
where
  time_created >= now() - interval '30' day;
```

### Get image details of each signature
Explore the details of each digital signature associated with a container image, including the identity of the image and its lifecycle state. This can be useful to understand the usage and status of these images, especially in scenarios where image verification and integrity are crucial.

```sql
select
  s.display_name,
  s.id,
  s.signature,
  s.signing_algorithm,
  s.image_id,
  i.digest,
  i.lifecycle_state,
  i.manifest_size_in_bytes,
  i.pull_count
from
  oci_artifacts_container_image_signature as s,
  oci_artifacts_container_image as i
where
  i.id = s.image_id;
```

### Get KMS key details used by each image signature
Determine the specific encryption key details associated with each image signature to gain insights into security measures. This can help identify any irregularities or potential vulnerabilities in the encryption process.

```sql
select
  s.display_name,
  s.id,
  s.kms_key_version_id,
  v.key_id,
  v.vault_id,
  v.public_key,
  v.origin
from
  oci_artifacts_container_image_signature as s,
  oci_kms_key_version as v
where
  v.id = s.kms_key_version_id;
```

### List signatures with RSA signining algorithm
Determine the areas in which the RSA signing algorithm is used for signatures. This is beneficial for assessing the security measures in place across different segments.

```sql
select
  display_name,
  id,
  message,
  signature,
  signing_algorithm
from
  oci_artifacts_container_image_signature
where
  signing_algorithm = 'RSA';
```
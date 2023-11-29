---
title: "Steampipe Table: oci_certificates_management_certificate - Query OCI Certificates Management Certificates using SQL"
description: "Allows users to query Certificates Management Certificates in Oracle Cloud Infrastructure."
---

# Table: oci_certificates_management_certificate - Query OCI Certificates Management Certificates using SQL

Oracle Cloud Infrastructure (OCI) Certificates Management is a service that helps in managing and controlling SSL/TLS certificates. It allows users to create, import, and manage SSL/TLS certificates for their cloud-based applications and resources. This service ensures that the certificates are valid, trusted, and can be used for secure connections.

## Table Usage Guide

The `oci_certificates_management_certificate` table provides insights into the SSL/TLS certificates managed by OCI Certificates Management. As a security analyst, you can delve into certificate-specific details through this table, including certificate type, validity, and associated metadata. Utilize it to uncover information about certificates, such as their expiration dates, the resources they're associated with, and their overall status within your OCI environment.

## Examples

### Basic info
Discover the details of your certificates, including their lifecycle state and associated metadata. This is useful to understand the status and configuration of your certificates for better management and security compliance.

```sql
select
  id,
  name,
  config_type,
  issuer_certificate_authority_id,
  description,
  certificate_rules,
  time_of_deletion,
  lifecycle_details,
  current_version,
  subject,
  certificate_revocation_list_details,
  key_algorithm,
  signature_algorithm,
  certificate_profile_type,
  lifecycle_state as state
from
  oci_certificates_management_certificate;
```

### List imported certificates
Explore imported certificates to understand their lifecycle state, rules, and signature algorithm. This can be helpful in managing and understanding the security and validity of certificates in your environment.

```sql
select
  name,
  id,
  certificate_rules,
  lifecycle_state,
  config_type,
  signature_algorithm
from
  oci_certificates_management_certificate
where
  config_type = 'IMPORTED';
```

### List failed certificates
Discover the segments that have failed certificates to gain insights into potential security risks or issues in your system. This is useful for quickly identifying and addressing problematic areas, enhancing the overall security of your system.

```sql
select
  name,
  id,
  current_version,
  subject,
  certificate_revocation_list_details,
  key_algorithm
from
  oci_certificates_management_certificate
where
  lifecycle_state = 'FAILED';
```

### List certificates created in the last 30 days
Discover the segments that have been certified within the last month. This is useful for tracking recent changes and ensuring system security.

```sql
select
  name,
  id,
  current_version,
  key_algorithm,
  time_created,
  time_of_deletion
from
  oci_certificates_management_certificate
where
  time_created >= now() - interval '30' day;
```

### Get the current version details of each certificate
Explore which certificates are currently in use and their respective details. This is useful for tracking certificate versions and their revocation status, which aids in maintaining secure connections.

```sql
select
  name,
  id,
  current_version ->> 'certificateId' as certificate_id,
  current_version ->> 'issuerCaVersionNumber' as issuer_ca_version_number,
  current_version ->> 'revocationStatus' as revocation_status,
  current_version ->> 'serialNumber' as serial_number,
  current_version ->> 'stages' as stages
from
  oci_certificates_management_certificate;
```
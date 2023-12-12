---
title: "Steampipe Table: oci_certificates_management_certificate_version - Query OCI Certificates Management Certificate Versions using SQL"
description: "Allows users to query OCI Certificates Management Certificate Versions."
---

# Table: oci_certificates_management_certificate_version - Query OCI Certificates Management Certificate Versions using SQL

The Oracle Cloud Infrastructure (OCI) Certificates Management service helps you manage SSL/TLS certificates for your services. It provides a centralized way to upload and manage certificates, including the ability to monitor certificates for expiration. The Certificates Management service supports both Oracle-managed certificates and customer-managed certificates.

## Table Usage Guide

The `oci_certificates_management_certificate_version` table provides insights into each version of a certificate in OCI Certificates Management. As a security administrator, you can explore certificate-specific details through this table, including certificate data, creation time, and associated metadata. Use it to uncover information about certificates, such as their validity period, the entities they are issued to, and the entities they are issued by.

## Examples

### Basic info
Explore the details of each certificate version in your Oracle Cloud Infrastructure, including its stages, validity, and revocation status. This can be useful in managing your certificates and ensuring they are up to date and secure.

```sql+postgres
select
  certificate_id,
  version_number,
  stages,
  serial_number,
  issuer_ca_version_number,
  version_name,
  subject_alternative_names,
  time_of_deletion,
  validity,
  revocation_status
from
  oci_certificates_management_certificate_version;
```

```sql+sqlite
select
  certificate_id,
  version_number,
  stages,
  serial_number,
  issuer_ca_version_number,
  version_name,
  subject_alternative_names,
  time_of_deletion,
  validity,
  revocation_status
from
  oci_certificates_management_certificate_version;
```

### Get all certificate versions
Discover the different versions of your certificates, including details such as stages, validity, and revocation status. This query is useful for maintaining an overview of your certificate statuses and identifying any potential issues or changes that might impact your system's security.

```sql+postgres
select
  cmcv.certificate_id,
  cmcv.version_number,
  cmcv.stages,
  cmcv.serial_number,
  cmcv.issuer_ca_version_number,
  cmcv.version_name,
  cmcv.subject_alternative_names,
  cmcv.time_of_deletion,
  cmcv.validity,
  cmcv.revocation_status
from
  oci_certificates_management_certificate_version cmcv
  inner join
    oci_certificates_management_certificate cmc
    on cmcv.certificate_id = cmc.id;
```

```sql+sqlite
select
  cmcv.certificate_id,
  cmcv.version_number,
  cmcv.stages,
  cmcv.serial_number,
  cmcv.issuer_ca_version_number,
  cmcv.version_name,
  cmcv.subject_alternative_names,
  cmcv.time_of_deletion,
  cmcv.validity,
  cmcv.revocation_status
from
  oci_certificates_management_certificate_version cmcv
  inner join
    oci_certificates_management_certificate cmc
    on cmcv.certificate_id = cmc.id;
```

### Count versions by certificate
Explore the number of versions associated with each certificate to effectively manage and track certificate updates. This can be particularly useful in maintaining security standards and ensuring optimal system performance.

```sql+postgres
select
  certificate_id,
  count(version_number)
from
  oci_certificates_management_certificate_version
group by
  certificate_id;
```

```sql+sqlite
select
  certificate_id,
  count(version_number)
from
  oci_certificates_management_certificate_version
group by
  certificate_id;
```

### List certificate versions created in the last 30 days
Explore which certificate versions were created in the last month. This can be helpful in tracking recent changes and additions to your certificates, ensuring you're up-to-date on your security configurations.

```sql+postgres
select
  certificate_id,
  version_number,
  time_of_deletion,
  time_created,
  serial_number
from
  oci_certificates_management_certificate_version
where
  time_created >= now() - interval '30' day;
```

```sql+sqlite
select
  certificate_id,
  version_number,
  time_of_deletion,
  time_created,
  serial_number
from
  oci_certificates_management_certificate_version
where
  time_created >= datetime('now', '-30 day');
```
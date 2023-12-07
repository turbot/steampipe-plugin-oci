---
title: "Steampipe Table: oci_certificates_management_certificate_authority_version - Query OCI Certificates Management Certificate Authority Versions using SQL"
description: "Allows users to query Certificate Authority Versions in the OCI Certificates Management service."
---

# Table: oci_certificates_management_certificate_authority_version - Query OCI Certificates Management Certificate Authority Versions using SQL

A Certificate Authority Version in OCI Certificates Management is a specific version of a certificate authority (CA). It contains the information about the CA, such as its public key and the associated private key. This CA version can be used to sign certificates or certificate revocation lists (CRLs).

## Table Usage Guide

The `oci_certificates_management_certificate_authority_version` table provides insights into the versions of certificate authorities within OCI Certificates Management. As a security administrator, you can explore CA version-specific details through this table, including the public key, private key, and other associated metadata. Utilize it to uncover information about CA versions, such as their signing status, the associated CAs, and the validity of each version.

## Examples

### Basic info
Explore the versions of certificate authorities to understand their revocation status, validity, and deletion times. This can be useful for maintaining security compliance and tracking changes in your certificate authorities.

```sql+postgres
select
  certificate_authority_id,
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
  oci_certificates_management_certificate_authority_version;
```

```sql+sqlite
select
  certificate_authority_id,
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
  oci_certificates_management_certificate_authority_version;
```

### Get all certificate authority versions
Determine the versions of all certificate authorities to assess their validity, revocation status, and other key details. This is useful for maintaining the security of your network by ensuring all certificates are up-to-date and valid.

```sql+postgres
select
  cmcav.certificate_authority_id,
  cmcav.version_number,
  cmcav.stages,
  cmcav.serial_number,
  cmcav.issuer_ca_version_number,
  cmcav.version_name,
  cmcav.subject_alternative_names,
  cmcav.time_of_deletion,
  cmcav.validity,
  cmcav.revocation_status
from
  oci_certificates_management_certificate_authority_version cmcav
  inner join
    oci_certificates_management_certificate_authority cmca
    on cmca.id = cmcav.certificate_authority_id;
```

```sql+sqlite
select
  cmcav.certificate_authority_id,
  cmcav.version_number,
  cmcav.stages,
  cmcav.serial_number,
  cmcav.issuer_ca_version_number,
  cmcav.version_name,
  cmcav.subject_alternative_names,
  cmcav.time_of_deletion,
  cmcav.validity,
  cmcav.revocation_status
from
  oci_certificates_management_certificate_authority_version cmcav
  inner join
    oci_certificates_management_certificate_authority cmca
    on cmca.id = cmcav.certificate_authority_id;
```

### Count certificate versions by certificate authority
Analyze the distribution of certificate versions across different certificate authorities. This can be useful in identifying authorities that have a high number of versions, potentially indicating frequent changes or updates.

```sql+postgres
select
  certificate_authority_id,
  count(version_number)
from
  oci_certificates_management_certificate_authority_version
group by
  certificate_authority_id;
```

```sql+sqlite
select
  certificate_authority_id,
  count(version_number)
from
  oci_certificates_management_certificate_authority_version
group by
  certificate_authority_id;
```

### List certificate versions created in the last 30 days
Identify recent certificate versions made within the past month. This can be useful for tracking the creation and management of certificate authorities over time.

```sql+postgres
select
  certificate_authority_id,
  version_number,
  serial_number,
  issuer_ca_version_number,
  version_name
from
  oci_certificates_management_certificate_authority_version
where
  time_created >= now() - interval '30' day;
```

```sql+sqlite
select
  certificate_authority_id,
  version_number,
  serial_number,
  issuer_ca_version_number,
  version_name
from
  oci_certificates_management_certificate_authority_version
where
  time_created >= datetime('now', '-30 day');
```
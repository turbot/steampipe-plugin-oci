---
title: "Steampipe Table: oci_certificates_management_certificate_authority - Query OCI Certificate Management Certificate Authorities using SQL"
description: "Allows users to query Certificate Authorities within the OCI Certificate Management service."
---

# Table: oci_certificates_management_certificate_authority - Query OCI Certificate Management Certificate Authorities using SQL

The Oracle Cloud Infrastructure (OCI) Certificate Management service is a scalable and secure solution for managing digital certificates. It offers a centralized way to create, deploy, and renew certificates, helping to ensure the security of web services and applications. A Certificate Authority (CA) in OCI Certificate Management is a trusted entity that issues digital certificates, which are data files used to cryptographically link an entity with a public key.

## Table Usage Guide

The `oci_certificates_management_certificate_authority` table provides insights into Certificate Authorities within the OCI Certificate Management service. As a security engineer, you can use this table to explore details about each Certificate Authority, including its status, type, and associated endpoints. This can be particularly useful for maintaining an up-to-date inventory of Certificate Authorities, ensuring they are valid and appropriately configured, and identifying any potential security risks.

## Examples

### Basic info
Explore the details of your organization's digital certificates to understand their current status and configuration. This can be useful for maintaining security standards and ensuring proper certificate management.

```sql+postgres
select
  id,
  name,
  config_type,
  issuer_certificate_authority_id,
  description,
  time_of_deletion,
  kms_key_id,
  lifecycle_details,
  certificate_authority_rules,
  current_version,
  certificate_revocation_list_details,
  subject,
  signing_algorithm,
  lifecycle_state as state
from
  oci_certificates_management_certificate_authority;
```

```sql+sqlite
select
  id,
  name,
  config_type,
  issuer_certificate_authority_id,
  description,
  time_of_deletion,
  kms_key_id,
  lifecycle_details,
  certificate_authority_rules,
  current_version,
  certificate_revocation_list_details,
  subject,
  signing_algorithm,
  lifecycle_state as state
from
  oci_certificates_management_certificate_authority;
```

### List inactive certificate authorities
Explore which Certificate Authorities are currently inactive. This is useful to identify potential areas of your system that may lack the necessary security due to inactive authorities.

```sql+postgres
select
  id,
  name,
  lifecycle_state,
  issuer_certificate_authority_id,
  current_version
from
  oci_certificates_management_certificate_authority
where
  lifecycle_state <> 'ACTIVE'
```

```sql+sqlite
select
  id,
  name,
  lifecycle_state,
  issuer_certificate_authority_id,
  current_version
from
  oci_certificates_management_certificate_authority
where
  lifecycle_state <> 'ACTIVE'
```

### List certificates with SHA256_WITH_RSA signed algorithm
Explore which certificates are using the SHA256_WITH_RSA signing algorithm. This is particularly useful for determining the security level of your certificates and identifying any potential risks.

```sql+postgres
select
  id,
  name,
  current_version,
  signing_algorithm,
  certificate_authority_rules
from
  oci_certificates_management_certificate_authority
where
  signing_algorithm = 'SHA256_WITH_RSA';
```

```sql+sqlite
select
  id,
  name,
  current_version,
  signing_algorithm,
  certificate_authority_rules
from
  oci_certificates_management_certificate_authority
where
  signing_algorithm = 'SHA256_WITH_RSA';
```

### List certificates created in the last 30 days
Explore which certificates have been created in the past month. This can be useful in managing and tracking newly issued certificates to ensure they are properly configured and authorized.

```sql+postgres
select
  name,
  id,
  time_created,
  lifecycle_state,
  config_type,
  issuer_certificate_authority_id,
  kms_key_id
from
  oci_certificates_management_certificate_authority
where
  time_created >= now() - interval '30' day;
```

```sql+sqlite
select
  name,
  id,
  time_created,
  lifecycle_state,
  config_type,
  issuer_certificate_authority_id,
  kms_key_id
from
  oci_certificates_management_certificate_authority
where
  time_created >= datetime('now', '-30 day');
```

### Get KMS key details of each certificate
Analyze the settings to understand the relationship between each certificate and its corresponding key in Oracle Cloud Infrastructure's Key Management service. This can help in identifying the certificates that are linked to a particular key, thereby aiding in key management and security compliance.

```sql+postgres
select
  a.id,
  a.name,
  a.kms_key_id,
  k.name as key_name,
  k.vault_id,
  k.current_key_version
from
  oci_certificates_management_certificate_authority as a,
  oci_kms_key as k
where
  k.id = a.kms_key_id;
```

```sql+sqlite
select
  a.id,
  a.name,
  a.kms_key_id,
  k.name as key_name,
  k.vault_id,
  k.current_key_version
from
  oci_certificates_management_certificate_authority as a,
  oci_kms_key as k
where
  k.id = a.kms_key_id;
```

### Get certificate authority rule details
Explore the rules of your certificate authority to understand the maximum validity duration for both the certificate authority itself and the leaf certificates. This helps in managing the lifecycle of your certificates, ensuring they remain valid and secure.

```sql+postgres
select
  name,
  id,
  r ->> 'certificateAuthorityMaxValidityDuration' as certificate_authority_max_validity_duration,
  r ->> 'leafCertificateMaxValidityDuration' as leaf_certificate_max_validity_duration,
  r ->> 'ruleType' as rule_type
from
  oci_certificates_management_certificate_authority,
  jsonb_array_elements(certificate_authority_rules) as r;
```

```sql+sqlite
select
  name,
  id,
  json_extract(r.value, '$.certificateAuthorityMaxValidityDuration') as certificate_authority_max_validity_duration,
  json_extract(r.value, '$.leafCertificateMaxValidityDuration') as leaf_certificate_max_validity_duration,
  json_extract(r.value, '$.ruleType') as rule_type
from
  oci_certificates_management_certificate_authority,
  json_each(certificate_authority_rules) as r;
```

### List certificates that are valid upto a certain date
Identify certificates that will remain valid until a specific date. This is useful for planning renewals and managing certificate lifecycles.

```sql+postgres
select
  name,
  id,
  validity ->> 'timeOfValidityNotAfter' as time_of_validity_not_after
from
  oci_certificates_management_certificate_authority
where
  time_of_validity_not_after <= '2033-06-04T00:00:00Z';
```

```sql+sqlite
select
  name,
  id,
  json_extract(validity, '$.timeOfValidityNotAfter') as time_of_validity_not_after
from
  oci_certificates_management_certificate_authority
where
  json_extract(validity, '$.timeOfValidityNotAfter') <= '2033-06-04T00:00:00Z';
```
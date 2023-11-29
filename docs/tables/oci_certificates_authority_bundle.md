---
title: "Steampipe Table: oci_certificates_authority_bundle - Query OCI Identity Certificate Authority Bundle using SQL"
description: "Allows users to query Certificate Authority Bundles in Oracle Cloud Infrastructure."
---

# Table: oci_certificates_authority_bundle - Query OCI Identity Certificate Authority Bundle using SQL

Oracle Cloud Infrastructure (OCI) provides a Certificate Authority Bundle, which is a collection of root and intermediate certificates. These certificates are used by OCI services and resources for secure communication. They are essential for establishing trusted connections between different parts of the OCI ecosystem.

## Table Usage Guide

The `oci_certificates_authority_bundle` table provides insights into the Certificate Authority Bundles within Oracle Cloud Infrastructure Identity service. As a security engineer, explore details of these bundles through this table, including their contents and associated metadata. Utilize it to uncover information about the trust relationships within OCI, and to verify the authenticity of OCI services and resources.

## Examples

### Basic info
Determine the status and details of a specific certificate authority in a specific version. This is useful for validating the authenticity of the certificate authority and ensuring it's not revoked, which is crucial for maintaining secure connections.

```sql
select
  certificate_authority_id,
  certificate_authority_name,
  serial_number,
  certificate_pem,
  version_number,
  validity,
  stages,
  cert_chain_pem,
  version_name,
  revocation_status
from
  oci_certificates_authority_bundle
where
  certificate_authority_id = 'ocid1.certificateauthority.oc1.ap-mumbai-1.amaaaaaa6igdexaatxksy32ktbtblp6knkei4xr2kl2zc46sdpxujs22momq'
  and
  version_number = 1;
```

### Get all certificate authority bundles
Explore the details of all certificate authority bundles to understand their validity and revocation status. This is useful for ensuring the security and integrity of your digital certificates.
```sql
select
  ccab.certificate_authority_id,
  ccab.certificate_authority_name,
  ccab.serial_number,
  ccab.certificate_pem,
  ccab.version_number,
  ccab.validity,
  ccab.stages,
  ccab.cert_chain_pem,
  ccab.version_name,
  ccab.revocation_status
from
    oci_certificates_authority_bundle ccab
inner join oci_certificates_management_certificate_authority cmca
on ccab.certificate_authority_id = cmca.id;
```
# Table: oci_certificates_management_ca_bundle

OCI Certificate Management is a service provided by Oracle Cloud Infrastructure (OCI) that enables you to manage and secure your SSL/TLS certificates. The OCI Certificate Management CA bundle, also known as the Root CA Bundle, is a collection of trusted root certificates provided by OCI. These certificates are used to validate the trustworthiness of SSL/TLS certificates issued by certificate authorities supported by OCI Certificate Management.

## Examples

### Basic info

```sql
select
  id,
  name,
  lifecycle_state,
  lifecycle_details,
  description,
  time_created
from
  oci_certificates_management_ca_bundle;
```

### List bundles created between a specific time

```sql
select
  name,
  id,
  lifecycle_state,
  time_created
from
  oci_certificates_management_ca_bundle
where
  time_created between '2023-05-01' and '2023-07-01';
```

### List failed bundles

```sql
select
  name,
  id,
  lifecycle_state,
  lifecycle_details
from
  oci_certificates_management_ca_bundle
where
  lifecycle_state = 'FAILED';
```
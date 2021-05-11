select retention_period_days, tenant_id
from oci.oci_audit_configuration
where tenant_id = '{{ output.tenancy_ocid.value }}::dummy';
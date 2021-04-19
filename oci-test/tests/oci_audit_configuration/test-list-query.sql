select retention_period_days,tenant_id
from oci.oci_audit_configuration
where retention_period_days = '{{ output.retention_period_days.value }}';
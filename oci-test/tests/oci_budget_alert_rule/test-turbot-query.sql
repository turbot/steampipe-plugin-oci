select title, tenant_id
from oci.oci_budget_alert_rule
where id = '{{ output.resource_id.value }}';
select id, display_name, lifecycle_state
from oci.oci_budget_alert_rule
where id = '{{ output.resource_id.value }}' and budget_id = '{{ output.budget_id.value }}';
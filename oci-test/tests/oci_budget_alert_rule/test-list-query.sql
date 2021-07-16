select id, display_name, lifecycle_state
from oci.oci_budget_alert_rule
where display_name = '{{ resourceName }}';
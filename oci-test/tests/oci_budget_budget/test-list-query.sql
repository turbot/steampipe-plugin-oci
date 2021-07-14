select id, display_name, lifecycle_state
from oci.oci_budget_budget
where display_name = '{{ resourceName }}';
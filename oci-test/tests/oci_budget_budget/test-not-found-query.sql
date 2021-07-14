select id, lifecycle_state
from oci.oci_budget_budget
where id = '{{ output.resource_id.value }}::dummy';
select title, tenant_id
from oci.oci_budget_budget
where id = '{{ output.resource_id.value }}';
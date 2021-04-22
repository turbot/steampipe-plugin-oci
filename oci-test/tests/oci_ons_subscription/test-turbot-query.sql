select tenant_id, region
from oci.oci_ons_subscription
where id = '{{ output.resource_id.value }}';
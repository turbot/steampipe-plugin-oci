select title, tenant_id
from oci.oci_autoscaling_auto_scaling_configuration
where id = '{{ output.resource_id.value }}';
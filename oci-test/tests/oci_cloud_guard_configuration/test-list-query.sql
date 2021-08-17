select reporting_region, self_manage_resources
from oci.oci_cloud_guard_configuration
where reporting_region = '{{ output.reporting_region.value }}';
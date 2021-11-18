select reporting_region
from oci.oci_cloud_guard_configuration
where reporting_region = '{{ output.reporting_region.value }}::dummy';
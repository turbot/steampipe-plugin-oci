select region
from oci.oci_cloud_guard_configuration
where region = '{{ output.reporting_region.value }}::dummy';
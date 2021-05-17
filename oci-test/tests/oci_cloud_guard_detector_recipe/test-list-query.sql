select id, name, lifecycle_state
from oci.oci_cloud_guard_detector_recipe
where name = '{{ resourceName }}';
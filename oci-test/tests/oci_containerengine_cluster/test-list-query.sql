select name, id, lifecycle_state
from oci.oci_containerengine_cluster
where name = '{{ resourceName }}';
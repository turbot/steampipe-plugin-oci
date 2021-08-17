select id, name, lifecycle_state, capacity_type, capacity_value, feature_set, license_type
from oci.oci_analytics_instance
where name = '{{ resourceName }}';
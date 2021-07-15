select
  name,
  object_lifecycle_policy -> 'items' as object_lifecycle_policy_rules
from
  oci.oci_objectstorage_bucket
where
  name = '{{ resourceName }}';
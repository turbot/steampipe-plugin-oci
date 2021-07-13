select title, tenant_id
from oci.oci_objectstorage_bucket
where name = '{{ resourceName }}';
select name, namespace, id, freeform_tags, versioning
from oci.oci_objectstorage_bucket
where id = '{{ output.resource_id.value }}';
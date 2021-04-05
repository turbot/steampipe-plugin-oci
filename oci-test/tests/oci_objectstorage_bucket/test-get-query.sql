select name, namespace, id, freeform_tags, versioning
from oci.oci_objectstorage_bucket
where name = '{{ resourceName }}';
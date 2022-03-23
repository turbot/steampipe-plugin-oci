select name, id, lifecycle_state
from oci.oci_streaming_stream
where name = '{{ resourceName }}';
# Table: oci_objectstorage_object

In the Oracle Cloud Infrastructure Object Storage service, an object is a file or unstructured data you upload to a bucket within a compartment  within an Object Storage namespace. The object can be any type of data, for example, multimedia files, data backups, static web content, or logs.

## Examples

### Basic info

```sql
select
  name,
  bucket_name,
  namespace,
  region,
  size,
  md5,
  time_created,
  time_modified
from
  oci_objectstorage_object;
```


### List archived objects

```sql
select
  name,
  bucket_name,
  namespace,
  region,
  archival_state
from
  oci_objectstorage_object
where
  archival_state = 'Archived';
```

# Table: oci_releasenote

This table offers the release notes published by Oracle for OCI. These notes describe all new releases, new services and major changes in any of the OCI services - including Cloud Shell, Cloud Advisor and Developer Tools.  

The data is retrieved ("scraped") from the release notes pages on the Oracle website, at https://docs.oracle.com/en-us/iaas/releasenotes/. Retrieving data in this way is a fairly slow process. Queries are faster if you include a limit clause to reduce the number of notes to be considered and therefore the number of pages scraped.

The most recent 50 release notes are published in this RSS Feed: https://docs.oracle.com/en-us/iaas/releasenotes/feed/

## Examples

### Basic info - list most recent 100 OCI release notes 

```sql
select 
  title,
  service, 
  summary, 
  release_date 
from 
  oci_releasenote 
limit 100  
```

### List release notes published in the last 30 days

```sql
select 
  title,
  service, 
  summary, 
  release_date 
from 
  oci_releasenote 
where 
  release_date  > current_date - interval '1' month;
```

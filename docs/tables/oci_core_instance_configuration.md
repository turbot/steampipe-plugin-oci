---
title: "Steampipe Table: oci_core_instance_configuration - Query OCI Core Instance Configurations using SQL"
description: "Allows users to query OCI Core Instance Configurations."
---

# Table: oci_core_instance_configuration - Query OCI Core Instance Configurations using SQL

An OCI Core Instance Configuration is a template that contains the settings to launch a compute instance in the Oracle Cloud Infrastructure Compute service. These configurations capture the instance's metadata, including the OCID, the compartment OCID, and the display name. They also provide details about the instance's launch parameters, such as the instance type, the shape, and the associated networking details.

## Table Usage Guide

The `oci_core_instance_configuration` table provides insights into the instance configurations within OCI's Compute service. As a cloud engineer, explore instance-specific details through this table, including instance type, shape, and associated networking details. Utilize it to uncover information about configurations, such as those associated with specific compartments, the metadata of instances, and the verification of launch parameters.

## Examples

### Basic info
Explore the general details of your OCI instance configurations to understand when and where they were created. This can aid in managing resources and tracking the deployment of instances across different regions.

```sql
select
  display_name,
  id,
  time_created,
  region
from
  oci_core_instance_configuration;
```

### List instance configurations along with the compartment details
Explore the configurations of various instances and their associated compartments to understand how resources are allocated and organized within different regions. This information can be useful for assessing resource distribution and planning for future infrastructure needs.

```sql
select
  inst.display_name,
  inst.id,
  inst.region,
  comp.name as compartment_name
from
  oci_core_instance_configuration inst
  inner join
    oci_identity_compartment comp
    on (inst.compartment_id = comp.id)
order by
  comp.name,
  inst.region;
```

### List instance configurations created in the last 30 days
Discover the segments that have been recently added in the past month. This query is useful for keeping track of newly created configurations, aiding in system management and monitoring.

```sql
select
  display_name,
  id,
  time_created,
  region
from
  oci_core_instance_configuration
where
  time_created >= now() - interval '30' day;
```
import pulumi_aws as aws

from config import APP_NAME, REDIS_CLUSTER_SIZE, REDIS_CLUSTER_NAME, REDIS_CLUSTER_PORT
from networking import default_vpc, default_vpc_subnets

ENGINE_VERSION = "7.0"

PARAMETER_GROUP_FAMILY = "redis7"

NODE_TYPE = "cache.t2.micro"

# Create a Security Group
db_security_group = aws.ec2.SecurityGroup(
    f"{APP_NAME}-redis-secgrp",
    vpc_id=default_vpc.id,
    ingress=[
        aws.ec2.SecurityGroupIngressArgs(
            protocol="tcp",
            from_port=REDIS_CLUSTER_PORT,
            to_port=REDIS_CLUSTER_PORT,
            cidr_blocks=["0.0.0.0/0"],
        ),
    ],
    egress=[
        aws.ec2.SecurityGroupEgressArgs(
            protocol="-1",
            from_port=0,
            to_port=0,
            cidr_blocks=["0.0.0.0/0"],
        )
    ],
)

# # Create an ElastiCache Subnet Group
subnet_group = aws.elasticache.SubnetGroup(
    f"{APP_NAME}-redis-subnetgrp", subnet_ids=default_vpc_subnets.ids
)

redis_cluster = aws.elasticache.ReplicationGroup(
    REDIS_CLUSTER_NAME,
    replication_group_id=f"{REDIS_CLUSTER_NAME}-group",
    description="TinyURL Redis cluster",
    engine="redis",
    node_type=NODE_TYPE,
    num_cache_clusters=REDIS_CLUSTER_SIZE,
    engine_version=ENGINE_VERSION,
    parameter_group_name=f"default.{PARAMETER_GROUP_FAMILY}",
    subnet_group_name=subnet_group.name,
    security_group_ids=[db_security_group.id],
)

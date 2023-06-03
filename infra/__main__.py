import json

from pulumi import export, ResourceOptions
import pulumi_aws as aws

from config import APP_NAME, DOCKER_IMAGE, APP_PORT
from networking import default_vpc_subnets, security_group, lb_target_group, load_balancer, cert, http_list, https_list

# Create an ECS cluster to run a container-based service
cluster = aws.ecs.Cluster(f'{APP_NAME}-cluster')

# Create an IAM role that can be used by our service's task
role = aws.iam.Role(f'{APP_NAME}-task-exec-role',
	assume_role_policy=json.dumps({
		'Version': '2008-10-17',
		'Statement': [{
			'Sid': '',
			'Effect': 'Allow',
			'Principal': {
				'Service': 'ecs-tasks.amazonaws.com'
			},
			'Action': 'sts:AssumeRole',
		}]
	}),
)

rpa = aws.iam.RolePolicyAttachment(f'{APP_NAME}-task-exec-policy',
	role=role.name,
	policy_arn='arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy',
)

# Spin up a load balanced service running our container image
task_definition = aws.ecs.TaskDefinition(f'{APP_NAME}-task',
    family='fargate-task-definition',
    cpu='512',
    memory='1024',
    network_mode='awsvpc',
    requires_compatibilities=['FARGATE'],
    execution_role_arn=role.arn,
    container_definitions=json.dumps([{
		'name': APP_NAME,
		'image': DOCKER_IMAGE,
		'portMappings': [{
			'containerPort': APP_PORT,
			'hostPort': APP_PORT,
			'protocol': 'tcp'
		}]
	}])
)

service = aws.ecs.Service(f'{APP_NAME}-service',
	cluster=cluster.arn,
    desired_count=1,
    launch_type='FARGATE',
    task_definition=task_definition.arn,
    network_configuration=aws.ecs.ServiceNetworkConfigurationArgs(
		assign_public_ip=True,
		subnets=default_vpc_subnets.ids,
		security_groups=[security_group.id],
	),
    load_balancers=[aws.ecs.ServiceLoadBalancerArgs(
		target_group_arn=lb_target_group.arn,
		container_name=APP_NAME,
		container_port=APP_PORT,
	)],
    opts=ResourceOptions(depends_on=[http_list, https_list]),
)

export('url', load_balancer.dns_name)
export("certificateArn", cert.arn)

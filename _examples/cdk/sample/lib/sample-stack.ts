import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";

export class Function extends Construct {
  readonly alias: cdk.aws_lambda.Alias;

  constructor(scope: Construct, id: string) {
    super(scope, id);

    const role = new cdk.aws_iam.Role(this, "Role", {
      assumedBy: new cdk.aws_iam.ServicePrincipal("lambda.amazonaws.com"),
      inlinePolicies: {
        LoggingPolicy: new cdk.aws_iam.PolicyDocument({
          statements: [
            new cdk.aws_iam.PolicyStatement({
              effect: cdk.aws_iam.Effect.ALLOW,
              actions: [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents",
              ],
              resources: ["arn:aws:logs:*:*:*"],
            }),
          ],
        }),
        FunctionPolicy: new cdk.aws_iam.PolicyDocument({
          statements: [
            new cdk.aws_iam.PolicyStatement({
              effect: cdk.aws_iam.Effect.ALLOW,
              actions: [
                "logs:DescribeLogGroups",
                "logs:PutRetentionPolicy",
                "logs:DeleteRetentionPolicy",
                "logs:DeleteLogGroup",
              ],
              resources: ["arn:aws:logs:*:*:*"],
            }),
          ],
        }),
      },
    });

    const fn = new cdk.aws_lambda.DockerImageFunction(this, "Function", {
      description:
        "Set retention period for log groups with no expiration date.",
      code: cdk.aws_lambda.DockerImageCode.fromImageAsset("src/lambda"),
      architecture: cdk.aws_lambda.Architecture.ARM_64,
      role: role,
      logRetention: cdk.aws_logs.RetentionDays.THREE_MONTHS,
      currentVersionOptions: {
        removalPolicy: cdk.RemovalPolicy.RETAIN,
      },
      timeout: cdk.Duration.minutes(5),
      environment: {
        FILTERS: "retention == infinite",
        DESIRED_STATE: "3months",
      },
    });
    this.alias = new cdk.aws_lambda.Alias(this, "Alias", {
      aliasName: "live",
      version: fn.currentVersion,
    });
  }
}

export interface ScheduleProps {
  alias: cdk.aws_lambda.Alias;
}

export class Schedule extends Construct {
  constructor(scope: Construct, id: string, props: ScheduleProps) {
    super(scope, id);

    const role = new cdk.aws_iam.Role(this, "Role", {
      assumedBy: new cdk.aws_iam.ServicePrincipal("scheduler.amazonaws.com"),
    });

    new cdk.aws_scheduler.CfnSchedule(this, "Schedule", {
      flexibleTimeWindow: {
        mode: "OFF",
      },
      scheduleExpression: "cron(0 0 1 * ? *)",
      target: {
        arn: props.alias.functionArn,
        roleArn: role.roleArn,
        input: undefined,
      },
    });
  }
}

export class SampleStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const fn = new Function(this, "Function");
    new Schedule(this, "Schedule", {
      alias: fn.alias,
    });
  }
}

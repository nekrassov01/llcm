#!/usr/bin/env node
import * as cdk from "aws-cdk-lib";
import { SampleStack } from "../lib/sample-stack";

const app = new cdk.App();
new SampleStack(app, "SampleStack", {
  env: {
    region: process.env.CDK_DEFAULT_REGION,
    account: process.env.CDK_DEFAULT_ACCOUNT,
  },
});

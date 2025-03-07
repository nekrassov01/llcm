import * as cdk from "aws-cdk-lib";
import { Template } from "aws-cdk-lib/assertions";
import { SampleStack } from "../lib/sample-stack";

const app = new cdk.App();
const stack = new SampleStack(app, "SampleStack");
const template = Template.fromStack(stack);

describe("Snapshot Tests", () => {
  test("Snapshot test", () => {
    expect(template.toJSON()).toMatchSnapshot();
  });
});

describe("Fine-grained Assertions Tests", () => {
  test("Lambda Function created", () => {
    template.resourceCountIs("AWS::Lambda::Function", 2); // including the custom resource
  });
  test("Lambda Version created", () => {
    template.resourceCountIs("AWS::Lambda::Version", 1);
  });
  test("Lambda Alias created", () => {
    template.resourceCountIs("AWS::Lambda::Alias", 1);
  });
  test("Scheduler Schedule created", () => {
    template.resourceCountIs("AWS::Scheduler::Schedule", 1);
  });
  test("IAM Role created", () => {
    template.resourceCountIs("AWS::IAM::Role", 3); // including the custom resource
  });
});

# llcm

[![CI](https://github.com/nekrassov01/llcm/actions/workflows/ci.yml/badge.svg)](https://github.com/nekrassov01/llcm/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nekrassov01/llcm)](https://goreportcard.com/report/github.com/nekrassov01/llcm)
![GitHub](https://img.shields.io/github/license/nekrassov01/llcm)
![GitHub](https://img.shields.io/github/v/release/nekrassov01/llcm)

llcm is a CLI tool to manage the lifecycle of Amazon CloudWatch Logs log groups. List, update, and delete log groups to manage their lifecycle. It handles multiple regions at high speed while avoiding throttling errors. It can also return simulation results based on desired states.

## Features

- **List**: Fast listing of log groups for specified multiple regions.
- **Preview**: By passing the desired state as an argument, the log group is listed with the results of the reduction simulation.
- **Apply**: The desired state passed in the argument is actually applied to the listed log groups.

All of these subcommands can be passed the filter expressions to narrow down the target log groups.

## Commands

```text
NAME:
   llcm - AWS Log groups lifecycle manager

USAGE:
   llcm [global options] command [command options]

VERSION:
   0.0.0

DESCRIPTION:
   A listing, updating, and deleting tool to manage the lifecycle of Amazon CloudWatch Logs.
   It handles multiple regions fast while avoiding throttling. It can also return simulation
   results based on the desired state.

COMMANDS:
   list        List log group entries with specified format
   preview     Preview simulation results based on desired state
   apply       Apply desired state to log group entries

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### List

```text
NAME:
   llcm list - List log group entries with specified format

USAGE:
   llcm list [command options]

DESCRIPTION:
   List collects basic information about log groups from multiple specified regions and
   returns it in a specified format.

OPTIONS:
   --profile value, -p value                              set aws profile [$AWS_PROFILE]
   --log-level value, -l value                            set log level (default: "info") [$LLCM_LOG_LEVEL]
   --region value, -r value [ --region value, -r value ]  set target regions (default: all regions with no opt-in)
   --filter value, -f value [ --filter value, -f value ]  set expressions to filter log groups
   --output value, -o value                               set output type (default: "compressedtext") [$LLCM_OUTPUT_TYPE]
   --help, -h                                             show help
```

### Preview

```text
NAME:
   llcm preview - Preview simulation results based on desired state

USAGE:
   llcm preview [command options]

DESCRIPTION:
   Preview performs a simple calculation based on `DesiredState` specified in the argument
   and returns a simulated list including `ReducibleBytes`, `RemainingBytes`, etc.

OPTIONS:
   --profile value, -p value                              set aws profile [$AWS_PROFILE]
   --log-level value, -l value                            set log level (default: "info") [$LLCM_LOG_LEVEL]
   --region value, -r value [ --region value, -r value ]  set target regions (default: all regions with no opt-in)
   --filter value, -f value [ --filter value, -f value ]  set expressions to filter log groups
   --desired value, -d value                              set the desired state
   --output value, -o value                               set output type (default: "compressedtext") [$LLCM_OUTPUT_TYPE]
   --help, -h                                             show help
```

### Apply

```text
NAME:
   llcm apply - Apply desired state to log group entries

USAGE:
   llcm apply [command options]

DESCRIPTION:
   Apply deletes and updates target log groups in batches based on `DesiredState`.
   It is fast across multiple regions, but cleverly avoids throttling.

OPTIONS:
   --profile value, -p value                              set aws profile [$AWS_PROFILE]
   --log-level value, -l value                            set log level (default: "info") [$LLCM_LOG_LEVEL]
   --region value, -r value [ --region value, -r value ]  set target regions (default: all regions with no opt-in)
   --filter value, -f value [ --filter value, -f value ]  set expressions to filter log groups
   --desired value, -d value                              set the desired state
   --help, -h                                             show help
```

## Options

The following values can be passed for each option.

| Option                                            | Values                                                                                                                                                                                                                                                                                                                                                                                                                                                                | Default value                                                                                                                             | Environment Variable |
| ------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------- | -------------------- |
| `--profile value` `-p value`                      | -                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | -                                                                                                                                         | `AWS_PROFILE`        |
| `--log-level value` `-l value`                    | `debug` `info` `warn` `error`                                                                                                                                                                                                                                                                                                                                                                                                                                         | `info`                                                                                                                                    | `LLCM_LOG_LEVEL`     |
| `--region value1,value2...` `-r value1,value2...` | `af-south-1` `ap-east-1` `ap-northeast-1` `ap-northeast-2` `ap-northeast-3` `ap-south-1` `ap-south-2` `ap-southeast-1` `ap-southeast-2` `ap-southeast-3` `ap-southeast-4` `ap-southeast-5` `ap-southeast-7` `ca-central-1` `ca-west-1` `eu-central-1` `eu-central-2` `eu-north-1` `eu-south-1` `eu-south-2` `eu-west-1` `eu-west-2` `eu-west-3` `il-central-1` `me-central-1` `me-south-1` `mx-central-1` `sa-east-1` `us-east-1` `us-east-2` `us-west-1` `us-west-2` | [All regions with no opt-in](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html#concepts-regionsz) | -                    |
| `--filter value1,value2...` `-f value1,value2...` | key, operator, and value separated by spaces, as in `bytes != 0`<br>key: `name` `class` `elapsed` `retention` `bytes`<br>operator: `>` `>=` `<` `<=` `==` `==*` `!=` `!=*` `=~` `=~*` `!~` `!~*`                                                                                                                                                                                                                                                                      | -                                                                                                                                         | -                    |
| `--desired value` `-d value`                      | `delete` `1day` `3days` `5days` `1week` `2weeks` `1month` `2months` `3months` `4months` `5months` `6months` `1year` `13months` `18months` `2years` `3years` `5years` `6years` `7years` `8years` `9years` `10years` `infinite`                                                                                                                                                                                                                                         | -                                                                                                                                         | -                    |
| `--output value` `-o value`                       | `json` `prettyjson` `text` `compressedtext` `markdown` `backlog` `tsv`                                                                                                                                                                                                                                                                                                                                                                                                | `compressedtext`                                                                                                                          | `LLCM_OUTPUT_TYPE`   |
| `--help` `-h`                                     | -                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | -                                                                                                                                         | -                    |
| `--version` `-v`                                  | -                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | -                                                                                                                                         | -                    |

## Examples

### Case 1

- List logs generated by Lambda from all regions that are not empty and have a retention period longer than one year, and list the reduction effect if the retention period is changed uniformly to one year. The output shall be in markdown table format.

> [!NOTE]
> Note that this simulation is simply a pro-rata calculation of the log bytes.

```sh
llcm preview --desired 1year --filter 'name =~ ^/aws/lambda/.*','bytes != 0','retention > 1year' --output markdown

# The following outputs are obtained
| Name                 | Region         | Class    | CreatedAt                 | ElapsedDays | RetentionInDays | StoredBytes  | BytesPerDay | DesiredState | ReductionInDays | ReducibleBytes | RemainingBytes |
| -------------------- | -------------- | -------- | ------------------------- | ----------- | --------------- | ------------ | ----------- | ------------ | --------------- | -------------- | -------------- |
| /aws/lambda/tokyo-1  | ap-northeast-1 | STANDARD | 2019-04-15T21:50:12+09:00 | 2107        | 731             | 161094000389 | 220374829   | 365          | 366             | 80657187414    | 80436812975    |
| /aws/lambda/tokyo-2  | ap-northeast-1 | STANDARD | 2020-08-26T23:45:50+09:00 | 1608        | 731             | 30273686566  | 41414071    | 365          | 366             | 15157549986    | 15116136580    |
| /aws/lambda/oregon-1 | us-west-2      | STANDARD | 2020-08-27T14:34:54+09:00 | 1607        | 731             | 28578246408  | 39094728    | 365          | 366             | 14308670448    | 14269575960    |
| /aws/lambda/oregon-2 | us-west-2      | STANDARD | 2020-08-26T23:48:51+09:00 | 1608        | 731             | 22822519036  | 31220956    | 365          | 366             | 11426869896    | 11395649140    |
...
```

- Apply the desired retention period to the log groups identified above.

```sh
llcm apply --desired 1year --filter 'name =~ ^/aws/lambda/.*','bytes != 0','retention > 1year'
```

### Case 2

- List log groups for Tokyo and Oregon that are empty and have been created more than one year ago in the backlog table format.

```sh
llcm list --filter 'bytes == 0','elapsed > 365' --region ap-northeast-1,us-west-2 --output backlog

# The following outputs are obtained
| Name   | Region         | Class    | CreatedAt                 | ElapsedDays | RetentionInDays | StoredBytes |h
| test-1 | ap-northeast-1 | STANDARD | 2017-12-07T13:16:02+09:00 |        2601 |             731 |           0 |
| test-2 | us-west-2      | STANDARD | 2017-12-07T12:44:45+09:00 |        2601 |             731 |           0 |
| test-3 | ap-northeast-1 | STANDARD | 2017-12-07T13:21:09+09:00 |        2601 |             731 |           0 |
| test-4 | us-west-2      | STANDARD | 2017-12-07T12:50:11+09:00 |        2601 |             731 |           0 |
...
```

- Delete the log groups identified above in a batch.

```sh
llcm apply --desired delete --filter 'bytes == 0','elapsed > 365' --region ap-northeast-1,us-west-2
```

### Case 3

- For a more serious investigation, copy all log groups for that account in TSV format to the clipboard and paste them into a spreadsheet.

```sh
# macOS
llcm list --output tsv | pbcopy

# Windows
llcm list --output tsv | Set-Clipboard
```

### Case 4

- Visualize the reductions in order to submit the rate reduction plan in a document. A simple stacked bar chart will appear in your browser in one stroke!

```sh
llcm preview --desired 1year --output chart
```

![chart](_assets/preview.png)

- If the output type is chart in the list command, a pie chart is displayed.

```sh
llcm list --output chart
```

![chart](_assets/list.png)

## Moreover

In addition to use from the command line, it is designed to be used as is as a package; see the example of running it with Lambda defined in the AWS CDK.

- [CDK stack sample](_examples/cdk/sample/lib/sample-stack.ts)
- [Lambda function sample](_examples/cdk/sample/src/lambda/main.go)

## Warnings

- Consider enclosing strings passed to the filter in single quotes. Unintended expansion may occur, e.g., history expansion by the shell (Try typing this command in your shell environment: `echo "name !~ ^test.*"`)
- The Preview command is the best used to simulate reductions, but note that it is only a simple calculation of the log volume pro-rated by day.
- The fields such as `ElapsedDays` and `ReductionInDays` represent the number of days, but are rounded down to the nearest whole number when cast to int64. This means that the reduction simulation will not be inflated beyond what is expected.
- The minimum value for `BytesPerDay` is 1. Note this specification if you have a large number of log groups that have just been created and are small in size.

## Installation

Install with homebrew

```sh
brew install nekrassov01/tap/llcm
```

Install with go

```sh
go install github.com/nekrassov01/llcm@latest
```

Or download binary from [releases](https://github.com/nekrassov01/llcm/releases)

## Shell completion

Supported Shells are as follows:

- bash
- zsh
- pwsh

```sh
llcm completion bash|zsh|pwsh
```

## Todo

- [x] Add to readme an example of implementing as a lambda function
- [ ] Implement logical operators: `&&` `||`
- [x] Implement visualization of simulation results
- [ ] Support streaming output (in `github.com/nekrassov01/mintab`)

## Author

[nekrassov01](https://github.com/nekrassov01)

## License

[MIT](https://github.com/nekrassov01/llcm/blob/main/LICENSE)

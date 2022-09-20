# OnePassENV

OnePassENV is a command proxy tool to get environment variables from one pass and execute the programs with that variables.

It does not store any secret data on the disk, it only reads configuration JSON for profiling, whitelisting for environment variables that come from 1pass, and last thing is the executable path of the op command.

## Installing

### op command

This software requires the 1Password CLI v2 to get secrets.  
You can download it from [site](https://developer.1password.com/docs/cli/get-started/) or you can use below brew command.

```bash
brew install --cask 1password/tap/1password-cli
```

### onepassenv

The prebuild binaries are not available for this project.
You can build with go

```bash
go build
sudo mv onepassenv /usr/local/bin/
sudo chown root:wheel /usr/local/bin/onepassenv
chmod 755 /usr/local/bin/onepassenv
```

Then you need to configure your config file. It contains profile information and op command path.

In the profiles, you will see two parameter, one of the profile name other is environment

```json
{
        "profileName": "aws-dev",
        "variables": ["CDK_DEFAULT_ACCOUNT","AWS_ACCESS_KEY_ID","AWS_SECRET_ACCESS_KEY","AWS_DEFAULT_REGION"]
}
```

In the 1password, you can store many key=value in one profile, to prevent unwanted data leaks, you may want to allow only specific ones, and this is done by setting the variables array in the profiles. So other values in the 1password profile will not pass to the application (such as username and password).

You need to specify which programs can be executed by this project. So any possible mistaken execution will prevented with this parameter.

```json
"allowedbins":[
      "/opt/homebrew/bin/aws",
      "/opt/homebrew/bin/cdk"
    ]
```

## Use

Two parameters are required for execution, to define the profile you need to set onepenv and to specify the executable command you must set onepenvbin, and then you can execute the file.

```bash
onepenv=aws-dev onepenvbin=aws onepassenv sts get-caller-identity --query "Account" --output text | cat
```

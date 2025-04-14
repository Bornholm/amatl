# On Linux

1. Download the latest Linux release: https://github.com/Bornholm/amatl/releases/tag/{{ .Vars.amatlVersion }}

   ```shell
   wget https://github.com/Bornholm/amatl/releases/download/{{ .Vars.amatlVersion }}/amatl_{{ trimPrefix "v" .Vars.amatlVersion }}_linux_amd64.tar.gz
   ```

   _Replace `amd64` by your architecture._

2. Extract the binary then move it to your preferred location, for example `/usr/local/bin`

   ```shell
   tar -xzf amatl_{{ trimPrefix "v" .Vars.amatlVersion }}_linux_amd64.tar.gz
   sudo mv ./amatl /usr/local/bin
   ```

3. Check that your installation is ok by running the `--help` command

   ```shell
    amatl help
   ```

   The command should return something like this:

   ```shell
    NAME:
      amatl - a markdown to markdown/html/pdf compiler

    USAGE:
      amatl [global options] command [command options]

    COMMANDS:
      render
      help, h  Shows a list of commands or help for one command

    GLOBAL OPTIONS:
      --debug            Enable debug mode (default: false) [$AMATL_DEBUG]
      --log-level value  Set logging level (default: "info") [$AMATL_LOG_LEVEL]
      --workdir value    The working directory [$AMATL_WORKDIR]
      --help, -h         show help
   ```

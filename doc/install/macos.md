# 🍎 On MacOS

1. Visit the [latest release page](https://github.com/Bornholm/amatl/releases/tag/{{ .Vars.amatlVersion }}) or use `curl`:

   ```sh
   curl -LO https://github.com/Bornholm/amatl/releases/download/{{ .Vars.amatlVersion }}/amatl_{{ trimPrefix "v" .Vars.amatlVersion }}_darwin_amd64.tar.gz
   ```

   _Replace `amd64` with your CPU architecture (e.g., `arm64` for Apple Silicon)._

2. Extract the archive and move the binary:

   ```sh
   tar -xzf amatl_{{ trimPrefix "v" .Vars.amatlVersion }}_darwin_amd64.tar.gz
   sudo mv ./amatl /usr/local/bin
   ```

3. Verify the installation:

   ```sh
   amatl help
   ```

   You should see the CLI help output, confirming that Amatl is installed.

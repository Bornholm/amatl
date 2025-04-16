# 🪟 On Windows

1. Visit the [latest release page](https://github.com/Bornholm/amatl/releases/tag/{{ .Vars.amatlVersion }}) and download the file named:

   ```
   amatl_{{ trimPrefix "v" .Vars.amatlVersion }}_windows_amd64.tar.gz
   ```

   _Replace `amd64` with your architecture if needed._

2. Uncompress the archive with your preferred archive extractor.

3. Move `amatl.exe` to a directory included in your `PATH` (e.g., `C:\Program Files\Amatl\`) or run it directly from the extracted folder.

4. Open **PowerShell** or **Command Prompt** and run:

   ```
   amatl.exe help
   ```

   You should see the list of commands and global options, confirming the installation was successful.

using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ZmkFlasher.Lib;

internal class Process
{
    public static async Task<string> Run(string command, string arguments)
    {
        var processDescription = new ProcessStartInfo()
        {
            FileName = command,
            Arguments = arguments,
            UseShellExecute = false,
            RedirectStandardOutput = true,
        };
        var process = System.Diagnostics.Process.Start(processDescription);
        if(process is null) throw new Exception("Failed to start process");
        await process.WaitForExitAsync();
        var output = await process.StandardOutput.ReadToEndAsync();
        if (process.ExitCode != 0) throw new Exception($"exit code != 0 ({process.ExitCode}). Output: {output}");
        return output;
    }
}

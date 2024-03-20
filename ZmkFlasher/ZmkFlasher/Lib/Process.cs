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
            RedirectStandardOutput = true,
        };
        var process = System.Diagnostics.Process.Start(processDescription);
        if(process is null) throw new Exception("Failed to start process");
        await process.WaitForExitAsync();
        var output = await process.StandardOutput.ReadToEndAsync();
        if (process.ExitCode != 0) throw new Exception($"exit code != 0 ({process.ExitCode}). Output: {output}");
        return output;
    }

    public static async Task Run2(string command, string arguments)
    {
        var startInfo = new ProcessStartInfo
        {
            FileName = command,
            Arguments = arguments,
            RedirectStandardOutput = true,
            RedirectStandardError = true,
            UseShellExecute = false,
            CreateNoWindow = true
        };
         var process = System.Diagnostics.Process.Start(startInfo);
        if (process is null) throw new Exception("Failed to start process");

        // Read the output and error streams
        var output = process.StandardOutput.ReadToEnd();
        var error = process.StandardError.ReadToEnd();

        await process.WaitForExitAsync();

        // Display the output and error
        Console.WriteLine("Output:");
        Console.WriteLine(output);

        Console.WriteLine("Error:");
        Console.WriteLine(error);
    }
}

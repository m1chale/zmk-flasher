
using Dotcore.FileSystem.Directory;
using ZmkFlasher.Records;

namespace ZmkFlasher.Lib;

internal static class Mount
{
    public static async Task Run(Device device, Info mountPoint, string password)
    {
        var result = await Process.Run($"echo {password} | sudo -S mount", $"/dev/{device.Name} {mountPoint.Path}");
        Console.WriteLine(result);
    }
}

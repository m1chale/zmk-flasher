
using Dotcore.FileSystem.Directory;
using ZmkFlasher.Records;

namespace ZmkFlasher.Lib;

internal static class Mount
{
    public static async Task Run(Device device, Info mountPoint)
    {
        var result = await Process.Run("sudo mount", $"/dev/{device.Name} {mountPoint.Path}");
        Console.WriteLine(result);
    }
}

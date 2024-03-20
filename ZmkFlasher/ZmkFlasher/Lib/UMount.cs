
using ZmkFlasher.Records;

namespace ZmkFlasher.Lib;

internal static class UMount
{
    public static async Task Run(Device device)
    {
        var result = await Process.Run("umount", $"/dev/{device.Name}");
        Console.WriteLine(result);
    }
}

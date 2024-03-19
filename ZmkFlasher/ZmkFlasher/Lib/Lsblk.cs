using Dotcore.FileSystem.Directory;
using System.Text.Json;
using ZmkFlasher.Records;

namespace ZmkFlasher.Lib;

internal static class Lsblk
{
    public static async Task<IEnumerable<Device>> Run()
    {
        var json = await Process.Run("lsblk", "-lf --json");
        Console.WriteLine(json);
        var data = JsonSerializer.Deserialize<LsblkData>(json);
        return data.BlockDevices.Select(d => new Device(d.Name, d.Label, d.Mountpoints.Select(m => m.ToDirectoryInfo())));
    }
}


internal record LsblkData(IEnumerable<BlockDevice> BlockDevices);
internal record BlockDevice(string Name, string Label, IEnumerable<string> Mountpoints);

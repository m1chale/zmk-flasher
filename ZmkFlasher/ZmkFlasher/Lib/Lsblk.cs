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
        if(data is null) throw new Exception("Failed to parse lsblk output");
        Console.WriteLine(data);
        return data.blockdevices.Select(d => new Device(d.name, d.label, d.mountpoints?.Select(m => m.ToDirectoryInfo()).ToArray() ?? []));
    }
}


internal record LsblkData(IEnumerable<BlockDevice> blockdevices);
internal record BlockDevice(string name, string label, string[]? mountpoints);

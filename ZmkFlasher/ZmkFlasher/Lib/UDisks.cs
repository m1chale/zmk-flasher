using Dotcore.FileSystem.Directory;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using ZmkFlasher.Records;

namespace ZmkFlasher.Lib;

internal static class UDisks
{
    public static async Task<Info> Mount(Device device)
    {
        var result = await Process.Run("udisksctl", $"mount -b /dev/{device.Name}");
        var path = result.Split("at").Last().Replace(" ", "");
        return path.ToDirectoryInfo();
    }

    public static async Task Unmount(Device device)
    {
        await Process.Run("udisksctl", $"unmount -b /dev/{device.Name}");
    }
}

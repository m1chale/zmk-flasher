
using File = Dotcore.FileSystem.File;
using Directory = Dotcore.FileSystem.Directory;
using System.Xml.Linq;
using Dotcore.FileSystem.File;
using ZmkFlasher.Lib;
using ZmkFlasher.Records;
using Dotcore.FileSystem.Directory;
using Dotcore.FileSystem;

namespace ZmkFlasher.WaitRemovableDevice;

internal interface IWaitAndCopy
{
    public static IWaitAndCopy Instance => OperatingSystem.IsWindows() ? new WaitAndCopyWindows() : new WaitAndCopyLinux();
    public Task WaitForDeviceAndCopy(string Label, File.Info Firmware, string password);
}

public class WaitAndCopyWindows : IWaitAndCopy
{
    public async Task WaitForDeviceAndCopy(string label, File.Info firmware, string password)
    {
        var directory = await WaitForDevice(label);
        firmware.CopyTo(directory);
        Console.WriteLine($"Flashed {label}");
    }

    private static Task<Directory.Info> WaitForDevice(string volumeLabel) => Task.Run(async () =>
    {
        while (true)
        {
            foreach (DriveInfo drive in DriveInfo.GetDrives())
            {
                if (drive.DriveType != DriveType.Removable) continue;
                if (drive.VolumeLabel != volumeLabel) continue;
                return new Directory.Info(drive.RootDirectory.FullName);
            }
            await Task.Delay(1000);
        }
    });
}

public class WaitAndCopyLinux : IWaitAndCopy
{
    public Task WaitForDeviceAndCopy(string Label, File.Info Firmware, string password) => TemporaryDirectory.With(async directory =>
    {
        
        Device? device;
        while (true)
        {
            var devices = await Lsblk.Run();
            device = devices.SingleOrDefault(d => d.Label == Label);
            if (device != null) break;
            await Task.Delay(1000);
        }
        Console.WriteLine($"Device {device.Label} found");

        if (device.MountPoints.Length == 0)
        {
            directory.EnsureExists();
            await Mount.Run(device, directory, password);
            Console.WriteLine($"Mounted {device.Label}");
        }else
        {
            Console.WriteLine($"Already mounted {device.Label}. {string.Join(",", device.MountPoints.Select(m => m.Path))}");
        }
        var json = await Process.Run("lsblk", "-lf --json");
        Console.WriteLine(json);
        //Firmware.CopyTo(directory);
        await Task.Delay(TimeSpan.FromMilliseconds(500));
    });
}

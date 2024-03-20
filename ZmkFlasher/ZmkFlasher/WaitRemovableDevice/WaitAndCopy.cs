
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
    public async Task WaitForDeviceAndCopy(string Label, File.Info Firmware, string password)
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

        var isMounted = device.MountPoints.Length > 0;
        if (isMounted) Console.WriteLine($"Already mounted {device.Label}. {string.Join(",", device.MountPoints.Select(m => m.Path))}");
        else
        {
            var directory = await UDisks.Mount(device);
            Console.WriteLine($"Mounted {device.Label}");
        }

        //Firmware.CopyTo(directory);
        await Task.Delay(TimeSpan.FromMilliseconds(500));
        if (isMounted)
        {
            try
            {
                await UDisks.Unmount(device);
                Console.WriteLine($"Unmounted {device.Label}");
            }
            catch (Exception any)
            {
                Console.WriteLine($"Failed to unmount {device.Label}: {any.Message}");
            }
        }
    }
}

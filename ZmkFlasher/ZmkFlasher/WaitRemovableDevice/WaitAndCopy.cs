
using File = Dotcore.FileSystem.File;
using Directory = Dotcore.FileSystem.Directory;
using Dotcore.FileSystem.File;
using ZmkFlasher.Lib;
using ZmkFlasher.Records;

namespace ZmkFlasher.WaitRemovableDevice;

internal interface IWaitAndCopy
{
    public static IWaitAndCopy Instance => OperatingSystem.IsWindows() ? new WaitAndCopyWindows() : new WaitAndCopyLinux();
    public Task WaitForDeviceAndCopy(string Label, File.Info Firmware);
}

public class WaitAndCopyWindows : IWaitAndCopy
{
    public async Task WaitForDeviceAndCopy(string label, File.Info firmware)
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
            await Task.Delay(100);
        }
    });
}

public class WaitAndCopyLinux : IWaitAndCopy
{
    public async Task WaitForDeviceAndCopy(string Label, File.Info Firmware)
    {

        Device? device;
        while (true)
        {
            var devices = await Lsblk.Run();
            device = devices.SingleOrDefault(d => d.Label == Label);
            if (device != null) break;
            await Task.Delay(100);
        }
        Console.WriteLine($"Device {device.Label} found");

        var isMounted = device.MountPoints.Length > 0;
        var directory = device.MountPoints.SingleOrDefault();
        if (isMounted) Console.WriteLine($"Already mounted {device.Label}. {string.Join(",", device.MountPoints.Select(m => m.Path))}");
        else
        {
            directory = await UDisks.Mount(device);
            Console.WriteLine($"Mounted {device.Label}");
        }
        if (directory == null) throw new Exception("Failed to mount device");

        Console.WriteLine($"Copying firmware to {directory}");
        Firmware.CopyTo(directory);
        try
        {
            await UDisks.Unmount(device);
            Console.WriteLine($"Unmounted {device.Label}");
        }
        catch (Exception any)
        {
            Console.WriteLine($"Failed to unmount {device.Label}: {any}");
        }
    }
}

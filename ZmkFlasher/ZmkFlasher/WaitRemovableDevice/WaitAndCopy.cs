
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
            await Task.Delay(1000);
        }
    });
}

public class WaitAndCopyLinux : IWaitAndCopy
{
    public Task WaitForDeviceAndCopy(string Label, File.Info Firmware) => TemporaryDirectory.With(async directory =>
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
            await Mount.Run(device, directory);
            Console.WriteLine($"Mounted {device.Label}");
        }

        //Firmware.CopyTo(directory);
    });
}



internal class WaitAndCopy
{
    public static Task WaitAndCopyFirmware( (string VolumeLabel, File.Info Firmware) left,  (string VolumeLabel, File.Info Firmware) right, bool verbose = false)
    {
        Console.WriteLine($"Connect left or right bootloader");
        return Task.WhenAll(
            WaitAndCopyFirmware("left", left.VolumeLabel, left.Firmware, verbose),
            WaitAndCopyFirmware("right", right.VolumeLabel, right.Firmware, verbose));
    }

    public static Task WaitAndCopyFirmware(string name, string volumeLabel, File.Info firmware, bool verbose = false) => Task.Run(async () =>
    {
        var directory = await Wait.ForDeviceWindows(volumeLabel, verbose);
        firmware.CopyTo(directory);
        Console.WriteLine($"Flashed {name}");
    });
}

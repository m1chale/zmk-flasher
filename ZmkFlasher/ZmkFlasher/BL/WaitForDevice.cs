
using Dotcore.FileSystem.Directory;
using System.Reflection.Emit;
using ZmkFlasher.Lib;
using ZmkFlasher.Records;
using Directory= Dotcore.FileSystem.Directory;

namespace ZmkFlasher.WaitRemovableDevice;

public interface IWaitForDevice
{
    public static IWaitForDevice Instance => OperatingSystem.IsWindows() ? new WaitForDeviceWindows() : new WaitForDeviceLinux();
    public Task<Device> WaitForDevice(string Label);
}

internal class WaitForDeviceLinux : IWaitForDevice
{
    public async Task<Device> WaitForDevice(string volumeLabel)
    {
        while (true)
        {
            var devices = await Lsblk.Run();
            var device = devices.SingleOrDefault(d => d.Label?.Equals(volumeLabel) ?? false);
            if (device != null)
            {
                Console.WriteLine($"Found {device.Label}");
                return new Device(device.Name, device.Label, device.MountPoints);
            }
            
            await Task.Delay(100);
        }
    }
}

internal class WaitForDeviceWindows : IWaitForDevice
{
    public async Task<Device> WaitForDevice(string volumeLabel)
    {

        while (true)
        {
            foreach (DriveInfo drive in DriveInfo.GetDrives())
            {
                if (drive.DriveType != DriveType.Removable) continue;
                if (drive.VolumeLabel != volumeLabel) continue;
                Console.WriteLine($"Found {drive.VolumeLabel}");
                return new Device(drive.Name, drive.VolumeLabel, [drive.RootDirectory]);
            }
            await Task.Delay(100);
        }
    }
}

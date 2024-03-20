
using ZmkFlasher.Lib;
using ZmkFlasher.Records;

namespace ZmkFlasher.BL;


internal interface ICleanUpDevice
{
    public static ICleanUpDevice Instance => OperatingSystem.IsWindows() ? new CleanUpDeviceWindows() : new CleanUpDeviceLinux();
    public Task CleanUp(Device device);
}

internal class CleanUpDeviceWindows : ICleanUpDevice
{
    public Task CleanUp(Device device) => Task.CompletedTask;
}

internal class CleanUpDeviceLinux : ICleanUpDevice
{
    public async Task CleanUp(Device device)
    {
        var devices = await Lsblk.Run();
        if (!devices.Any(device => device.Label?.Equals(device.Label) ?? false)) return;
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

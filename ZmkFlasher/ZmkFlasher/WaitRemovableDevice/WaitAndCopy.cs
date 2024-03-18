
using Dotcore.FileSystem.File;
using File = Dotcore.FileSystem.File;

namespace ZmkFlasher.WaitRemovableDevice;

internal class WaitAndCopy
{
    public static Task WaitAndCopyFirmware( (string VolumeLabel, File.Info Firmware) left,  (string VolumeLabel, File.Info Firmware) right)
    {
        Console.WriteLine($"Connect left or right bootloader");
        return Task.WhenAll(
            WaitAndCopyFirmware("left", left.VolumeLabel, left.Firmware),
                       WaitAndCopyFirmware("right", right.VolumeLabel, right.Firmware));
    }

    public static Task WaitAndCopyFirmware(string name, string volumeLabel, File.Info firmware) => Task.Run(async () =>
    {
        var directory = await Wait.ForDevice(volumeLabel);
        firmware.CopyTo(directory);
        Console.WriteLine($"Flashed {name}");
    });
}

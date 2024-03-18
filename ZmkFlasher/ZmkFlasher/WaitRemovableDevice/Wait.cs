using Directory = Dotcore.FileSystem.Directory;

namespace ZmkFlasher.WaitRemovableDevice;

internal static class Wait
{
    public static Task<Directory.Info> ForDevice(string volumeLabel, bool verbose = false) => Task.Run(async () =>
    {
        while (true)
        {
            foreach (DriveInfo drive in DriveInfo.GetDrives())
            {
                if(drive.DriveType != DriveType.Removable) continue;
                if (verbose) Console.WriteLine($"Drive-Label: {drive.VolumeLabel}");
                if (drive.VolumeLabel != volumeLabel) continue;
                return new Directory.Info(drive.RootDirectory.FullName);
            }
            await Task.Delay(1000);
        }
    });
}


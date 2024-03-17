using Directory = Dotcore.FileSystem.Directory;

namespace ZmkFlasher.WaitRemovableDevice;

internal static class Wait
{
    public static Task<Directory.Info> ForDevice(string volumeLabel) => Task.Run(async () =>
    {
        while (true)
        {
            foreach (DriveInfo drive in DriveInfo.GetDrives())
            {
                if (drive.DriveType == DriveType.Removable && drive.VolumeLabel == volumeLabel)
                {
                    return new Directory.Info(drive.RootDirectory.FullName);
                }
            }
            await Task.Delay(1000);
        }
    });
}


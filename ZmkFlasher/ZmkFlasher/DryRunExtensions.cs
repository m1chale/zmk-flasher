
using Dotcore.FileSystem.File;
using File= Dotcore.FileSystem.File;
using Directory = Dotcore.FileSystem.Directory;

namespace ZmkFlasher;

internal static class DryRunExtensions
{
    public static bool IsDryRun { get; set; } = false;

    public static void ThrowIfNotExistsOrDryRun(this File.Info info)
    {
        if (IsDryRun)
        {
            Console.WriteLine($"Dry run: Would have checked if {info.Path} exists");
            return;
        }
        info.ThrowIfNotExists();
    }

    public static void CopyToOrDryRun(this File.Info source, Directory.Info destination)
    {
        if (IsDryRun)
        {
            Console.WriteLine($"Dry run: Would have copied {source.Path} to {destination.Path}");
            return;
        }
        source.CopyTo(destination);
    }
}

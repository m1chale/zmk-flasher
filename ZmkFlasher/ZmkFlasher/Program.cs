
using Dotcore.FileSystem.File;
using ZmkFlasher.Arguments;
using ZmkFlasher.Lib;
using ZmkFlasher.WaitRemovableDevice;

var drives = await Lsblk.Run();
foreach (var drive in drives)
{
    Console.WriteLine(drive);
}

var result = CommandLine.Parser.Default.ParseArguments<StringArguments>(args);
if (result.Errors.Any()) throw new Exception($"invalid arguments {result}");
var arguments = result.Value.ToTypedArguments();

arguments.LeftFirmware.ThrowIfNotExists();
arguments.RightFirmware.ThrowIfNotExists();

await WaitAndCopy.WaitAndCopyFirmware(("GLV80LHBOOT", arguments.LeftFirmware), ("GLV80RHBOOT", arguments.RightFirmware), arguments.Verbose);
Console.WriteLine("Done");